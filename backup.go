// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// backup.go contains functions for dumping a backup database.

package main

import (
	"log"
	"os"
	"path"

	"github.com/boltdb/bolt"
)

func dumpFingerprints(group string) error {
	err := os.MkdirAll("dump-"+group, 0664)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(path.Join("dump-"+group, "learning.fingerprints"), os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}

	defer f.Close()

	db, err := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0664, nil)
	if err != nil {
		log.Fatal(err)
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
	db.Close()
	return nil
}
