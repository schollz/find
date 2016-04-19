package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Regular expression matching lines of the Wifi signals results
// of the shape `22:86:8c:d5:30:d8 -82`
const darwinRssiExpr = macExpr + " -\\d+"

var darwinRssiRegexp = regexp.MustCompile(darwinRssiExpr)

func darwinFindMac(line string) (string, bool) {
	return linuxFindMac(line)
}

func darwinFindRssi(line string) (int, bool) {
	rawRssi := darwinRssiRegexp.FindString(line)
	if rawRssi == "" {
		return 0, false
	}

	components := strings.Split(rawRssi, " ")

	// We can safely access element 1 given that the line
	// was accepted by the rssi regexp that contains 1 space
	signal, err := strconv.Atoi(components[1])

	return signal, (err == nil)
}
