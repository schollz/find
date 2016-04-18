package main

import (
	"regexp"
	"strconv"
	"strings"
)

func processOutputLinux(out string) ([]WifiData, error) {
	data := []WifiData{}
	entry := WifiData{}
	macRegexp := regexp.MustCompile("([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})")
	rssiRegexp := regexp.MustCompile("signal: ((\\+?)|-)\\d*\\.?\\d* dBm")
	var err error

	for _, line := range strings.Split(out, "\n") {
		macAddress := macRegexp.FindString(line)
		if macAddress != "" {
			entry.Mac = macAddress
		}

		if entry.Mac == "" {
			// A mac address for the current entry has not been found yet
			continue
		}

		rawRssi := rssiRegexp.FindString(line)
		if rawRssi == "" {
			// We have a mac address but not a rssi yet
			continue
		}
		components := strings.Split(rawRssi, " ")
		var signal float64

		// We can safely access element 1 given that the line
		// was accepted by the rssi regexp that contains 2 spaces
		signal, err = strconv.ParseFloat(components[1], 10)
		if err != nil {
			continue
		}

		entry.Rssi = int(signal)

		data = append(data, entry)
		entry.Mac = ""
	}

	return data, err
}
