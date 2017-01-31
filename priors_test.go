package main

import (
	"log"
	"path"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func TestPriorsThreaded(t *testing.T) {
	assert.Equal(t, optimizePriorsThreaded("testdb"), nil)
}

// func ExampleTestPriors() {
// 	// optimizePriors("testdb")
// 	fmt.Println("OK")
// 	// Output: OK
// }

//go test -test.bench BenchmarkOptimizePriors -test.benchmem
func BenchmarkOptimizePriors(b *testing.B) {
	for i := 0; i < b.N; i++ {
		optimizePriors("testdb")
	}
}

func BenchmarkOptimizePriorsThreaded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		optimizePriorsThreaded("testdb")
	}
}

func BenchmarkOptimizePriorsThreadedNot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		optimizePriorsThreadedNot("testdb")
	}
}

func BenchmarkCrossValidation(b *testing.B) {
	group := "testdb"

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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for n := range ps.Priors {
			ps.Priors[n].Special["MixIn"] = 0.5
			ps.Priors[n].Special["VarabilityCutoff"] = 0.005
			crossValidation(group, n, &ps, fingerprintsInMemory, fingerprintsOrdering)
			break
		}
	}
}

func BenchmarkCalculatePriors(b *testing.B) {
	group := "testdb"
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
	getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calculatePriors(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	}
}
