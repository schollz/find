package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

func scanCommandOSX() string {
	return "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -s"
}

func processOutputOSX(out string) []WifiData {
	w     := []WifiData{}
	wTemp := WifiData{Mac: "none", Rssi: 0}
	re    := regexp.MustCompile("^\\s*.+ ((?:[a-f0-9]{2}:){5}[a-f0-9]{2}) (-?\\d+)")

	for _, line := range strings.Split(out, "\n") {
		mac_signal := re.FindStringSubmatch(line)

		if mac_signal == nil {
			continue
		}

		wTemp.Mac = mac_signal[1]
		val, err := strconv.ParseFloat(mac_signal[2], 10) // strings.Replace(mac_signal[2], "%", "", 1)
		if err != nil {
			fmt.Println(line, val, err)
		}
		wTemp.Rssi = int(val)
		if wTemp.Mac != "none" && wTemp.Rssi != 0 {
			w = append(w, wTemp)
		}
		wTemp = WifiData{Mac: "none", Rssi: 0}
		}
	return w
}
