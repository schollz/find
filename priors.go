// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// priors.go contains variables for calcualting priors.

package main

import (
	"log"
	"math"
	"path"

	"github.com/boltdb/bolt"
)

// PdfType dictates the width of gaussian smoothing
var PdfType []float32

// MaxRssi is the maximum level of signal
var MaxRssi int

// MinRssi is the minimum level of signal
var MinRssi int

// RssiPartitions are the calculated number of partitions from MinRssi and MaxRssi
var RssiPartitions int

// Absentee is the base level of probability for any signal
var Absentee float32

// RssiRange is the calculated partitions in array form
var RssiRange []float32

// FoldCrossValidation is the amount of data left out during learning to be used in cross validation
var FoldCrossValidation float64

func init() {
	PdfType = []float32{.1995, .1760, .1210, .0648, .027, 0.005}
	Absentee = 1e-6
	MinRssi = -110
	MaxRssi = 5
	RssiPartitions = MaxRssi - MinRssi + 1
	RssiRange = make([]float32, RssiPartitions)
	for i := 0; i < len(RssiRange); i++ {
		RssiRange[i] = float32(MinRssi + i)
	}
	FoldCrossValidation = 4
}

// deprecated
func optimizePriors(group string) {
	// generate the fingerprintsInMemory
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
			// fmt.Println(fingerprintsInMemory[string(k)].Location, string(k))
			fingerprintsOrdering = append(fingerprintsOrdering, string(k))
		}
		return nil
	})
	db.Close()

	var ps = *NewFullParameters()
	getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	calculatePriors(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	// fmt.Println(string(dumpParameters(ps)))
	// ps, _ = openParameters("findtest")
	var results = *NewResultsParameters()
	for n := range ps.Priors {
		ps.Results[n] = results
	}
	// fmt.Println(ps.Results)
	// ps.Priors["0"].Special["MixIn"] = 1.0
	// fmt.Println(crossValidation(group, "0", &ps))
	// fmt.Println(ps.Results)

	// loop through these parameters
	mixins := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	cutoffs := []float64{0.005}

	for n := range ps.Priors {
		bestResult := float64(0)
		bestMixin := float64(0)
		bestCutoff := float64(0)
		for _, cutoff := range cutoffs {
			for _, mixin := range mixins {
				ps.Priors[n].Special["MixIn"] = mixin
				ps.Priors[n].Special["VarabilityCutoff"] = cutoff
				avgAccuracy := crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
				// avgAccuracy := crossValidation(group, n, &ps)
				if avgAccuracy > bestResult {
					bestResult = avgAccuracy
					bestCutoff = cutoff
					bestMixin = mixin
				}
			}
		}
		ps.Priors[n].Special["MixIn"] = bestMixin
		ps.Priors[n].Special["VarabilityCutoff"] = bestCutoff
		// Final validation
		crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
		// crossValidation(group, n, &ps)
	}

	go saveParameters(group, ps)
	go setPsCache(group, ps)
}

func regenerateEverything(group string) {
	// generate the fingerprintsInMemory
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
			fingerprintsInMemory[string(v)] = loadFingerprint(v)
			fingerprintsOrdering = append(fingerprintsOrdering, string(v))
		}
		return nil
	})
	db.Close()

	var ps = *NewFullParameters()
	ps, _ = openParameters(group)
	getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	calculatePriors(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	var results = *NewResultsParameters()
	for n := range ps.Priors {
		ps.Results[n] = results
	}
	for n := range ps.Priors {
		crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
	}
	saveParameters(group, ps)
}

func crossValidation(group string, n string, ps *FullParameters, fingerprintsInMemory map[string]Fingerprint, fingerprintsOrdering []string) float64 {
	for loc := range ps.NetworkLocs[n] {
		ps.Results[n].TotalLocations[loc] = 0
		ps.Results[n].CorrectLocations[loc] = 0
		ps.Results[n].Accuracy[loc] = 0
		ps.Results[n].Guess[loc] = make(map[string]int)
	}

	it := float64(-1)
	for _, v1 := range fingerprintsOrdering {
		v2 := fingerprintsInMemory[v1]
		it++
		if math.Mod(it, FoldCrossValidation) == 0 {
			if len(v2.WifiFingerprint) == 0 {
				continue
			}
			if _, ok := ps.NetworkLocs[n][v2.Location]; ok {
				locationGuess, _ := calculatePosterior(v2, *ps)
				ps.Results[n].TotalLocations[v2.Location]++
				if locationGuess == v2.Location {
					ps.Results[n].CorrectLocations[v2.Location]++
				}
				if _, ok := ps.Results[n].Guess[v2.Location]; !ok {
					ps.Results[n].Guess[v2.Location] = make(map[string]int)
				}
				if _, ok := ps.Results[n].Guess[v2.Location][locationGuess]; !ok {
					ps.Results[n].Guess[v2.Location][locationGuess] = 0
				}
				ps.Results[n].Guess[v2.Location][locationGuess]++
			}
		}
	}

	average := float64(0)
	for loc := range ps.NetworkLocs[n] {
		if ps.Results[n].TotalLocations[loc] > 0 {
			// fmt.Println(ps.Results[n].CorrectLocations[loc], ps.Results[n].TotalLocations[loc])
			ps.Results[n].Accuracy[loc] = int(100.0 * ps.Results[n].CorrectLocations[loc] / ps.Results[n].TotalLocations[loc])
			average += float64(ps.Results[n].Accuracy[loc])
		}
	}
	average = average / float64(len(ps.NetworkLocs[n]))

	return average
}

// calculatePriors generates the prior data for Naive-Bayes classification. Now deprecated, use calculatePriorsThreaded instead.
func calculatePriors(group string, ps *FullParameters, fingerprintsInMemory map[string]Fingerprint, fingerprintsOrdering []string) {
	// defer timeTrack(time.Now(), "calculatePriors")
	ps.Priors = make(map[string]PriorParameters)
	for n := range ps.NetworkLocs {
		var newPrior = *NewPriorParameters()
		ps.Priors[n] = newPrior
	}

	// Initialization
	ps.MacVariability = make(map[string]float32)
	for n := range ps.Priors {
		ps.Priors[n].Special["MacFreqMin"] = float64(100)
		ps.Priors[n].Special["NMacFreqMin"] = float64(100)
		for loc := range ps.NetworkLocs[n] {
			ps.Priors[n].P[loc] = make(map[string][]float32)
			ps.Priors[n].NP[loc] = make(map[string][]float32)
			ps.Priors[n].MacFreq[loc] = make(map[string]float32)
			ps.Priors[n].NMacFreq[loc] = make(map[string]float32)
			for mac := range ps.NetworkMacs[n] {
				ps.Priors[n].P[loc][mac] = make([]float32, RssiPartitions)
				ps.Priors[n].NP[loc][mac] = make([]float32, RssiPartitions)
			}
		}
	}

	it := float64(-1)
	for _, v1 := range fingerprintsOrdering {
		v2 := fingerprintsInMemory[v1]
		it++
		if math.Mod(it, FoldCrossValidation) != 0 { // cross-validation
			macs := []string{}
			for _, router := range v2.WifiFingerprint {
				macs = append(macs, router.Mac)
			}

			networkName, inNetwork := hasNetwork(ps.NetworkMacs, macs)
			if inNetwork {
				for _, router := range v2.WifiFingerprint {
					if router.Rssi > MinRssi {
						ps.Priors[networkName].P[v2.Location][router.Mac][router.Rssi-MinRssi] += PdfType[0]
						for i, val := range PdfType {
							if i > 0 {
								ps.Priors[networkName].P[v2.Location][router.Mac][router.Rssi-MinRssi-i] += val
								ps.Priors[networkName].P[v2.Location][router.Mac][router.Rssi-MinRssi+i] += val
							}
						}
					} else {
						Warning.Println(router.Rssi)
					}
				}
			}

		}
	}

	// Calculate the nP
	for n := range ps.Priors {
		for locN := range ps.NetworkLocs[n] {
			for loc := range ps.NetworkLocs[n] {
				if loc != locN {
					for mac := range ps.NetworkMacs[n] {
						for i := range ps.Priors[n].P[locN][mac] {
							if ps.Priors[n].P[loc][mac][i] > 0 {
								ps.Priors[n].NP[locN][mac][i] += ps.Priors[n].P[loc][mac][i]
							}
						}
					}
				}
			}
		}
	}

	// Add in absentee, normalize P and nP and determine MacVariability
	for n := range ps.Priors {
		macAverages := make(map[string][]float32)

		for loc := range ps.NetworkLocs[n] {
			for mac := range ps.NetworkMacs[n] {
				for i := range ps.Priors[n].P[loc][mac] {
					ps.Priors[n].P[loc][mac][i] += Absentee
					ps.Priors[n].NP[loc][mac][i] += Absentee
				}
				total := float32(0)
				for _, val := range ps.Priors[n].P[loc][mac] {
					total += val
				}
				averageMac := float32(0)
				for i, val := range ps.Priors[n].P[loc][mac] {
					if val > float32(0) {
						ps.Priors[n].P[loc][mac][i] = val / total
						averageMac += RssiRange[i] * ps.Priors[n].P[loc][mac][i]
					}
				}
				if averageMac < float32(0) {
					if _, ok := macAverages[mac]; !ok {
						macAverages[mac] = []float32{}
					}
					macAverages[mac] = append(macAverages[mac], averageMac)
				}

				total = float32(0)
				for i := range ps.Priors[n].NP[loc][mac] {
					total += ps.Priors[n].NP[loc][mac][i]
				}
				if total > 0 {
					for i := range ps.Priors[n].NP[loc][mac] {
						ps.Priors[n].NP[loc][mac][i] = ps.Priors[n].NP[loc][mac][i] / total
					}
				}
			}
		}

		// Determine MacVariability
		for mac := range macAverages {
			if len(macAverages[mac]) <= 2 {
				ps.MacVariability[mac] = float32(1)
			} else {
				maxVal := float32(-10000)
				for _, val := range macAverages[mac] {
					if val > maxVal {
						maxVal = val
					}
				}
				for i, val := range macAverages[mac] {
					macAverages[mac][i] = maxVal / val
				}
				ps.MacVariability[mac] = standardDeviation(macAverages[mac])
			}
		}
	}

	// Determine mac frequencies and normalize
	for n := range ps.Priors {
		for loc := range ps.NetworkLocs[n] {
			maxCount := 0
			for mac := range ps.MacCountByLoc[loc] {
				if ps.MacCountByLoc[loc][mac] > maxCount {
					maxCount = ps.MacCountByLoc[loc][mac]
				}
			}
			for mac := range ps.MacCountByLoc[loc] {
				ps.Priors[n].MacFreq[loc][mac] = float32(ps.MacCountByLoc[loc][mac]) / float32(maxCount)
				if float64(ps.Priors[n].MacFreq[loc][mac]) < ps.Priors[n].Special["MacFreqMin"] {
					ps.Priors[n].Special["MacFreqMin"] = float64(ps.Priors[n].MacFreq[loc][mac])
				}
			}
		}
	}

	// Deteremine negative mac frequencies and normalize
	for n := range ps.Priors {
		for loc1 := range ps.Priors[n].MacFreq {
			sum := float32(0)
			for loc2 := range ps.Priors[n].MacFreq {
				if loc2 != loc1 {
					for mac := range ps.Priors[n].MacFreq[loc2] {
						ps.Priors[n].NMacFreq[loc1][mac] += ps.Priors[n].MacFreq[loc2][mac]
						sum++
					}
				}
			}
			// Normalize
			if sum > 0 {
				for mac := range ps.Priors[n].MacFreq[loc1] {
					ps.Priors[n].NMacFreq[loc1][mac] = ps.Priors[n].NMacFreq[loc1][mac] / sum
					if float64(ps.Priors[n].NMacFreq[loc1][mac]) < ps.Priors[n].Special["NMacFreqMin"] {
						ps.Priors[n].Special["NMacFreqMin"] = float64(ps.Priors[n].NMacFreq[loc1][mac])
					}
				}
			}
		}
	}

	for n := range ps.Priors {
		ps.Priors[n].Special["MixIn"] = 0.5
		ps.Priors[n].Special["VarabilityCutoff"] = 0
	}

}
