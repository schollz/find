// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// api.go handles functions that return JSON responses.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"uptime": time.Since(startTime).Seconds(), "registered": startTime.String(), "status": "standard", "num_cores": runtime.NumCPU(), "success": true})
}

// UserPositionJSON stores the a users time, location and bayes after calculatePosterior()
type UserPositionJSON struct {
	Time     interface{}        `json:"time"`
	Location interface{}        `json:"location"`
	Bayes    map[string]float64 `json:"bayes"`
	Svm      map[string]float64 `json:"svm"`
	Rf       map[string]float64 `json:"rf"`
}

func getHistoricalUserPositions(group string, user string, n int) []UserPositionJSON {
	group = strings.ToLower(group)
	user = strings.ToLower(user)

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var fingerprints []Fingerprint
	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("fingerprints-track"))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		numFound := 0
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			v2 := loadFingerprint(v)
			if v2.Username == user {
				timestampString := string(k)
				timestampUnixNano, _ := strconv.ParseInt(timestampString, 10, 64)
				UTCfromUnixNano := time.Unix(0, timestampUnixNano)
				v2.Timestamp = UTCfromUnixNano.UnixNano()
				fingerprints = append(fingerprints, v2)
				numFound++
				if numFound >= n {
					return nil
				}
			}
		}
		return fmt.Errorf("User " + user + " not found")
	})
	db.Close()

	Debug.Printf("Got history of %d fingerprints\n", len(fingerprints))
	userJSONs := make([]UserPositionJSON, len(fingerprints))
	for i, fingerprint := range fingerprints {
		var userJSON UserPositionJSON
		UTCfromUnixNano := time.Unix(0, fingerprint.Timestamp)
		userJSON.Time = UTCfromUnixNano.String()
		location, bayes := calculatePosterior(fingerprint, *NewFullParameters())
		userJSON.Location = location
		userJSON.Bayes = bayes
		// Process SVM if needed
		if RuntimeArgs.Svm {
			_, userJSON.Svm = classify(fingerprint)
		}
		// Process RF if needed
		if RuntimeArgs.RandomForests {
			userJSON.Rf = rfClassify(group, fingerprint)
		}
		userJSONs[i] = userJSON
	}
	return userJSONs
}

func getCurrentPositionOfUser(group string, user string) UserPositionJSON {
	group = strings.ToLower(group)
	user = strings.ToLower(user)
	val, ok := getUserPositionCache(group + user)
	if ok {
		return val
	}
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var userFingerprint Fingerprint
	var userJSON UserPositionJSON
	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("fingerprints-track"))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			v2 := loadFingerprint(v)
			if v2.Username == user {
				timestampString := string(k)
				timestampUnixNano, _ := strconv.ParseInt(timestampString, 10, 64)
				UTCfromUnixNano := time.Unix(0, timestampUnixNano)
				userJSON.Time = UTCfromUnixNano.String()
				userFingerprint = v2
				return nil
			}
		}
		return fmt.Errorf("User " + user + " not found")
	})
	db.Close()
	location, bayes := calculatePosterior(userFingerprint, *NewFullParameters())
	userJSON.Location = location
	userJSON.Bayes = bayes
	// Process SVM if needed
	if RuntimeArgs.Svm {
		_, userJSON.Svm = classify(userFingerprint)
	}
	if RuntimeArgs.RandomForests {
		userJSON.Rf = rfClassify(group, userFingerprint)
		log.Println("%+v", userJSON.Rf)
	}
	go setUserPositionCache(group+user, userJSON)
	return userJSON
}

func calculate(c *gin.Context) {
	group := c.DefaultQuery("group", "noneasdf")
	if group != "noneasdf" {
		if !groupExists(group) {
			c.JSON(http.StatusOK, gin.H{"message": "You should insert a fingerprint first, see documentation", "success": false})
			return
		}
		group = strings.ToLower(group)
		optimizePriorsThreaded(group)
		if RuntimeArgs.Svm {
			dumpFingerprintsSVM(group)
			err := calculateSVM(group)
			if err != nil {
				Warning.Println("Encountered error when calculating SVM")
				Warning.Println(err)
			}
		}
		if RuntimeArgs.RandomForests {
			rfLearn(group)
		}
		go resetCache("userPositionCache")
		c.JSON(http.StatusOK, gin.H{"message": "Parameters optimized.", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func userLocations(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	group := c.DefaultQuery("group", "noneasdf")
	group = strings.ToLower(group)
	if group != "noneasdf" {
		if !groupExists(group) {
			c.JSON(http.StatusOK, gin.H{"message": "You should insert fingerprints before tracking, see documentation", "success": false})
			return
		}
		users := getUsers(group)
		people := make(map[string]UserPositionJSON)
		for _, user := range users {
			people[user] = getCurrentPositionOfUser(group, user)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Correctly found", "success": true, "users": people})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}
func getUserLocations(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	group := c.DefaultQuery("group", "noneasdf")
	userQuery := c.DefaultQuery("user", "noneasdf")
	usersQuery := c.DefaultQuery("users", "noneasdf")
	nQuery := c.DefaultQuery("n", "noneasdf")
	group = strings.ToLower(group)
	if group != "noneasdf" {
		if !groupExists(group) {
			c.JSON(http.StatusOK, gin.H{"message": "You should insert fingerprints before tracking, see documentation", "success": false})
			return
		}
		people := make(map[string][]UserPositionJSON)
		allusers := getUsers(group)
		if userQuery != "noneasdf" {
			usersQuery = userQuery
		}
		users := strings.Split(strings.ToLower(usersQuery), ",")
		if users[0] == "noneasdf" {
			users = allusers
		}
		for _, user := range users {
			if !stringInSlice(user, allusers) {
				continue
			}
			if _, ok := people[user]; !ok {
				people[user] = []UserPositionJSON{}
			}
			if nQuery != "noneasdf" {
				number, _ := strconv.ParseInt(nQuery, 10, 0)
				Debug.Println("Getting history for " + user)
				people[user] = append(people[user], getHistoricalUserPositions(group, user, int(number))...)
			} else {
				people[user] = append(people[user], getCurrentPositionOfUser(group, user))
			}
		}
		message := "Correctly found locations."
		if len(people) == 0 {
			message = "No users found for username " + strings.Join(users, " or ")
			people = nil
		}

		c.JSON(http.StatusOK, gin.H{"message": message, "success": true, "users": people})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func migrateDatabase(c *gin.Context) {
	fromDB := strings.ToLower(c.DefaultQuery("from", "noneasdf"))
	toDB := strings.ToLower(c.DefaultQuery("to", "noneasdf"))
	Debug.Printf("Migrating %s to %s.\n", fromDB, toDB)
	if !exists(path.Join(RuntimeArgs.SourcePath, fromDB+".db")) {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Can't migrate from " + fromDB + ", it does not exist."})
		return
	}
	if !exists(path.Join(RuntimeArgs.SourcePath, toDB)) {
		CopyFile(path.Join(RuntimeArgs.SourcePath, fromDB+".db"), path.Join(RuntimeArgs.SourcePath, toDB+".db"))
	} else {
		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, fromDB+".db"), 0664, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		db2, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, toDB+".db"), 0664, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db2.Close()

		db2.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("fingerprints"))
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					bucket.Put(k, v)
				}
				return nil
			})
			return nil
		})

		db2.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints-track"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			db.View(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("fingerprints-track"))
				c := b.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					bucket.Put(k, v)
				}
				return nil
			})
			return nil
		})
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully migrated " + fromDB + " to " + toDB})
}

func deleteDatabase(c *gin.Context) {
	group := strings.TrimSpace(strings.ToLower(c.DefaultQuery("group", "noneasdf")))
	if exists(path.Join(RuntimeArgs.SourcePath, group+".db")) {
		os.Remove(path.Join(RuntimeArgs.SourcePath, group+".db"))
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Successfully deleted " + group})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Group does not exist"})
	}
}

func putMixinOverride(c *gin.Context) {
	group := strings.ToLower(c.DefaultQuery("group", "noneasdf"))
	newMixin := c.DefaultQuery("mixin", "none")
	if group != "noneasdf" {
		newMixinFloat, err := strconv.ParseFloat(newMixin, 64)
		if err == nil {
			err2 := setMixinOverride(group, newMixinFloat)
			if err2 == nil {
				optimizePriorsThreaded(strings.ToLower(group))
				c.JSON(http.StatusOK, gin.H{"success": true, "message": "Overriding mixin for " + group + ", now set to " + newMixin})
			} else {
				c.JSON(http.StatusOK, gin.H{"success": false, "message": err2.Error()})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"success": false, "message": err.Error()})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func editNetworkName(c *gin.Context) {
	group := c.DefaultQuery("group", "noneasdf")
	oldname := c.DefaultQuery("oldname", "none")
	newname := c.DefaultQuery("newname", "none")
	if group != "noneasdf" {
		Debug.Println("Attempting renaming ", group, oldname, newname)
		renameNetwork(group, oldname, newname)
		optimizePriors(group)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Finished"})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func editName(c *gin.Context) {
	group := c.DefaultQuery("group", "noneasdf")
	location := c.DefaultQuery("location", "none")
	newname := c.DefaultQuery("newname", "none")
	if group != "noneasdf" {
		toUpdate := make(map[string]string)
		numChanges := 0

		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Location == location {
						v2.Location = newname
						toUpdate[string(k)] = string(dumpFingerprint(v2))
					}
				}
			}
			return nil
		})

		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for k, v := range toUpdate {
				bucket.Put([]byte(k), []byte(v))
			}
			return nil
		})

		numChanges += len(toUpdate)

		toUpdate = make(map[string]string)

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints-track"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Location == location {
						v2.Location = newname
						toUpdate[string(k)] = string(dumpFingerprint(v2))
					}
				}
			}
			return nil
		})

		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints-track"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for k, v := range toUpdate {
				bucket.Put([]byte(k), []byte(v))
			}
			return nil
		})

		db.Close()
		numChanges += len(toUpdate)
		optimizePriorsThreaded(strings.ToLower(group))

		c.JSON(http.StatusOK, gin.H{"message": "Changed name of " + strconv.Itoa(numChanges) + " things", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func editUserName(c *gin.Context) {
	group := strings.ToLower(c.DefaultQuery("group", "noneasdf"))
	user := strings.ToLower(c.DefaultQuery("user", "none"))
	newname := strings.ToLower(c.DefaultQuery("newname", "none"))
	if group != "noneasdf" {
		toUpdate := make(map[string]string)
		numChanges := 0

		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Username == user {
						v2.Username = newname
						toUpdate[string(k)] = string(dumpFingerprint(v2))
					}
				}
			}
			return nil
		})

		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for k, v := range toUpdate {
				bucket.Put([]byte(k), []byte(v))
			}
			return nil
		})

		numChanges += len(toUpdate)

		toUpdate = make(map[string]string)

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints-track"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Username == user {
						v2.Username = newname
						toUpdate[string(k)] = string(dumpFingerprint(v2))
					}
				}
			}
			return nil
		})

		db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists([]byte("fingerprints-track"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}

			for k, v := range toUpdate {
				bucket.Put([]byte(k), []byte(v))
			}
			return nil
		})

		db.Close()
		numChanges += len(toUpdate)

		// reset the cache (cache.go)
		go resetCache("usersCache")
		go resetCache("userPositionCache")

		c.JSON(http.StatusOK, gin.H{"message": "Changed name of " + strconv.Itoa(numChanges) + " things", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func deleteLocation(c *gin.Context) {
	group := strings.ToLower(c.DefaultQuery("group", "noneasdf"))
	location := strings.ToLower(c.DefaultQuery("location", "none"))
	if group != "noneasdf" {
		numChanges := 0

		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Location == location {
						b.Delete(k)
						numChanges++
					}
				}
			}
			return nil
		})

		db.Close()
		optimizePriorsThreaded(strings.ToLower(group))

		c.JSON(http.StatusOK, gin.H{"message": "Deleted " + strconv.Itoa(numChanges) + " locations", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func deleteLocations(c *gin.Context) {
	group := strings.ToLower(c.DefaultQuery("group", "noneasdf"))
	locationsQuery := strings.ToLower(c.DefaultQuery("names", "none"))
	if group != "noneasdf" && locationsQuery != "none" {
		locations := strings.Split(strings.ToLower(locationsQuery), ",")
		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		numChanges := 0
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					for _, location := range locations {
						if v2.Location == location {
							b.Delete(k)
							numChanges++
							break
						}
					}
				}
			}
			return nil
		})
		db.Close()
		optimizePriorsThreaded(strings.ToLower(group))
		c.JSON(http.StatusOK, gin.H{"message": "Deleted " + strconv.Itoa(numChanges) + " locations", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Need to provide group and location list. DELETE /locations?group=X&names=Y,Z,W"})
	}
}

func deleteUser(c *gin.Context) {
	group := strings.ToLower(c.DefaultQuery("group", "noneasdf"))
	user := strings.ToLower(c.DefaultQuery("user", "noneasdf"))
	if group != "noneasdf" && user != "noneasdf" {
		numChanges := 0

		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("fingerprints-track"))
			if b != nil {
				c := b.Cursor()
				for k, v := c.Last(); k != nil; k, v = c.Prev() {
					v2 := loadFingerprint(v)
					if v2.Username == user {
						b.Delete(k)
						numChanges++
					}
				}
			}
			return nil
		})

		db.Close()

		// reset the cache (cache.go)
		go resetCache("usersCache")
		go resetCache("userPositionCache")

		c.JSON(http.StatusOK, gin.H{"message": "Deletes " + strconv.Itoa(numChanges) + " things " + " with user " + user, "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

type whereAmIJson struct {
	Group string `json:"group"`
	User  string `json:"user"`
}

func whereAmI(c *gin.Context) {
	var jsonData whereAmIJson
	if c.BindJSON(&jsonData) == nil {
		defer timeTrack(time.Now(), "getUniqueMacs")
		db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, jsonData.Group+".db"), 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		locations := []string{}
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("fingerprints-track"))
			c := b.Cursor()
			for k, v := c.Last(); k != nil; k, v = c.Prev() {
				v2 := loadFingerprint(v)
				if v2.Username == jsonData.User {
					locations = append(locations, v2.Location)
				}
				if len(locations) > 2 {
					break
				}
			}
			return nil
		})
		// jsonLocations, _ := json.Marshal(locations)
		message := "Found user"
		if len(locations) == 0 {
			message = "No locations found."
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": message, "group": jsonData.Group, "user": jsonData.User, "locations": locations})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Could not bind JSON - did you not send it as a JSON?", "success": false})
	}
}
