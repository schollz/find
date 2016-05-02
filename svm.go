package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// # sudo apt-get install g++
// # wget http://www.csie.ntu.edu.tw/~cjlin/cgi-bin/libsvm.cgi?+http://www.csie.ntu.edu.tw/~cjlin/libsvm+tar.gz
// # tar -xvf libsvm-3.18.tar.gz
// # cd libsvm-3.18
// # make
//
// cp ~/Documents/find/svm ./
// cat svm | shuf > svm.shuffled
// ./svm-scale -l 0 -u 1 svm.shuffled > svm.shuffled.scaled
// head -n 500 svm.shuffled.scaled > learning
// tail -n 1500 svm.shuffled.scaled > testing
// ./svm-train -s 0 -t 0 -b 1 learning > /dev/null
// ./svm-predict -b 1 testing learning.model out

type Svm struct {
	Data     string
	Mac      map[string]string
	Location map[string]string
}

func dumpFingerprintsSVM(group string) error {
	macs := make(map[string]int)
	locations := make(map[string]int)
	macsFromID := make(map[string]string)
	locationsFromID := make(map[string]string)
	macI := 1
	locationI := 1

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0755, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v2 := loadFingerprint(v)
			for _, fingerprint := range v2.WifiFingerprint {
				if _, ok := macs[fingerprint.Mac]; !ok {
					macs[fingerprint.Mac] = macI
					macsFromID[strconv.Itoa(macI)] = fingerprint.Mac
					macI++
				}
			}
			if _, ok := locations[v2.Location]; !ok {
				locations[v2.Location] = locationI
				locationsFromID[strconv.Itoa(locationI)] = v2.Location
				locationI++
			}
		}
		return nil
	})

	svmData := ""
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v2 := loadFingerprint(v)
			svmData = svmData + makeSVMLine(v2, macs, locations)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("resources"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		err = bucket.Put([]byte("svmData"), []byte(svmData))
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}

		s, _ := json.Marshal(macsFromID)
		err = bucket.Put([]byte("macsFromID"), s)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}

		s, _ = json.Marshal(locationsFromID)
		err = bucket.Put([]byte("locationsFromID"), s)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}

		s, _ = json.Marshal(macs)
		err = bucket.Put([]byte("macs"), s)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}

		s, _ = json.Marshal(locations)
		err = bucket.Put([]byte("locations"), s)
		if err != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}

		return err
	})

	return err
}

func calculateSVM(group string) error {
	defer timeTrack(time.Now(), "TIMEING")
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0755, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	svmData := ""
	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("svmData"))
		svmData = string(v)
		return err
	})
	if err != nil {
		panic(err)
	}
	if len(svmData) == 0 {
		return fmt.Errorf("No data")
	}

	lines := strings.Split(svmData, "\n")
	list := rand.Perm(len(lines))
	learningSet := ""
	testingSet := ""
	for i, _ := range list {
		if len(lines[list[i]]) == 0 {
			continue
		}
		if i < len(list)/2 {
			learningSet = learningSet + lines[list[i]] + "\n"
		} else {
			testingSet = testingSet + lines[list[i]] + "\n"
		}
	}

	tempFileTrain := RandStringBytesMaskImprSrc(6) + ".learning"
	tempFileTest := RandStringBytesMaskImprSrc(6) + ".testing"
	tempFileOut := RandStringBytesMaskImprSrc(6) + ".out"
	d1 := []byte(learningSet)
	err = ioutil.WriteFile(tempFileTrain, d1, 0644)
	if err != nil {
		panic(err)
	}

	d1 = []byte(testingSet)
	err = ioutil.WriteFile(tempFileTest, d1, 0644)
	if err != nil {
		panic(err)
	}

	cmd := "svm-scale"
	args := "-l 0 -u 1 " + tempFileTrain
	Debug.Println(cmd, args)
	outCmd, err := exec.Command(cmd, strings.Split(args, " ")...).Output()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(tempFileTrain+".scaled", outCmd, 0644)
	if err != nil {
		panic(err)
	}

	cmd = "svm-scale"
	args = "-l 0 -u 1 " + tempFileTest
	Debug.Println(cmd, args)
	outCmd, err = exec.Command(cmd, strings.Split(args, " ")...).Output()
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(tempFileTest+".scaled", outCmd, 0644)
	if err != nil {
		panic(err)
	}

	cmd = "svm-train"
	args = "-s 0 -t 0 -b 1 " + tempFileTrain + ".scaled data/" + group + ".model"
	Debug.Println(cmd, args)
	if _, err = exec.Command(cmd, strings.Split(args, " ")...).Output(); err != nil {
		panic(err)
	}

	cmd = "svm-predict"
	args = "-b 1 " + tempFileTest + ".scaled data/" + group + ".model " + tempFileOut
	Debug.Println(cmd, args)
	outCmd, err = exec.Command(cmd, strings.Split(args, " ")...).Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(outCmd))

	os.Remove(tempFileTrain)
	os.Remove(tempFileTrain + ".scaled")
	os.Remove(tempFileTest)
	os.Remove(tempFileTest + ".scaled")
	os.Remove(tempFileOut)
	return nil
}

func classify(jsonFingerprint Fingerprint) {
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, jsonFingerprint.Group+".db"), 0755, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var locations map[string]int
	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("locations"))
		json.Unmarshal(v, &locations)
		return err
	})
	if err != nil {
		panic(err)
	}

}

func makeSVMLine(v2 Fingerprint, macs map[string]int, locations map[string]int) string {
	svmData := ""
	svmData = svmData + strconv.Itoa(locations[v2.Location]) + " "

	m := make(map[int]int)
	for _, fingerprint := range v2.WifiFingerprint {
		m[macs[fingerprint.Mac]] = fingerprint.Rssi
	}
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		svmData = svmData + strconv.Itoa(k) + ":" + strconv.Itoa(m[k]) + " "
	}
	return svmData + "\n"
}

// cp ~/Documents/find/svm ./
// cat svm | shuf > svm.shuffled
// ./svm-scale -l 0 -u 1 svm.shuffled > svm.shuffled.scaled
// head -n 500 svm.shuffled.scaled > learning
// tail -n 1500 svm.shuffled.scaled > testing
// ./svm-train -s 0 -t 0 -b 1 learning > /dev/null
// ./svm-predict -b 1 testing learning.model out
