package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func scanCommandWindows() string {
	return "bin\\windows-wlan-util.exe query"
}

func processOutputWindows(out string) []WifiData {
	w := []WifiData{}
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "Error") {
			continue
		}
		parts := strings.Fields(line)
		if len(line) < 4 {
			continue
		}
		
		val, err := strconv.ParseFloat(parts[0], 10)
		if err != nil {
			fmt.Println(line, val, err)
			continue
		}
		
		w = append(w, WifiData{Mac: parts[2], Rssi: int(val)})
	}
	
	// This is probably not the best place for this
	time.Sleep(10000 * time.Millisecond)
	
	return w
}
