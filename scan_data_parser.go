package main

import (
	"strings"
)

// Core parser for Wifi APs scanned results. This will process
// the output line by line calling `FindMac` and `FindRssi` for
// the given configuration.
func ParseOutput(config ScanParsingConfig, out string) ([]WifiData, error) {
	data := []WifiData{}
	entry := WifiData{}
	var err error

	for _, line := range strings.Split(out, "\n") {
		macAddress, found := config.FindMac(line)
		if found {
			entry.Mac = macAddress
		}

		if entry.Mac == "" {
			// A mac address for the current entry has not been found yet
			continue
		}

		rawRssi, valid := config.FindRssi(line)
		if !valid {
			// We have a mac address but not a rssi yet
			continue
		}
		entry.Rssi = rawRssi

		data = append(data, entry)
		entry.Mac = ""
	}

	return data, err
}
