package main

import (
	"strconv"
	"strings"
	"time"
)

// windows2 uses a third party alternative to netsh
// See https://github.com/ScottSWu/windows-wlan-util for details

func windows2FindMac(line string) (string, bool) {
	return linuxFindMac(line)
}

func windows2FindRssi(line string) (int, bool) {
	line = strings.Trim(line, " ")
	space := strings.Index(line, " ")
	
	if space == -1 {
		return 0, false
	}
	
	rssi, err := strconv.ParseFloat(line[:space], 10)
	
	// This is probably not the best place for this
	// Rough timings suggest that finding ~30 access points takes around 10
	// seconds. Thus for each MAC/Rssi value found, we should wait about 300 ms.
	time.Sleep(300 * time.Millisecond)
	
	return int(rssi), (err == nil)
}
