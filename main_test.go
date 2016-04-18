package main

import (
	"encoding/json"
	"fmt"
)

func ExampleSendFingerprintServerOnlineLearn() {
	fingerprintString := []byte(`{"username": "zack", "group": "testdb", "wifi-fingerprint": [{"rssi": -43, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -56, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -60, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -66, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -69, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -79, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -84, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -83, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -86, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}, {"rssi": -74, "mac": "34:68:95:f8:25:fd"}, {"rssi": -81, "mac": "00:ac:e0:b8:ea:a0"}, {"rssi": -83, "mac": "b4:75:0e:03:cd:66"}, {"rssi": -85, "mac": "4c:60:de:fe:e5:24"}], "location": "zakhome floor 2 office", "time": 1439596537334, "password": "frusciante_0128"}`)
	var f Fingerprint
	err := json.Unmarshal(fingerprintString, &f)
	if err != nil {
		fmt.Println(err)
	}
	response, err := sendFingerprint("https://ml.internalpositioning.com/learn", f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output: Inserted 22 fingerprints for zack at zakhome floor 2 office
}

func ExampleSendFingerprintServerOnlineTrack() {
	fingerprintString := []byte(`{"username": "zack", "group": "testdb", "wifi-fingerprint": [{"rssi": -43, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -56, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -60, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -66, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -69, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -79, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -84, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -83, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -86, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}, {"rssi": -74, "mac": "34:68:95:f8:25:fd"}, {"rssi": -81, "mac": "00:ac:e0:b8:ea:a0"}, {"rssi": -83, "mac": "b4:75:0e:03:cd:66"}, {"rssi": -85, "mac": "4c:60:de:fe:e5:24"}], "location": "zakhome floor 2 office", "time": 1439596537334, "password": "frusciante_0128"}`)
	var f Fingerprint
	err := json.Unmarshal(fingerprintString, &f)
	if err != nil {
		fmt.Println(err)
	}
	response, err := sendFingerprint("https://ml.internalpositioning.com/track", f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	// Output: Calculated location: zakhome floor 2 office
}

func ExampleSendFingerprintServerOffline() {
	fingerprintString := []byte(`{"username": "zack", "group": "testdb", "wifi-fingerprint": [{"rssi": -43, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -56, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -60, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -66, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -69, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -79, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -84, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -83, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -86, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -81, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}, {"rssi": -74, "mac": "34:68:95:f8:25:fd"}, {"rssi": -81, "mac": "00:ac:e0:b8:ea:a0"}, {"rssi": -83, "mac": "b4:75:0e:03:cd:66"}, {"rssi": -85, "mac": "4c:60:de:fe:e5:24"}], "location": "zakhome floor 2 office", "time": 1439596537334, "password": "frusciante_0128"}`)
	var f Fingerprint
	err := json.Unmarshal(fingerprintString, &f)
	if err != nil {
		fmt.Println(err)
	}
	_, err = sendFingerprint("https://asldkfjaklsdjfklasjdfljaslkdjf.com", f)
	fmt.Println(err == nil)
	// Output: false
}
