package main

// Utilities to deal with multi-platform target. In order to add
// a new target OS (as represented by the string `runtime.GOOS`)
// one must add the corresponding entry in the `osConfigurations`
// map inside the `populateConfigurations` function. One must
// associate the `runtime.GOOS` string for the target os to
// the command used in that os to retrieve a list of scanned Wifi
// APs and two functions (`FindMac` and `FindRssi`, grouped in the
// `ScanParsingConfig` that given a line will return the MAC address
// and the RSSI signal plus a boolean flag indicating if such content
// was found in the given line.

import (
	"fmt"
	"regexp"
)

// Regular expression matching standard (IEEE 802) MAC-48 addresses
const macExpr = "([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})"

var macRegexp = regexp.MustCompile(macExpr)

type OSConfig struct {
	WifiScanCommand string
	ScanConfig      ScanParsingConfig
}

type ScanParsingConfig struct {
	FindMac  func(string) (string, bool)
	FindRssi func(string) (int, bool)
}

var osConfigurations = make(map[string]OSConfig)

func GetConfiguration(os string, wlanInterface string) (OSConfig, bool) {
	if len(osConfigurations) == 0 {
		populateConfigurations(wlanInterface)
	}

	config, ok := osConfigurations[os]
	return config, ok
}

func populateConfigurations(wlanInterface string) {
	osConfigurations["darwin"] = OSConfig{
		WifiScanCommand: "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -s",
		ScanConfig:      ScanParsingConfig{darwinFindMac, darwinFindRssi},
	}

	fmt.Println(useIwlist)
	if !useIwlist {
		linuxCommand := "/sbin/iw dev " + wlanInterface + " scan -u"
		osConfigurations["linux"] = OSConfig{
			WifiScanCommand: linuxCommand,
			ScanConfig:      ScanParsingConfig{linuxFindMac, linuxFindRssi},
		}
	} else {
		linuxCommand := "/sbin/iwlist " + wlanInterface + " scan"
		osConfigurations["linux"] = OSConfig{
			WifiScanCommand: linuxCommand,
			ScanConfig:      ScanParsingConfig{linuxFindMac, linuxFindRssiIwList},
		}
	}

	osConfigurations["windows"] = OSConfig{
		WifiScanCommand: "netsh wlan show network mode=bssid",
		ScanConfig:      ScanParsingConfig{windowsFindMac, windowsFindRssi},
	}
}
