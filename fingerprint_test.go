package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	Init(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard)
}

func TestTrackFingerprintRoute(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`

	router := gin.New()
	router.POST("/foo", trackFingerprintPOST)

	req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(jsonTest))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.Contains(resp.Body.String(), "Current location: zakhome floor 2 office"), true)
}

func BenchmarkTrackFingerprintRoute(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "testdb", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`

	router := gin.New()
	router.POST("/foo", trackFingerprintPOST)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(jsonTest))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

func BenchmarkLearnFingerprintRoute(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "testdb", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`

	router := gin.New()
	router.POST("/foo", learnFingerprintPOST)

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("POST", "/foo", bytes.NewBufferString(jsonTest))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}

// Example1 measures compressed vs uncompressed size
func TestCompression(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	bFingerprint := dumpFingerprint(res)
	loadedFingerprint := loadFingerprint(bFingerprint)
	dumped, _ := json.Marshal(loadedFingerprint)
	bFingerprint2 := dumpFingerprint(loadedFingerprint)
	isGood := 100*len(dumped)/len(bFingerprint2) == 244 || 100*len(dumped)/len(bFingerprint2) == 246
	assert.Equal(t, isGood, true)
}

func BenchmarkPutFingerprintInDatabase(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "testdbfoo", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	os.Remove(path.Join(RuntimeArgs.SourcePath, "testdbfoo.db"))
	putFingerprintIntoDatabase(res, "fingerprints") // create it for the first time
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		putFingerprintIntoDatabase(res, "fingerprints")
	}
}

func BenchmarkGetFingerprintInDatabase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		group := "testdb"
		db, _ := bolt.Open(path.Join(RuntimeArgs.SourcePath, group+".db"), 0600, nil)
		db.View(func(tx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := tx.Bucket([]byte("fingerprints"))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				loadFingerprint(v)
				break
			}
			return nil
		})
		db.Close()
	}
}

func BenchmarkLoadFingerprint(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	bFingerprint, _ := res.MarshalJSON()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := Fingerprint{}
		res.UnmarshalJSON(bFingerprint)
	}
}

func BenchmarkLoadCompressedFingerprint(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	bFingerprintUncompressed, _ := res.MarshalJSON()
	bFingerprint := compressByte(bFingerprintUncompressed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res := Fingerprint{}
		decompressed := decompressByte(bFingerprint)
		res.UnmarshalJSON(decompressed)
	}
}

func BenchmarkDumpFingerprint(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res.MarshalJSON()
	}
}

func BenchmarkDumpCompressedFingerprint(b *testing.B) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dumped, _ := res.MarshalJSON()
		compressByte(dumped)
	}
}

func TestLearnFingerprintPOST(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`

	router := gin.New()
	router.POST("/foo", learnFingerprintPOST)

	req, _ := http.NewRequest("POST", "/foo", bytes.NewBuffer([]byte(jsonTest)))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Inserted fingerprint containing 18 APs for zack (find) at zakhome floor 2 office\",\"success\":true}"
	assert.Equal(t, response, strings.TrimSpace(resp.Body.String()))
}

func TestLearnFingerprint(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	message, _ := learnFingerprint(res)
	assert.Equal(t, strings.TrimSpace(message), "Inserted fingerprint containing 18 APs for zack (find) at zakhome floor 2 office")
}

func TestTrackFingerprintPOST(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`

	router := gin.New()
	router.POST("/foo", learnFingerprintPOST)

	req, _ := http.NewRequest("POST", "/foo", bytes.NewBuffer([]byte(jsonTest)))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Inserted fingerprint containing 18 APs for zack (find) at zakhome floor 2 office\",\"success\":true}"
	assert.Equal(t, response, strings.TrimSpace(resp.Body.String()))
}

func TestTrackFingerprintFunction(t *testing.T) {
	jsonTest := `{"username": "zack", "group": "Find", "wifi-fingerprint": [{"rssi": -45, "mac": "80:37:73:ba:f7:d8"}, {"rssi": -58, "mac": "80:37:73:ba:f7:dc"}, {"rssi": -61, "mac": "a0:63:91:2b:9e:65"}, {"rssi": -68, "mac": "a0:63:91:2b:9e:64"}, {"rssi": -70, "mac": "70:73:cb:bd:9f:b5"}, {"rssi": -75, "mac": "d4:05:98:57:b3:10"}, {"rssi": -75, "mac": "00:23:69:d4:47:9f"}, {"rssi": -76, "mac": "30:46:9a:a0:28:c4"}, {"rssi": -81, "mac": "2c:b0:5d:36:e3:b8"}, {"rssi": -82, "mac": "00:1a:1e:46:cd:10"}, {"rssi": -82, "mac": "20:aa:4b:b8:31:c8"}, {"rssi": -83, "mac": "e8:ed:05:55:21:10"}, {"rssi": -83, "mac": "ec:1a:59:4a:9c:ed"}, {"rssi": -88, "mac": "b8:3e:59:78:35:99"}, {"rssi": -84, "mac": "e0:46:9a:6d:02:ea"}, {"rssi": -84, "mac": "00:1a:1e:46:cd:11"}, {"rssi": -84, "mac": "f8:35:dd:0a:da:be"}, {"rssi": -84, "mac": "b4:75:0e:03:cd:69"}], "location": "zakhome floor 2 office", "time": 1439596533831, "password": "frusciante_0128"}`
	res := Fingerprint{}
	json.Unmarshal([]byte(jsonTest), &res)
	message, _, _, _, _, _ := trackFingerprint(res)
	assert.Equal(t, strings.TrimSpace(message), "Current location: zakhome floor 2 office")
}

//
// func BenchmarkParseDate(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ParseDate("02/17/2016 3:05pm")
// 	}
// }
