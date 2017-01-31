// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// priorsThreaded.go contains the main Prior-calculation function which is multi-threaded

package main

import (
	"fmt"
	"log"
	"math"
	"path"
	"runtime"

	"github.com/boltdb/bolt"
)

// following this:https://play.golang.org/p/hK2h-irKyz
type resultA struct {
	mixin         float64
	locationGuess string
	locationTrue  string
	n             string
}

type jobA struct {
	mixin        float64
	locs         []string
	bayes1       []float64
	bayes2       []float64
	n            string
	locationTrue string
}

// MaxParallelism returns the maximum parallelism https://stackoverflow.com/questions/13234749/golang-how-to-verify-number-of-processors-on-which-a-go-program-is-running
func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func worker(id int, jobs <-chan jobA, results chan<- resultA) {
	for j := range jobs {
		maxVal := float64(-1)
		locationGuess := ""
		for i, loc := range j.locs {
			PBayesMix := j.mixin*j.bayes1[i] + (1-j.mixin)*j.bayes2[i]
			if PBayesMix > maxVal {
				maxVal = PBayesMix
				locationGuess = loc
			}
		}
		results <- resultA{locationGuess: locationGuess,
			locationTrue: j.locationTrue,
			mixin:        j.mixin,
			n:            j.n}
	}
}

// optimizePriorsThreaded generates the optimized prior data for Naive-Bayes classification.
func optimizePriorsThreaded(group string) error {
	// Debug.Println("Optimizing priors for " + group)
	// generate the fingerprintsInMemory
	fingerprintsInMemory := make(map[string]Fingerprint)
	var fingerprintsOrdering []string
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		if b == nil {
			return fmt.Errorf("No fingerprint bucket")
		}
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fingerprintsInMemory[string(k)] = loadFingerprint(v)
			fingerprintsOrdering = append(fingerprintsOrdering, string(k))
		}
		return nil
	})
	db.Close()
	if err != nil {
		return err
	}

	var ps = *NewFullParameters()
	getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	calculatePriors(group, &ps, fingerprintsInMemory, fingerprintsOrdering)

	var results = *NewResultsParameters()
	for n := range ps.Priors {
		ps.Results[n] = results
	}

	// loop through these parameters
	mixins := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	mixinOverride, _ := getMixinOverride(group)
	if mixinOverride >= 0 && mixinOverride <= 1 {
		mixins = []float64{mixinOverride}
	}

	// Choose cutoff
	cutoffs := []float64{0.005, 0.05, 0.1}
	cutoffOverride, _ := getCutoffOverride(group)
	if cutoffOverride >= 0 && cutoffOverride <= 1 {
		cutoffs = []float64{cutoffOverride}
	}

	bestMixin := make(map[string]float64)
	bestResult := make(map[string]float64)
	bestCutoff := make(map[string]float64)
	for n := range ps.Priors {
		bestResult[n] = 0
		bestMixin[n] = 0
		bestCutoff[n] = 0
	}

	for _, cutoff := range cutoffs {

		//                 network      id      loc    value
		PBayes1 := make(map[string]map[string]map[string]float64)
		PBayes2 := make(map[string]map[string]map[string]float64)
		totalJobs := 0
		for n := range ps.Priors {
			it := float64(-1)
			PBayes1[n] = make(map[string]map[string]float64)
			PBayes2[n] = make(map[string]map[string]float64)
			PBayes1[n] = make(map[string]map[string]float64)
			PBayes2[n] = make(map[string]map[string]float64)
			for _, v1 := range fingerprintsOrdering {
				it++
				if math.Mod(it, FoldCrossValidation) != 0 {
					_, ok := ps.NetworkLocs[n][fingerprintsInMemory[v1].Location]
					if len(fingerprintsInMemory[v1].WifiFingerprint) == 0 || !ok {
						continue
					}
					totalJobs++
					PBayes1[n][v1], PBayes2[n][v1] = calculatePosteriorThreadSafe(fingerprintsInMemory[v1], ps, cutoff)
				}
			}
		}

		numJobs := len(mixins) * totalJobs
		runtime.GOMAXPROCS(MaxParallelism())
		chanJobs := make(chan jobA, 1+numJobs)
		chanResults := make(chan resultA, 1+numJobs)
		for w := 1; w <= MaxParallelism(); w++ {
			go worker(w, chanJobs, chanResults)
		}

		finalResults := make(map[string]map[float64]ResultsParameters)
		for n := range ps.Priors {
			finalResults[n] = make(map[float64]ResultsParameters)
			for _, mixin := range mixins {

				finalResults[n][mixin] = *NewResultsParameters()
				for loc := range ps.NetworkLocs[n] {
					finalResults[n][mixin].TotalLocations[loc] = 0
					finalResults[n][mixin].CorrectLocations[loc] = 0
					finalResults[n][mixin].Accuracy[loc] = 0
					finalResults[n][mixin].Guess[loc] = make(map[string]int)
				}
				// Loop through each fingerprint
				for id := range PBayes1[n] {
					locs := []string{}
					bayes1 := []float64{}
					bayes2 := []float64{}
					for key := range PBayes1[n][id] {
						locs = append(locs, key)
						bayes1 = append(bayes1, PBayes1[n][id][key])
						bayes2 = append(bayes2, PBayes2[n][id][key])
					}
					trueLoc := fingerprintsInMemory[id].Location
					chanJobs <- jobA{n: n,
						mixin:        mixin,
						locs:         locs,
						locationTrue: trueLoc,
						bayes1:       bayes1,
						bayes2:       bayes2}
				}
			}
		}
		close(chanJobs)

		for a := 1; a <= numJobs; a++ {
			t := <-chanResults
			finalResults[t.n][t.mixin].TotalLocations[t.locationTrue]++
			if t.locationGuess == t.locationTrue {
				finalResults[t.n][t.mixin].CorrectLocations[t.locationTrue]++
			}
			if _, ok := finalResults[t.n][t.mixin].Guess[t.locationTrue]; !ok {
				finalResults[t.n][t.mixin].Guess[t.locationTrue] = make(map[string]int)
			}
			if _, ok := finalResults[t.n][t.mixin].Guess[t.locationTrue][t.locationGuess]; !ok {
				finalResults[t.n][t.mixin].Guess[t.locationTrue][t.locationGuess] = 0
			}
			finalResults[t.n][t.mixin].Guess[t.locationTrue][t.locationGuess]++
		}

		for n := range ps.Priors {
			for mixin := range finalResults[n] {
				average := float64(0)
				it := 0
				for loc := range finalResults[n][mixin].TotalLocations {
					if finalResults[n][mixin].TotalLocations[loc] > 0 {
						finalResults[n][mixin].Accuracy[loc] = int(100.0 * finalResults[n][mixin].CorrectLocations[loc] / finalResults[n][mixin].TotalLocations[loc])
						// Debug.Println(n, mixin, cutoff, loc, finalResults[n][mixin].Accuracy[loc])
						average += float64(finalResults[n][mixin].Accuracy[loc])
						it++
					}
				}
				average = average / float64(it)
				// fmt.Println(mixin, average)
				if average > bestResult[n] {
					bestResult[n] = average
					bestMixin[n] = mixin
					bestCutoff[n] = cutoff
				}
			}
		}

	}

	// Load new priors and calculate new cross Validation
	for n := range ps.Priors {
		ps.Priors[n].Special["MixIn"] = bestMixin[n]
		ps.Priors[n].Special["VarabilityCutoff"] = bestCutoff[n]
		crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
	}

	// Debug.Println(getUsers(group))
	go resetCache("usersCache")
	go saveParameters(group, ps)
	go setPsCache(group, ps)

	return nil
}

func optimizePriorsThreadedNot(group string) {
	// generate the fingerprintsInMemory
	// Debug.Println("Optimizing priors for " + group)
	fingerprintsInMemory := make(map[string]Fingerprint)
	var fingerprintsOrdering []string
	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("fingerprints"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fingerprintsInMemory[string(k)] = loadFingerprint(v)
			fingerprintsOrdering = append(fingerprintsOrdering, string(k))
		}
		return nil
	})
	db.Close()

	var ps = *NewFullParameters()
	getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	calculatePriors(group, &ps, fingerprintsInMemory, fingerprintsOrdering)

	var results = *NewResultsParameters()
	for n := range ps.Priors {
		ps.Results[n] = results
	}

	// loop through these parameters
	mixins := []float64{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
	// cutoff := 0.1
	cutoffs := []float64{0.005, 0.05, 0.1}
	bestMixin := make(map[string]float64)
	bestResult := make(map[string]float64)
	bestCutoff := make(map[string]float64)
	for n := range ps.Priors {
		bestResult[n] = 0
		bestMixin[n] = 0
		bestCutoff[n] = 0
	}

	for _, cutoff := range cutoffs {

		//                 network      id      loc    value
		PBayes1 := make(map[string]map[string]map[string]float64)
		PBayes2 := make(map[string]map[string]map[string]float64)
		totalJobs := 0
		for n := range ps.Priors {
			it := float64(-1)
			PBayes1[n] = make(map[string]map[string]float64)
			PBayes2[n] = make(map[string]map[string]float64)
			PBayes1[n] = make(map[string]map[string]float64)
			PBayes2[n] = make(map[string]map[string]float64)
			for _, v1 := range fingerprintsOrdering {
				it++
				if math.Mod(it, FoldCrossValidation) != 0 {
					_, ok := ps.NetworkLocs[n][fingerprintsInMemory[v1].Location]
					if len(fingerprintsInMemory[v1].WifiFingerprint) == 0 || !ok {
						continue
					}
					totalJobs++
					PBayes1[n][v1], PBayes2[n][v1] = calculatePosteriorThreadSafe(fingerprintsInMemory[v1], ps, cutoff)
				}
			}
		}

		finalResults := make(map[string]map[float64]ResultsParameters)
		bestMixin := make(map[string]float64)
		bestResult := make(map[string]float64)
		for n := range ps.Priors {
			bestResult[n] = 0
			bestMixin[n] = 0
			finalResults[n] = make(map[float64]ResultsParameters)
			for _, mixin := range mixins {

				finalResults[n][mixin] = *NewResultsParameters()
				for loc := range ps.NetworkLocs[n] {
					finalResults[n][mixin].TotalLocations[loc] = 0
					finalResults[n][mixin].CorrectLocations[loc] = 0
					finalResults[n][mixin].Accuracy[loc] = 0
					finalResults[n][mixin].Guess[loc] = make(map[string]int)
				}
				// Loop through each fingerprint
				for id := range PBayes1[n] {
					maxVal := float64(-1)
					locationGuess := ""
					for key := range PBayes1[n][id] {
						PBayesMix := mixin*PBayes1[n][id][key] + (1-mixin)*PBayes2[n][id][key]
						if PBayesMix > maxVal {
							maxVal = PBayesMix
							locationGuess = key
						}
						locationTrue := fingerprintsInMemory[id].Location
						finalResults[n][mixin].TotalLocations[locationTrue]++
						if locationGuess == locationTrue {
							finalResults[n][mixin].CorrectLocations[locationTrue]++
						}
						if _, ok := finalResults[n][mixin].Guess[locationTrue]; !ok {
							finalResults[n][mixin].Guess[locationTrue] = make(map[string]int)
						}
						if _, ok := finalResults[n][mixin].Guess[locationTrue][locationGuess]; !ok {
							finalResults[n][mixin].Guess[locationTrue][locationGuess] = 0
						}
						finalResults[n][mixin].Guess[locationTrue][locationGuess]++
					}
				}
				average := float64(0)
				it := 0
				for loc := range finalResults[n][mixin].TotalLocations {
					if finalResults[n][mixin].TotalLocations[loc] > 0 {
						finalResults[n][mixin].Accuracy[loc] = int(100.0 * finalResults[n][mixin].CorrectLocations[loc] / finalResults[n][mixin].TotalLocations[loc])
						average += float64(finalResults[n][mixin].Accuracy[loc])
						it++
					}
				}
				average = average / float64(it)
				// fmt.Println(mixin, average, a)
				if average > bestResult[n] {
					bestResult[n] = average
					bestMixin[n] = mixin
					bestCutoff[n] = cutoff
				}
			}
		}
	}

	// Load new priors and calculate new cross Validation
	for n := range ps.Priors {
		ps.Priors[n].Special["MixIn"] = bestMixin[n]
		ps.Priors[n].Special["VarabilityCutoff"] = bestCutoff[n]
		crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
	}
	go saveParameters(group, ps)
	go setPsCache(group, ps)
	// Debug.Println("Analyzed ", totalJobs, " fingerprints")
}
