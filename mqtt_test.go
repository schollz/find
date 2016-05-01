package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessMQTTFingerprintTrack(t *testing.T) {
	fingerprint, _, _ := mqttBuildFingerprint("group/track/user", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773baf7d842"))
	json, _ := json.Marshal(fingerprint)
	assert.Equal(t, string(json), "{\"group\":\"group\",\"username\":\"user\",\"location\":\"\",\"timestamp\":0,\"wifi-fingerprint\":[{\"mac\":\"00:23:69:d4:47:9f\",\"rssi\":-73},{\"mac\":\"80:37:73:87:56:36\",\"rssi\":-66},{\"mac\":\"00:1a:1e:46:cd:11\",\"rssi\":-78},{\"mac\":\"00:1a:1e:46:cd:10\",\"rssi\":-75},{\"mac\":\"a0:63:91:2b:9e:65\",\"rssi\":-56},{\"mac\":\"80:37:73:ba:f7:d8\",\"rssi\":-42}]}")
}

func TestProcessMQTTFingerprintTrackError(t *testing.T) {
	fingerprint, _, _ := mqttBuildFingerprint("group/track/user", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773ba"))
	json, _ := json.Marshal(fingerprint)
	assert.Equal(t, string(json), "{\"group\":\"group\",\"username\":\"user\",\"location\":\"\",\"timestamp\":0,\"wifi-fingerprint\":[{\"mac\":\"00:23:69:d4:47:9f\",\"rssi\":-73},{\"mac\":\"80:37:73:87:56:36\",\"rssi\":-66},{\"mac\":\"00:1a:1e:46:cd:11\",\"rssi\":-78},{\"mac\":\"00:1a:1e:46:cd:10\",\"rssi\":-75},{\"mac\":\"a0:63:91:2b:9e:65\",\"rssi\":-56}]}")
}

func TestProcessMQTTFingerprintLearn(t *testing.T) {
	fingerprint, _, _ := mqttBuildFingerprint("group/learn/user/location", []byte("002369d4479f7380377387563666001a1e46cd1178001a1e46cd1075a063912b9e6556803773baf7d842"))
	json, _ := json.Marshal(fingerprint)
	assert.Equal(t, string(json), "{\"group\":\"group\",\"username\":\"user\",\"location\":\"location\",\"timestamp\":0,\"wifi-fingerprint\":[{\"mac\":\"00:23:69:d4:47:9f\",\"rssi\":-73},{\"mac\":\"80:37:73:87:56:36\",\"rssi\":-66},{\"mac\":\"00:1a:1e:46:cd:11\",\"rssi\":-78},{\"mac\":\"00:1a:1e:46:cd:10\",\"rssi\":-75},{\"mac\":\"a0:63:91:2b:9e:65\",\"rssi\":-56},{\"mac\":\"80:37:73:ba:f7:d8\",\"rssi\":-42}]}")
}
