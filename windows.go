package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Regular expression matching lines of the Wifi signals results
// of the shape `Signal : 16%`
const windowsRssiExpr = "Signal\\s*: \\d+"

var windowsRssiRegexp = regexp.MustCompile(windowsRssiExpr)

func windowsFindMac(line string) (string, bool) {
	return linuxFindMac(line)
}

func windowsFindRssi(line string) (int, bool) {
	rawRssi := windowsRssiRegexp.FindString(line)
	if rawRssi == "" {
		return 0, false
	}

	components := strings.Split(rawRssi, ": ")

	// We can safely access element 1 given that the line
	// was accepted by the rssi regexp that contains 1 `: `
	signal, err := strconv.ParseFloat(components[1], 10)

	// Convert from percentage to rssi
	signal = (signal / 2) - 100

	return int(signal), (err == nil)
}
