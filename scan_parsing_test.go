package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLinuxOutEmptyResultWhenEmptyScan(t *testing.T) {
	data, _ := processOutputLinux("")
	assert.True(t, len(data) == 0, "Result must be empty")
}

func TestLinuxOutSkipWhenInvalidMac(t *testing.T) {
	out := "BSS 801:37:73:ba:f7:dc\n" +
		"signal: dBm\n" +
		"BSS 11:11:11:aa:bb:cc\n" +
		"signal: -65.00 dBm"
	data, _ := processOutputLinux(out)

	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -65},
	}

	assert.Equal(t, expected, data)
}

func TestLinuxOutSkipWhenInvalidSignal(t *testing.T) {
	out := "BSS 80:37:73:ba:f7:dc\n" +
		"signal: dBm\n" +
		"BSS 11:11:11:aa:bb:cc\n" +
		"signal: -65.00 dBm"
	data, _ := processOutputLinux(out)
	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -65},
	}

	assert.Equal(t, expected, data)
}

func TestLinuxFullOutput(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/linuxOutput.txt")
	data, err := processOutputLinux(string(dat))

	expected := []WifiData{
		WifiData{"80:37:73:ba:f7:dc", -25},
		WifiData{"80:37:73:87:46:82", -59},
		WifiData{"98:6b:3d:d7:84:e0", -60},
		WifiData{"08:95:2a:b1:e9:55", -76},
		WifiData{"2c:b0:5d:36:e3:b8", -54},
		WifiData{"58:20:b1:21:63:9f", -62},
		WifiData{"70:73:cb:bd:9f:b5", -78},
		WifiData{"b8:3e:59:78:35:99", -75},
		WifiData{"a0:63:91:2b:9e:64", -59},
		WifiData{"e0:3f:49:03:fd:38", -61},
		WifiData{"30:8d:99:71:95:c5", -78},
		WifiData{"80:37:73:ba:f7:d8", -37},
		WifiData{"a0:63:91:2b:9e:65", -55},
		WifiData{"80:37:73:87:56:36", -52},
		WifiData{"00:1a:1e:46:cd:11", -59},
		WifiData{"00:23:69:d4:47:9f", -59},
		WifiData{"54:65:de:6f:7e:d5", -70},
		WifiData{"70:73:cb:bd:9f:b6", -77},
		WifiData{"f8:35:dd:0a:da:be", -82},
		WifiData{"00:1a:1e:46:cd:10", -58},
		WifiData{"d4:05:98:57:b3:15", -79},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}

func TestPi3FullOutput(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/pi3Output.txt")
	data, err := processOutputLinux(string(dat))

	expected := []WifiData{
		{"70:73:cb:bd:9f:b5", -72},
		{"4c:60:de:fe:e5:24", -80},
		{"80:37:73:ba:f7:d8", -16},
		{"a0:63:91:2b:9e:65", -43},
		{"00:23:69:d4:47:9f", -81},
		{"80:37:73:87:56:36", -68},
		{"2c:b0:5d:36:e3:b8", -75},
		{"58:20:b1:21:63:9f", -75},
		{"30:8d:99:71:95:c5", -81},
		{"c8:b3:73:25:22:51", -85},
		{"00:1a:1e:46:cd:10", -76},
		{"e0:46:9a:6d:02:ea", -91},
		{"08:95:2a:b1:e9:55", -81},
		{"00:1d:d4:7c:bd:30", -91},
		{"8c:09:f4:d3:84:50", -90},
		{"00:1a:1e:46:cd:11", -76},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
