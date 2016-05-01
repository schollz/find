package main

import "fmt"

func ExampleProcessMQTTFingerprintTrack() {
	fmt.Println(mqttBuildFingerprint("group/track/user", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773baf7d842")))
	// Output: {group user  0 [{00:23:69:d4:47:9f -73} {80:37:73:87:56:36 -66} {00:1a:1e:46:cd:11 -78} {00:1a:1e:46:cd:10 -75} {a0:63:91:2b:9e:65 -56} {80:37:73:ba:f7:d8 -42}]} track <nil>
}

func ExampleProcessMQTTFingerprintTrackError() {
	fmt.Println(mqttBuildFingerprint("group/track/user", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773ba")))
	// Output: {group user  0 [{00:23:69:d4:47:9f -73} {80:37:73:87:56:36 -66} {00:1a:1e:46:cd:11 -78} {00:1a:1e:46:cd:10 -75} {a0:63:91:2b:9e:65 -56}]} track <nil>
}

func ExampleProcessMQTTFingerprintLearn() {
	fmt.Println(mqttBuildFingerprint("group/learn/user/location", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773baf7d842")))
	// Output: {group user location 0 [{00:23:69:d4:47:9f -73} {80:37:73:87:56:36 -66} {00:1a:1e:46:cd:11 -78} {00:1a:1e:46:cd:10 -75} {a0:63:91:2b:9e:65 -56} {80:37:73:ba:f7:d8 -42}]} learn <nil>
}
