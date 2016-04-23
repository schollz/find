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
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime"
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

func ComputeSHA1Sum(filename string) (string, error) {
	fd, err := os.Open(filename)

	if err != nil {
		return "", err
	}
	defer fd.Close()

	h := sha1.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return "", err
	}

	result := h.Sum(nil)

	return hex.EncodeToString(result), nil
}

func DownloadVerifyFile(url string, target string, sha1sum string) error {
	needsDownload := true

	_, err := os.Stat(target)
	if err == nil {
		// File already exists, check if the hash matches
		fileSum, _ := ComputeSHA1Sum(target)
		if fileSum == sha1sum {
			needsDownload = false
		}
	}

	if needsDownload {
		fd, err := os.Create(target)
		if err != nil {
			fd.Close()
			return err
		}

		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		n, err := io.Copy(fd, response.Body)
		if err != nil {
			return err
		}

		log.Info(n, "bytes downloaded.")
		fd.Close()

		fileSum, err := ComputeSHA1Sum(target)

		if fileSum != sha1sum {
			log.Warning("3rd party utility hash (" + fileSum + ") does NOT match expected hash (" + sha1sum + ").")
			os.Remove(target)
			return errors.New("Download failed")
		}
	}

	return nil
}

func populateConfigurations(wlanInterface string) {
	osConfigurations["darwin"] = OSConfig{
		WifiScanCommand: "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -s",
		ScanConfig:      ScanParsingConfig{darwinFindMac, darwinFindRssi},
	}

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

	// Check for an alternative Windows wifi utility
	windowsConfig := OSConfig{
		WifiScanCommand: "netsh wlan show network mode=bssid",
		ScanConfig:      ScanParsingConfig{windowsFindMac, windowsFindRssi},
	}
	if runtime.GOOS == "windows" {
		_, err := os.Stat("bin")
		if os.IsNotExist(err) {
			// Create a bin folder
			os.Mkdir("bin", os.ModePerm)
		}
		needsPrompt := true
		// Check if the alternative is downloaded
		altSum := "5209b821e93d9a407b725b229223e53bb52495c9"
		_, err = os.Stat("bin/windows-wlan-util.exe")
		if err == nil {
			// Verify the file
			fileSum, _ := ComputeSHA1Sum("bin/windows-wlan-util.exe")
			if fileSum == altSum {
				windowsConfig = OSConfig{
					WifiScanCommand: "bin/windows-wlan-util.exe query",
					ScanConfig:      ScanParsingConfig{windows2FindMac, windows2FindRssi},
				}
				needsPrompt = false
			} else {
				log.Warning("3rd party utility hash (" + fileSum + ") does NOT match expected hash (" + altSum + ").")
			}
		} else {
			_, err := os.Stat("bin/windows-use-netsh")
			if err == nil {
				needsPrompt = false
			}
		}

		if needsPrompt {
			// Ask if the user wants to download the alternative
			altUtil := getInput("Do you want to download a 3rd party utility for better wifi capture? (y/n)")
			if altUtil == "y" {
				err = DownloadVerifyFile(
					"https://github.com/ScottSWu/windows-wlan-util/releases/download/v1.0/windows-wlan-util.exe",
					"bin/windows-wlan-util.exe", "5209b821e93d9a407b725b229223e53bb52495c9")
				if err == nil {
					windowsConfig = OSConfig{
						WifiScanCommand: "bin/windows-wlan-util.exe query",
						ScanConfig:      ScanParsingConfig{windows2FindMac, windows2FindRssi},
					}
				} else {
					log.Warning("Failed to download 3rd party utility. You will be asked again next time.")
				}
			} else {
				fd, _ := os.Create("bin/windows-use-netsh")
				fd.Close()
			}
		}

	}
	osConfigurations["windows"] = windowsConfig
}
