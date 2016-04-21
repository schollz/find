package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Regular expression matching lines of the Wifi signals results
// of the shape `signal: -56.00 dBm`
const linuxRssiExpr = "signal: -\\d*\\.?\\d* dBm"

var rssiRegexp = regexp.MustCompile(linuxRssiExpr)

func linuxFindMac(line string) (string, bool) {
	mac := macRegexp.FindString(line)

	return strings.ToLower(mac), (mac != "")
}

func linuxFindRssi(line string) (int, bool) {
	rawRssi := rssiRegexp.FindString(line)
	if rawRssi == "" {
		return 0, false
	}

	components := strings.Split(rawRssi, " ")

	// We can safely access element 1 given that the line
	// was accepted by the rssi regexp that contains 2 spaces
	signal, err := strconv.ParseFloat(components[1], 10)
	if err != nil {
		return 0, false
	}

	return int(signal), (signal <= 0)
}
