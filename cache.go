// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// cache.go handles the global variables for caching and the clearing.

package main

import (
	"strings"
	"sync"
	"time"
)

var psCache = struct {
	sync.RWMutex
	m map[string]FullParameters
}{m: make(map[string]FullParameters)}

var usersCache = struct {
	sync.RWMutex
	m map[string][]string
}{m: make(map[string][]string)}

var userPositionCache = struct {
	sync.RWMutex
	m map[string]UserPositionJSON
}{m: make(map[string]UserPositionJSON)}

var isLearning = struct {
	sync.RWMutex
	m map[string]bool
}{m: make(map[string]bool)}

func init() {
	go clearCache()
	go clearCacheFast()
}

func clearCacheFast() {
	for {
		go resetCache("userCache")
		time.Sleep(time.Second * 30)
	}
}

func clearCache() {
	for {
		//Debug.Println("Resetting cache")
		go resetCache("isLearning")
		go resetCache("psCache")
		go resetCache("userPositionCache")
		time.Sleep(time.Second * 600)
	}
}

func resetCache(cache string) {
	if cache == "userCache" {
		usersCache.Lock()
		usersCache.m = make(map[string][]string)
		usersCache.Unlock()
	} else if cache == "userPositionCache" {
		userPositionCache.Lock()
		userPositionCache.m = make(map[string]UserPositionJSON)
		userPositionCache.Unlock()
	} else if cache == "psCache" {
		psCache.Lock()
		psCache.m = make(map[string]FullParameters)
		psCache.Unlock()
	} else if cache == "isLearning" {
		isLearning.Lock()
		isLearning.m = make(map[string]bool)
		isLearning.Unlock()
	}
}

func getLearningCache(group string) (bool, bool) {
	//Debug.Println("getLearningCache")
	isLearning.RLock()
	cached, ok := isLearning.m[group]
	isLearning.RUnlock()
	return cached, ok
}

func setLearningCache(group string, val bool) {
	isLearning.Lock()
	isLearning.m[group] = val
	isLearning.Unlock()
}

func getUserCache(group string) ([]string, bool) {
	//Debug.Println("Getting userCache")
	usersCache.RLock()
	cached, ok := usersCache.m[group]
	usersCache.RUnlock()
	return cached, ok
}

func setUserCache(group string, users []string) {
	usersCache.Lock()
	usersCache.m[group] = users
	usersCache.Unlock()
}

func appendUserCache(group string, user string) {
	usersCache.Lock()
	if _, ok := usersCache.m[group]; ok {
		if len(usersCache.m[group]) == 0 {
			usersCache.m[group] = append([]string{}, strings.ToLower(user))
		}
	}
	usersCache.Unlock()
}

func getPsCache(group string) (FullParameters, bool) {
	//Debug.Println("Getting pscache")
	psCache.RLock()
	psCached, ok := psCache.m[group]
	psCache.RUnlock()
	return psCached, ok
}

func setPsCache(group string, ps FullParameters) {
	//Debug.Println("Setting pscache")
	psCache.Lock()
	psCache.m[group] = ps
	psCache.Unlock()
	return
}

func getUserPositionCache(group string) (UserPositionJSON, bool) {
	//Debug.Println("getUserPositionCache")
	userPositionCache.RLock()
	cached, ok := userPositionCache.m[group]
	userPositionCache.RUnlock()
	return cached, ok
}

func setUserPositionCache(group string, p UserPositionJSON) {
	//Debug.Println("setUserPositionCache")
	userPositionCache.Lock()
	userPositionCache.m[group] = p
	userPositionCache.Unlock()
	return
}
