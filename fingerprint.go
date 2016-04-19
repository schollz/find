package main

import (
	"fmt"
	"log"
	"strconv"

	"net/http"
	"path"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

// Fingerprint is the prototypical information from the fingerprinting device
type Fingerprint struct {
	Group           string   `json:"group"`
	Username        string   `json:"username"`
	Location        string   `json:"location"`
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
	return compressByte(dumped)
}

// compression 30 us -> 600 us
func loadFingerprint(jsonByte []byte) Fingerprint {
	res := Fingerprint{}
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

		timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
		err = bucket.Put([]byte(timestamp), dumpFingerprint(res))
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}
		return err
	})
	return err
}

func handleFingerprint(c *gin.Context) {
	var jsonFingerprint Fingerprint
	if c.BindJSON(&jsonFingerprint) == nil {
		cleanFingerprint(&jsonFingerprint)
		if jsonFingerprint.Location != "" {
			go putFingerprintIntoDatabase(jsonFingerprint, "fingerprints")
			isLearning[strings.ToLower(jsonFingerprint.Group)] = true
			Debug.Println("Inserted fingerprint for " + jsonFingerprint.Username + " (" + jsonFingerprint.Group + ") at " + jsonFingerprint.Location)
			c.JSON(http.StatusOK, gin.H{"message": "Inserted fingerprint containing " + strconv.Itoa(len(jsonFingerprint.WifiFingerprint)) + " APs for " + jsonFingerprint.Username + " at " + jsonFingerprint.Location, "success": true})
		} else {
			trackFingerprint(c)
		}
	}
}

func trackFingerprint(c *gin.Context) {
	var jsonFingerprint Fingerprint
	if c.BindJSON(&jsonFingerprint) == nil {
		cleanFingerprint(&jsonFingerprint)
		if wasLearning, ok := isLearning[strings.ToLower(jsonFingerprint.Group)]; ok {
			if wasLearning {
				Debug.Println("Was learning, calculating priors")
				group := strings.ToLower(jsonFingerprint.Group)
				isLearning[group] = false
				optimizePriorsThreaded(group)
				if _, ok := usersCache[group]; ok {
					if len(usersCache[group]) == 0 {
						usersCache[group] = append([]string{}, strings.ToLower(jsonFingerprint.Username))
					}
				}
			}
		}
		locationGuess, bayes := calculatePosterior(jsonFingerprint, *NewFullParameters())
		jsonFingerprint.Location = locationGuess
		go putFingerprintIntoDatabase(jsonFingerprint, "fingerprints-track")
		positions := [][]string{}
		positions1 := []string{}
		positions2 := []string{}
		positions1 = append(positions1, locationGuess)
		positions2 = append(positions2, " ")
		positions = append(positions, positions1)
		positions = append(positions, positions2)
		var userJSON UserPositionJSON
		userJSON.Location = locationGuess
		userJSON.Bayes = bayes
		userJSON.Time = time.Now().String()
		userPositionCache[strings.ToLower(jsonFingerprint.Group)+strings.ToLower(jsonFingerprint.Username)] = userJSON
		Debug.Println("Tracking fingerprint for " + jsonFingerprint.Username + " (" + jsonFingerprint.Group + ") at " + jsonFingerprint.Location + " (guess)")
		c.JSON(http.StatusOK, gin.H{"message": "Calculated location: " + locationGuess, "success": true, "location": locationGuess})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Something went wrong", "success": false})
	}
}
