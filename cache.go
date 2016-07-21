// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// cache.go handles the global variables for caching and the clearing.

package main

import (
	"sync"
	"time"
)

var psCache = struct {
	sync.RWMutex
	m map[string]FullParameters
}{m: make(map[string]FullParameters)}

var usersCache map[string][]string
var userPositionCache map[string]UserPositionJSON
var isLearning map[string]bool

func init() {
	go clearCache()
}

func clearCache() {
	for {
		Debug.Println("Resetting cache")
		isLearning = make(map[string]bool)
		psCache.Lock()
		psCache.m = make(map[string]FullParameters)
		psCache.Unlock()
		usersCache = make(map[string][]string)
		userPositionCache = make(map[string]UserPositionJSON)
		time.Sleep(time.Second * 10)
	}
}

func getPsCache(group string) (FullParameters, bool) {
	Debug.Println("Getting pscache")
	psCache.RLock()
	psCached, ok := psCache.m[group]
	psCache.RUnlock()
	return psCached, ok
}

func setPsCache(group string, ps FullParameters) {
	Debug.Println("Setting pscache")
	psCache.Lock()
	psCache.m[group] = ps
	psCache.Unlock()
	return
}
