package main

import (
	"fmt"
	"log"
	"net/http"
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
	c.JSON(http.StatusOK, gin.H{"uptime": time.Since(startTime).Seconds(), "registered": startTime.String(), "status": "standard", "num_cores": runtime.NumCPU()})
}

// UserPositionJSON stores the a users time, location and bayes after calculatePosterior()
type UserPositionJSON struct {
	Time     interface{}        `json:"time"`
	Location interface{}        `json:"location"`
	Bayes    map[string]float64 `json:"bayes"`
}

func getPositionBreakdown(group string, user string) UserPositionJSON {
	group = strings.ToLower(group)
	user = strings.ToLower(user)
	if val, ok := userPositionCache[group+user]; ok {
		return val
	}
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var userJSON UserPositionJSON
	var fullJSON Fingerprint
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
				fullJSON = v2
				return nil
			}
		}
		return fmt.Errorf("User " + user + " not found")
	})
	db.Close()
	if err == nil {
		location, bayes := calculatePosterior(fullJSON, *NewFullParameters())
		userJSON.Location = location
		userJSON.Bayes = bayes
	}
	userPositionCache[group+user] = userJSON
	return userJSON
}

func calculate(c *gin.Context) {
	group := c.DefaultQuery("group", "noneasdf")
	if group != "noneasdf" {
		optimizePriorsThreaded(strings.ToLower(group))
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
		users := getUsers(group)
		people := make(map[string]UserPositionJSON)
		for _, user := range users {
			people[user] = getPositionBreakdown(group, user)
		}
		c.JSON(http.StatusOK, gin.H{"message": "Correctly found", "success": true, "users": people})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func editNetworkName(c *gin.Context) {
	group := c.DefaultQuery("group", "noneasdf")
	oldname := c.DefaultQuery("oldname", "none")
	newname := c.DefaultQuery("newname", "none")
	if group != "noneasdf" {
		fmt.Println("Attempting renaming ", group, oldname, newname)
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
		usersCache = make(map[string][]string)
		userPositionCache = make(map[string]UserPositionJSON)

		c.JSON(http.StatusOK, gin.H{"message": "Changed name of " + strconv.Itoa(numChanges) + " things", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
	}
}

func deleteName(c *gin.Context) {
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

		c.JSON(http.StatusOK, gin.H{"message": "Changed name of " + strconv.Itoa(numChanges) + " things", "success": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "Error parsing request"})
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
		usersCache = make(map[string][]string)
		userPositionCache = make(map[string]UserPositionJSON)

		c.JSON(http.StatusOK, gin.H{"message": "Changed name of " + strconv.Itoa(numChanges) + " things", "success": true})
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
				fmt.Println(string(k), v2.Username)
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
