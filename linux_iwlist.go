package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Regular expression matching lines of the Wifi signals results
// of the shape `level=-56.00 dBm`
const linuxIwListRssiExpr = "level=-\\d*\\.?\\d* dBm"

var rssiRegexpIwlist = regexp.MustCompile(linuxIwListRssiExpr)

func linuxFindRssiIwList(line string) (int, bool) {
	rawRssi := rssiRegexpIwlist.FindString(line)
	if rawRssi == "" {
		return 0, false
	}

	rawRssi = strings.Replace(rawRssi, "level=", "", 1)
	rawRssi = strings.Replace(rawRssi, " dBm", "", 1)

	// We can safely access element 1 given that the line
	// was accepted by the rssi regexp that contains 2 spaces
	signal, err := strconv.ParseFloat(rawRssi, 10)
	if err != nil {
		return 0, false
	}

	return int(signal), (signal <= 0)
}
