package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func rfLearn(group string) float64 {
	tempFile := group + ".rf.json"

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0664, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	Debug.Println("Writing " + tempFile)
	f, err := os.OpenFile(path.Join(RuntimeArgs.SourcePath, tempFile), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return -1
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if _, err = f.WriteString(string(decompressByte(v)) + "\n"); err != nil {
				panic(err)
			}
		}
		return nil
	})
	f.Close()

	// Do learning
	out, err := exec.Command("python3", "rf.py", group).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	classificationSuccess, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		panic(err)
	}
	Debug.Printf("RF classification success for '%s' is %2.2f", group, classificationSuccess)
	os.Remove(tempFile)
	return classificationSuccess
}

func rfClassify(group string, fingerprint Fingerprint) map[string]float64 {
	var m map[string]float64
	tempFile := RandomString(10) + ".json"
	d1, _ := json.Marshal(fingerprint)
	err := ioutil.WriteFile(tempFile, d1, 0644)
	if err != nil {
		return m
	}
	out, err := exec.Command("python3", "rf.py", group, tempFile).Output()
	if err != nil {
		return m
	}
	err = json.Unmarshal(out, &m)
	if err != nil {
		return m
	}
	os.Remove(tempFile)
	log.Println(m)
	return m
}
