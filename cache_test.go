package main

import (
	"fmt"
	"path"
	"testing"

	"github.com/boltdb/bolt"
)

// BenchmarkCache needs to have precomputed parameters for testdb (run Optimize after loading testdb.sh)
func BenchmarkGetPSCache(b *testing.B) {
	var ps FullParameters
	db, err := bolt.Open(path.Join("data", "testdb.db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
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
	db.Close()
	setPsCache("testdb", ps)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getPsCache("testdb")
	}

}

// BenchmarkCache needs to have precomputed parameters for testdb (run Optimize after loading testdb.sh)
func BenchmarkSetPSCache(b *testing.B) {
	var ps FullParameters
	db, err := bolt.Open(path.Join("data", "testdb.db"), 0600, nil)
	if err != nil {
		Error.Println(err)
	}
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
	db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setPsCache("testdb", ps)
	}

}
