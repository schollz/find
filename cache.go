// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// cache.go handles the global variables for caching and the clearing.

package main

import "time"

var psCache map[string]FullParameters
var usersCache map[string][]string
var userPositionCache map[string]UserPositionJSON
var isLearning map[string]bool

func init() {
	go clearCache()
}

func clearCache() {
	for {
		isLearning = make(map[string]bool)
		psCache = make(map[string]FullParameters)
		usersCache = make(map[string][]string)
		userPositionCache = make(map[string]UserPositionJSON)
		time.Sleep(time.Minute * 10)
	}
}
