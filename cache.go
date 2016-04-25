package main

import "time"

var psCache map[string]FullParameters
var usersCache map[string][]string
var userPositionCache map[string]UserPositionJSON
var isLearning map[string]bool
var mixinOverrideCache map[string]float64

func init() {
	go clearCache()
}

func clearCache() {
	for {
		isLearning = make(map[string]bool)
		psCache = make(map[string]FullParameters)
		usersCache = make(map[string][]string)
		userPositionCache = make(map[string]UserPositionJSON)
		mixinOverrideCache = make(map[string]float64)
		time.Sleep(time.Minute * 10)
	}
}
