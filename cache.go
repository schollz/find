// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// cache.go handles the global variables for caching and the clearing.

package main

import (
	"sync"
	"time"
)

var counter = struct {
	sync.RWMutex
	ps           map[string]FullParameters
	users        map[string][]string
	userPosition map[string]UserPositionJSON
	isLearning   map[string]bool
}{isLearning: make(map[string]bool),
	ps:           make(map[string]FullParameters),
	users:        make(map[string][]string),
	userPosition: make(map[string]UserPositionJSON),
}

func init() {
	go clearCache()
}

func clearCache() {
	for {
		time.Sleep(time.Minute * 10)
	}
}
