package main

import (
	"fmt"
	"log"
	"path"
	"testing"

	"github.com/boltdb/bolt"
)

func BenchmarkLoadParameters(b *testing.B) {
	var ps FullParameters = *NewFullParameters()
	db, err := bolt.Open(path.Join("data", "testdb.db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
	defer db.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
		if err != nil {
			Error.Println(err)
		}

	}
}

func BenchmarkGetParameters(b *testing.B) {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getParameters(group, &ps, fingerprintsInMemory, fingerprintsOrdering)
	}
}
