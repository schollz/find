// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// parameters.go contains structures and functions for setting and getting Naive-Bayes parameters.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
)

// PersistentParameters are not reloaded each time
type PersistentParameters struct {
	NetworkRenamed map[string][]string
}

// PriorParameters contains the network-specific bayesian priors and Mac frequency, as well as special variables
type PriorParameters struct {
	P        map[string]map[string][]float32 // standard P
	NP       map[string]map[string][]float32 // standard nP
	MacFreq  map[string]map[string]float32   // Frequency of a mac in a certain location
	NMacFreq map[string]map[string]float32   // Frequency of a mac, in everywhere BUT a certain location
	Special  map[string]float64
}

// ResultsParameters contains the information about the accuracy from crossValidation
type ResultsParameters struct {
	Accuracy         map[string]int            // accuracy measurement for a given location
	TotalLocations   map[string]int            // number of locations
	CorrectLocations map[string]int            // number of times guessed correctly
	Guess            map[string]map[string]int // correct -> guess -> times
}

// FullParameters is the full parameter set for a given group
type FullParameters struct {
	NetworkMacs    map[string]map[string]bool // map of networks and then the associated macs in each
	NetworkLocs    map[string]map[string]bool // map of the networks, and then the associated locations in each
	MacVariability map[string]float32         // variability of macs
	MacCount       map[string]int             // number of each mac
	MacCountByLoc  map[string]map[string]int  // number of each mac, by location
	UniqueLocs     []string
	UniqueMacs     []string
	Priors         map[string]PriorParameters   // generate priors for each network
	Results        map[string]ResultsParameters // generate priors for each network
	Loaded         bool                         // flag to determine if parameters have been loaded
}

// NewFullParameters generates a blank FullParameters
func NewFullParameters() *FullParameters {
	return &FullParameters{
		NetworkMacs:    make(map[string]map[string]bool),
		NetworkLocs:    make(map[string]map[string]bool),
		MacCount:       make(map[string]int),
		MacCountByLoc:  make(map[string]map[string]int),
		UniqueMacs:     []string{},
		UniqueLocs:     []string{},
		Priors:         make(map[string]PriorParameters),
		MacVariability: make(map[string]float32),
		Results:        make(map[string]ResultsParameters),
		Loaded:         false,
	}
}

// NewPriorParameters generates a blank PriorParameters
func NewPriorParameters() *PriorParameters {
	return &PriorParameters{
		P:        make(map[string]map[string][]float32),
		NP:       make(map[string]map[string][]float32),
		MacFreq:  make(map[string]map[string]float32),
		NMacFreq: make(map[string]map[string]float32),
		Special:  make(map[string]float64),
	}
}

// NewResultsParameters generates a blank ResultsParameters
func NewResultsParameters() *ResultsParameters {
	return &ResultsParameters{
		Accuracy:         make(map[string]int),
		TotalLocations:   make(map[string]int),
		CorrectLocations: make(map[string]int),
		Guess:            make(map[string]map[string]int),
	}
}

// NewPersistentParameters returns the peristent parameters initialization
func NewPersistentParameters() *PersistentParameters {
	return &PersistentParameters{
		NetworkRenamed: make(map[string][]string),
	}
}

func dumpParameters(res FullParameters) []byte {
	jsonByte, _ := res.MarshalJSON()
	return compressByte(jsonByte)
}

func loadParameters(jsonByte []byte) FullParameters {
	var res2 FullParameters
	res2.UnmarshalJSON(decompressByte(jsonByte))
	return res2
}

func saveParameters(group string, res FullParameters) error {
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err2 := tx.CreateBucketIfNotExists([]byte("resources"))
		if err2 != nil {
			return fmt.Errorf("create bucket: %s", err2)
		}

		err2 = bucket.Put([]byte("fullParameters"), dumpParameters(res))
		if err2 != nil {
			return fmt.Errorf("could add to bucket: %s", err2)
		}
		return err2
	})
	return err
}

func openParameters(group string) (FullParameters, error) {
	psCached, ok := getPsCache(group)
	if ok {
		return psCached, nil
	}

	var ps = *NewFullParameters()
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("fullParameters"))
		ps = loadParameters(v)
		return nil
	})

	go setPsCache(group, ps)
	return ps, err
}

func openPersistentParameters(group string) (PersistentParameters, error) {
	var persistentPs = *NewPersistentParameters()
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("persistentParameters"))
		json.Unmarshal(v, &persistentPs)
		return nil
	})
	return persistentPs, err
}

func savePersistentParameters(group string, res PersistentParameters) error {
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err2 := tx.CreateBucketIfNotExists([]byte("resources"))
		if err2 != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		jsonByte, _ := json.Marshal(res)
		err2 = bucket.Put([]byte("persistentParameters"), jsonByte)
		if err2 != nil {
			return fmt.Errorf("could add to bucket: %s", err)
		}
		return err2
	})
	Debug.Println("Saved")
	return err
}

func getParameters(group string, ps *FullParameters, fingerprintsInMemory map[string]Fingerprint, fingerprintsOrdering []string) {
	persistentPs, err := openPersistentParameters(group)
	ps.NetworkMacs = make(map[string]map[string]bool)
	ps.NetworkLocs = make(map[string]map[string]bool)
	ps.UniqueMacs = []string{}
	ps.UniqueLocs = []string{}
	ps.MacCount = make(map[string]int)
	ps.MacCountByLoc = make(map[string]map[string]int)
	ps.Loaded = true
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get all parameters that don't need a network graph
	for _, v1 := range fingerprintsOrdering {
		v2 := fingerprintsInMemory[v1]

		// unique locs
		if !stringInSlice(v2.Location, ps.UniqueLocs) {
			ps.UniqueLocs = append(ps.UniqueLocs, v2.Location)
		}

		// mac by location count
		if _, ok := ps.MacCountByLoc[v2.Location]; !ok {
			ps.MacCountByLoc[v2.Location] = make(map[string]int)
		}

		// building network
		macs := []string{}

		for _, router := range v2.WifiFingerprint {
			// building network
			macs = append(macs, router.Mac)

			// unique macs
			if !stringInSlice(router.Mac, ps.UniqueMacs) {
				ps.UniqueMacs = append(ps.UniqueMacs, router.Mac)
			}

			// mac count
			if _, ok := ps.MacCount[router.Mac]; !ok {
				ps.MacCount[router.Mac] = 0
			}
			ps.MacCount[router.Mac]++

			// mac by location count
			if _, ok := ps.MacCountByLoc[v2.Location][router.Mac]; !ok {
				ps.MacCountByLoc[v2.Location][router.Mac] = 0
			}
			ps.MacCountByLoc[v2.Location][router.Mac]++
		}

		// building network
		ps.NetworkMacs = buildNetwork(ps.NetworkMacs, macs)
	}

	ps.NetworkMacs = mergeNetwork(ps.NetworkMacs)

	// Rename the NetworkMacs
	if len(persistentPs.NetworkRenamed) > 0 {
		newNames := []string{}
		for k := range persistentPs.NetworkRenamed {
			newNames = append(newNames, k)
		}
		for n := range ps.NetworkMacs {
			renamed := false
			for mac := range ps.NetworkMacs[n] {
				for renamedN := range persistentPs.NetworkRenamed {
					if stringInSlice(mac, persistentPs.NetworkRenamed[renamedN]) && !stringInSlice(n, newNames) {
						ps.NetworkMacs[renamedN] = make(map[string]bool)
						for k, v := range ps.NetworkMacs[n] {
							ps.NetworkMacs[renamedN][k] = v
						}
						delete(ps.NetworkMacs, n)
						renamed = true
					}
					if renamed {
						break
					}
				}
				if renamed {
					break
				}
			}
		}
	}

	// Get the locations for each graph (Has to have network built first)
	for _, v1 := range fingerprintsOrdering {
		v2 := fingerprintsInMemory[v1]
		macs := []string{}
		for _, router := range v2.WifiFingerprint {
			macs = append(macs, router.Mac)
		}
		networkName, inNetwork := hasNetwork(ps.NetworkMacs, macs)
		if inNetwork {
			if _, ok := ps.NetworkLocs[networkName]; !ok {
				ps.NetworkLocs[networkName] = make(map[string]bool)
			}
			if _, ok := ps.NetworkLocs[networkName][v2.Location]; !ok {
				ps.NetworkLocs[networkName][v2.Location] = true
			}
		}
	}

}

func getMixinOverride(group string) (float64, error) {
	group = strings.ToLower(group)
	override := float64(-1)
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	defer db.Close()
	if err != nil {
		Error.Println(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("mixinOverride"))
		if len(v) == 0 {
			return fmt.Errorf("No mixin override")
		}
		override, err = strconv.ParseFloat(string(v), 64)
		return err
	})
	return override, err
}

func getCutoffOverride(group string) (float64, error) {
	group = strings.ToLower(group)
	override := float64(-1)
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	defer db.Close()
	if err != nil {
		Error.Println(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("resources"))
		if b == nil {
			return fmt.Errorf("Resources dont exist")
		}
		v := b.Get([]byte("cutoffOverride"))
		if len(v) == 0 {
			return fmt.Errorf("No mixin override")
		}
		override, err = strconv.ParseFloat(string(v), 64)
		return err
	})
	return override, err
}

func setMixinOverride(group string, mixin float64) error {
	if (mixin < 0 || mixin > 1) && mixin != -1 {
		return fmt.Errorf("mixin must be between 0 and 1")
	}
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	defer db.Close()
	if err != nil {
		Error.Println(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err2 := tx.CreateBucketIfNotExists([]byte("resources"))
		if err2 != nil {
			return fmt.Errorf("create bucket: %s", err2)
		}

		err2 = bucket.Put([]byte("mixinOverride"), []byte(strconv.FormatFloat(mixin, 'E', -1, 64)))
		if err2 != nil {
			return fmt.Errorf("could add to bucket: %s", err2)
		}
		return err2
	})
	return err
}

func setCutoffOverride(group string, cutoff float64) error {
	if (cutoff < 0 || cutoff > 1) && cutoff != -1 {
		return fmt.Errorf("cutoff must be between 0 and 1")
	}
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	defer db.Close()
	if err != nil {
		Error.Println(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err2 := tx.CreateBucketIfNotExists([]byte("resources"))
		if err2 != nil {
			return fmt.Errorf("create bucket: %s", err2)
		}

		err2 = bucket.Put([]byte("cutoffOverride"), []byte(strconv.FormatFloat(cutoff, 'E', -1, 64)))
		if err2 != nil {
			return fmt.Errorf("could add to bucket: %s", err2)
		}
		return err2
	})
	return err
}
