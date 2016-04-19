package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var darwinConfig = ScanParsingConfig{darwinFindMac, darwinFindRssi}

func TestDarwinOutEmptyResultWhenEmptyScan(t *testing.T) {
	data, _ := ParseOutput(darwinConfig, "")
	assert.True(t, len(data) == 0, "Result must be empty")
}

func TestDarwinOutSkipWhenInvalidMac(t *testing.T) {
	out := "XFSETUP-30DA 1u0:86:8c:d5:30:d8 -81  11" +
		"      Y  US WPA(PSK/AES,TKIP/TKIP) WPA2(PSK/AES,TKIP/TKIP)\n" +
		"traviata 11:11:11:aa:bb:cc -65  11,-1   Y  -- WPA2(PSK/AES/AES)"

	data, _ := ParseOutput(darwinConfig, out)

	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -65},
	}

	assert.Equal(t, expected, data)
}

func TestDarwinOutSkipWhenInvalidSignal(t *testing.T) {
	out := "XFSETUP-30DA 10:86:8c:d5:30:d8 - 11" +
		"      Y  US WPA(PSK/AES,TKIP/TKIP) WPA2(PSK/AES,TKIP/TKIP)\n" +
		"traviata 11:11:11:aa:bb:cc -65  11,-1   Y  -- WPA2(PSK/AES/AES)"

	data, _ := ParseOutput(darwinConfig, out)
	expected := []WifiData{
		WifiData{"11:11:11:aa:bb:cc", -65},
	}

	assert.Equal(t, expected, data)
}

func TestDarwinFullOutput(t *testing.T) {
	dat, _ := ioutil.ReadFile("test/osxOutput.txt")
	data, err := ParseOutput(darwinConfig, string(dat))

	expected := []WifiData{
		WifiData{"22:86:8c:d5:30:d8", -82},
		WifiData{"10:86:8c:d5:30:d8", -81},
		WifiData{"00:7f:28:8b:0c:1d", -74},
		WifiData{"74:44:01:35:46:34", -75},
		WifiData{"20:e5:2a:16:79:d4", -73},
		WifiData{"4e:7a:8a:1d:7e:cc", -70},
		WifiData{"e6:89:2c:1a:02:e0", -91},
		WifiData{"e8:89:2c:1a:02:e0", -92},
		WifiData{"3c:7a:8a:1d:7e:cc", -70},
		WifiData{"50:6a:03:93:89:e3", -91},
		WifiData{"08:86:3b:6d:bc:16", -53},
		WifiData{"c0:ff:d4:e6:dd:2b", -67},
		WifiData{"60:a4:4c:29:8b:44", -88},
		WifiData{"08:86:3b:6d:bc:18", -66},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
