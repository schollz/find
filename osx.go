package main

import (
	"fmt"
	"strconv"
	"strings"
)

func scanCommandOSX() string {
	return "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I en0"
}

func processOutputOSX(out string) []WifiData {
	w := []WifiData{}
	wTemp := WifiData{Mac: "none", Rssi: 0}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 3 {
			continue
		}
		parts := strings.Fields(line)
		if parts[0] == "BSSID" { // ???????????????????
			wTemp.Mac = parts[3]
		}
		if parts[0] == "Signal" { // ???????????????????
			val, err := strconv.ParseFloat(strings.Replace(parts[2], "%", "", 1), 10) // ???????????????????
			if err != nil {
				fmt.Println(line, val, err)
			}
			wTemp.Rssi = int(val)
			if wTemp.Mac != "none" && wTemp.Rssi != 0 {
				w = append(w, wTemp)
			}
			wTemp = WifiData{Mac: "none", Rssi: 0}
		}
	}
	return w
}
