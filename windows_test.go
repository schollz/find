package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var windowsConfig = ScanParsingConfig{windowsFindMac, windowsFindRssi}

func TestWindowsOutEmptyResultWhenEmptyScan(t *testing.T) {
	data, _ := ParseOutput(windowsConfig, "")
	assert.True(t, len(data) == 0, "Result must be empty")
}

func TestWindowsOutSkipWhenInvalidMac(t *testing.T) {
	out := "BSSID 1                 : cu0:c1:c0:f0:6f:cd\n" +
		"Signal             : 16%\n" +
		"BSSID 1                 : 11:11:11:aa:bb:cc\n" +
		"Signal             : 38%"

	data, _ := ParseOutput(windowsConfig, out)

	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -81},
	}

	assert.Equal(t, expected, data)
}

func TestWindowsOutSkipWhenInvalidSignal(t *testing.T) {
	out := "BSSID 1                 : c0:c1:c0:f0:6f:cd\n" +
		"Signal             : \n" +
		"BSSID 1                 : 11:11:11:aa:bb:cc\n" +
		"Signal             : 38%"

	data, _ := ParseOutput(windowsConfig, out)
	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -81},
	}

	assert.Equal(t, expected, data)
}

func TestWindowsFullOutput(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/windowsOutput.txt")
	data, err := ParseOutput(windowsConfig, string(dat))

	expected := []WifiData{
		WifiData{"c0:c1:c0:f0:6f:cd", -92},
		WifiData{"f8:35:dd:0a:da:be", -81},
		WifiData{"00:1a:1e:46:cd:11", -80},
		WifiData{"00:1a:1e:46:cd:10", -80},
		WifiData{"2c:b0:5d:36:e3:b8", -81},
		WifiData{"58:20:b1:21:63:9f", -62},
		WifiData{"98:6b:3d:d7:84:e0", -91},
		WifiData{"80:37:73:87:56:36", -82},
		WifiData{"00:23:69:d4:47:9f", -82},
		WifiData{"80:37:73:ba:f7:d8", -55},
		WifiData{"a0:63:91:2b:9e:65", -58},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
