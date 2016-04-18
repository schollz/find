package main

import (
	"fmt"
	"strconv"
	"strings"
)

func scanCommandLinux(i string) string {
	return "/sbin/iw dev " + i + " scan -u"
}

func processOutputLinux(out string) []WifiData {
	w := []WifiData{}
	wTemp := WifiData{Mac: "none", Rssi: 0}
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 3 {
			continue
		}
		if line[0:3] == "BSS" {
			wTemp.Mac = strings.Fields(strings.Replace(line, "(", " (", -1))[1]
		}
		if strings.Contains(line, "signal") && strings.Contains(line, "dBm") {
			val, err := strconv.ParseFloat(strings.Fields(line)[1], 10)
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
