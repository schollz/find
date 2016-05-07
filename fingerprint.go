// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// fingerprint.go contains structures and functions for handling fingerprints.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"

	"net/http"
	"path"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

// Fingerprint is the prototypical information from the fingerprinting device
// IF you change Fingerprint, follow these steps to re-generate fingerprint_ffjson.go
// find ./ -name "*.go" -type f | xargs sed -i  's/package main/package main/g'
// Uncomment json.Marshal/Unmarshal functions
// $GOPATH/bin/ffjson fingerprint.go
// find ./ -name "*.go" -type f | xargs sed -i  's/package main/package main/g'
// Comment json.Marshal/Unmarshal functions
type Fingerprint struct {
	Group           string   `json:"group"`
	Username        string   `json:"username"`
	Location        string   `json:"location"`
	Timestamp       int64    `json:"timestamp"`
	WifiFingerprint []Router `json:"wifi-fingerprint"`
}

// Router is the router information for each invdividual mac address
type Router struct {
	Mac  string `json:"mac"`
	Rssi int    `json:"rssi"`
}

var jsonExample = `{
	"group": "whatevergroup",
	"username": "iamauser",
	"location": null,
	"wififingerprint": [{
		"mac": "AA:AA:AA:AA:AA:AA",
		"rssi": -45
	}, {
		"mac": "BB:BB:BB:BB:BB:BB",
		"rssi": -55
	}]
}`

// compression 9 us -> 900 us
func dumpFingerprint(res Fingerprint) []byte {
	dumped, _ := res.MarshalJSON()
	//dumped, _ := json.Marshal(res)
	return compressByte(dumped)
}

// compression 30 us -> 600 us
func loadFingerprint(jsonByte []byte) Fingerprint {
	res := Fingerprint{}
	//json.Unmarshal(decompressByte(jsonByte), res)
	res.UnmarshalJSON(decompressByte(jsonByte))
	return res
}

func cleanFingerprint(res *Fingerprint) {
	res.Group = strings.ToLower(res.Group)
	res.Location = strings.ToLower(res.Location)
	res.Username = strings.ToLower(res.Username)
	deleteIndex := -1
	for r := range res.WifiFingerprint {
		if res.WifiFingerprint[r].Rssi >= 0 { // https://stackoverflow.com/questions/15797920/how-to-convert-wifi-signal-strength-from-quality-percent-to-rssi-dbm
			res.WifiFingerprint[r].Rssi = int(res.WifiFingerprint[r].Rssi/2) - 100
		}
		if res.WifiFingerprint[r].Mac == "00:00:00:00:00:00" {
			deleteIndex = r
		}
	}
	if deleteIndex > -1 {
		res.WifiFingerprint[deleteIndex] = res.WifiFingerprint[len(res.WifiFingerprint)-1]
		res.WifiFingerprint = res.WifiFingerprint[:len(res.WifiFingerprint)-1]
	}
}

func putFingerprintIntoDatabase(res Fingerprint, database string) error {
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, res.Group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(database))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		if res.Timestamp == 0 {
			res.Timestamp = time.Now().UnixNano()
		}
		err = bucket.Put([]byte(strconv.FormatInt(res.Timestamp, 10)), dumpFingerprint(res))
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}
		return err
	})
	return err
}

func trackFingerprintPOST(c *gin.Context) {
	var jsonFingerprint Fingerprint
	if c.BindJSON(&jsonFingerprint) == nil {
		message, success, locationGuess, bayes := trackFingerprint(jsonFingerprint)
		if success {
			c.JSON(http.StatusOK, gin.H{"message": message, "success": true, "location": locationGuess, "bayes": bayes})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": message, "success": false})
		}
	} else {
		Warning.Println("Could not bind JSON")
		c.JSON(http.StatusOK, gin.H{"message": "Could not bind JSON", "success": false})
	}
}

func learnFingerprintPOST(c *gin.Context) {
	var jsonFingerprint Fingerprint
	if c.BindJSON(&jsonFingerprint) == nil {
		message, success := learnFingerprint(jsonFingerprint)
		Debug.Println(message)
		if !success {
			Debug.Println(jsonFingerprint)
		}
		c.JSON(http.StatusOK, gin.H{"message": message, "success": success})
	} else {
		Warning.Println("Could not bind JSON")
		c.JSON(http.StatusOK, gin.H{"message": "Could not bind JSON", "success": false})
	}
}

func learnFingerprint(jsonFingerprint Fingerprint) (string, bool) {
	cleanFingerprint(&jsonFingerprint)
	if len(jsonFingerprint.Group) == 0 {
		return "Need to define your group name in request, see API", false
	}
	if len(jsonFingerprint.WifiFingerprint) == 0 {
		return "No fingerprints found to insert, see API", false
	}
	putFingerprintIntoDatabase(jsonFingerprint, "fingerprints")
	isLearning[strings.ToLower(jsonFingerprint.Group)] = true
	message := "Inserted fingerprint containing " + strconv.Itoa(len(jsonFingerprint.WifiFingerprint)) + " APs for " + jsonFingerprint.Username + " at " + jsonFingerprint.Location
	return message, true
}

func trackFingerprint(jsonFingerprint Fingerprint) (string, bool, string, map[string]float64) {
	bayes := make(map[string]float64)
	cleanFingerprint(&jsonFingerprint)
	if !groupExists(jsonFingerprint.Group) || len(jsonFingerprint.Group) == 0 {
		return "You should insert fingerprints before tracking", false, "", bayes
	}
	if len(jsonFingerprint.WifiFingerprint) == 0 {
		return "No fingerprints found to track, see API", false, "", bayes
	}
	if len(jsonFingerprint.Username) == 0 {
		return "No username defined, see API", false, "", bayes
	}
	if wasLearning, ok := isLearning[strings.ToLower(jsonFingerprint.Group)]; ok {
		if wasLearning {
			Debug.Println("Was learning, calculating priors")
			group := strings.ToLower(jsonFingerprint.Group)
			isLearning[group] = false
			optimizePriorsThreaded(group)
			dumpFingerprintsSVM(group)
			calculateSVM(group)
			if _, ok := usersCache[group]; ok {
				if len(usersCache[group]) == 0 {
					usersCache[group] = append([]string{}, strings.ToLower(jsonFingerprint.Username))
				}
			}
		}
	}
	locationGuess1, bayes := calculatePosterior(jsonFingerprint, *NewFullParameters())
	locationGuess2, bayes2 := classify(jsonFingerprint)
	jsonFingerprint.Location = locationGuess1
	putFingerprintIntoDatabase(jsonFingerprint, "fingerprints-track")
	positions := [][]string{}
	positions1 := []string{}
	positions2 := []string{}
	positions1 = append(positions1, locationGuess1)
	positions2 = append(positions2, " ")
	positions = append(positions, positions1)
	positions = append(positions, positions2)
	var userJSON UserPositionJSON
	userJSON.Location = locationGuess1
	userJSON.Bayes = bayes
	userJSON.Time = time.Now().String()
	userPositionCache[strings.ToLower(jsonFingerprint.Group)+strings.ToLower(jsonFingerprint.Username)] = userJSON
	Debug.Println("Tracking fingerprint containing " + strconv.Itoa(len(jsonFingerprint.WifiFingerprint)) + " APs for " + jsonFingerprint.Username + " (" + jsonFingerprint.Group + ") at " + jsonFingerprint.Location + " (guess)")
	if RuntimeArgs.Mqtt {
		type FingerprintResponse struct {
			LocationGuess string             `json:"location"`
			Bayes         map[string]float64 `json:"bayes"`
		}
		mqttMessage, _ := json.Marshal(FingerprintResponse{
			LocationGuess: locationGuess1,
			Bayes:         bayes,
		})
		go sendMQTTLocation(string(mqttMessage), jsonFingerprint.Group, jsonFingerprint.Username)
	}
	percentGuess1 := float64(0)
	total := float64(0)
	for _, locBayes := range bayes {
		total += math.Exp(locBayes)
		if locBayes > percentGuess1 {
			percentGuess1 = locBayes
		}
	}
	percentGuess1 = math.Exp(bayes[locationGuess1]) / total * 100.0
	percentGuess2 := int(100 * math.Exp(bayes2[locationGuess2]))
	if percentGuess2 > 100 {
		percentGuess2 = percentGuess2 / 10
	}
	return "NB: " + locationGuess1 + " (" + strconv.Itoa(int(percentGuess1)) + "%)" + ", SVM: " + locationGuess2 + " (" + strconv.Itoa(int(percentGuess2)) + "%)", true, locationGuess1, bayes
}
