package main

type OSConfig struct {
	WifiScanCommand string
	ProcessOutput   func(string) ([]WifiData, error)
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
		ProcessOutput:   processOutputDarwin,
	}

	linuxCommand := "/sbin/iw dev " + wlanInterface + " scan -u"
	osConfigurations["linux"] = OSConfig{
		WifiScanCommand: linuxCommand,
		ProcessOutput:   processOutputLinux,
	}

	osConfigurations["windows"] = OSConfig{
		WifiScanCommand: "netsh wlan show network mode=bssid",
		ProcessOutput:   processOutputWindows,
	}
}
