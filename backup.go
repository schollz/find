// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// backup.go contains functions for dumping a backup database.

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/boltdb/bolt"
)

func dumpFingerprints(group string) error {
	err := os.MkdirAll("dump-"+group, 0664)
	if err != nil {
		return err
	}

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0664, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Dump the learning fingerprints
	f, err := os.OpenFile(path.Join("dump-"+group, "learning"), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
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

	// Dump the tracking fingerprints
	f, err = os.OpenFile(path.Join("dump-"+group, "tracking"), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints-track"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if _, err = f.WriteString(string(decompressByte(v)) + "\n"); err != nil {
				panic(err)
			}
		}
		return nil
	})
	f.Close()

	return nil
}

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

func dumpFingerprintsSVM(group string) error {
	err := os.MkdirAll("dump-"+group, 0664)
	if err != nil {
		return err
	}

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0755, nil)
	if err != nil {
		log.Fatal(err)
	}

	macs := make(map[string]int)
	locations := make(map[string]int)
	macI := 1
	locationI := 1
	// Dump the learning fingerprints
	if err != nil {
		return err
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v2 := loadFingerprint(v)
			for _, fingerprint := range v2.WifiFingerprint {
				if _, ok := macs[fingerprint.Mac]; !ok {
					macs[fingerprint.Mac] = macI
					macI++
				}
			}
			if _, ok := locations[v2.Location]; !ok {
				locations[v2.Location] = locationI
				locationI++
			}
		}
		return nil
	})

	fmt.Println(locations)
	fmt.Println(macs)
	// Dump the tracking fingerprints
	f, err := os.OpenFile("svm", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v2 := loadFingerprint(v)
			_, err := f.WriteString(strconv.Itoa(locations[v2.Location]) + " ")
			if err != nil {
				panic(err)
			}

			// To create a map as input
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
				f.WriteString(strconv.Itoa(k) + ":" + strconv.Itoa(m[k]) + " ")
			}

			f.WriteString("\n")
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	f.Close()
	db.Close()
	return nil
}
