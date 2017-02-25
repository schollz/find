package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	_, err := exec.Command("cp", []string{"data/testdb.db.backup", "data/testdb.db"}...).Output()
	if err != nil {
		log.Fatal(err)
	}
}

func TestGetStatus(t *testing.T) {
	router := gin.New()
	router.PUT("/foo", getStatus)

	req, _ := http.NewRequest("PUT", "/foo", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.Contains(resp.Body.String(), "\"success\":true"), true)
}

func TestMigrateDatabase(t *testing.T) {
	router := gin.New()
	router.PUT("/foo", migrateDatabase)

	req, _ := http.NewRequest("PUT", "/foo?from=testdb&to=newdb", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.TrimSpace(resp.Body.String()), "{\"message\":\"Successfully migrated testdb to newdb\",\"success\":true}")
	os.Remove("data/newdb.db")
}

func TestDeleteDatabase(t *testing.T) {
	CopyFile("./data/testdb.db", "./data/deleteme.db")

	router := gin.New()
	router.DELETE("/foo", deleteDatabase)

	req, _ := http.NewRequest("DELETE", "/foo?group=deleteme", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.TrimSpace(resp.Body.String()), "{\"message\":\"Successfully deleted deleteme\",\"success\":true}")
}

func TestCalculate(t *testing.T) {
	router := gin.New()
	router.GET("/foo", calculate)

	req, _ := http.NewRequest("GET", "/foo?group=testdb", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.TrimSpace(resp.Body.String()), "{\"message\":\"Parameters optimized.\",\"success\":true}")
}

func TestUserLocations(t *testing.T) {
	router := gin.New()
	router.GET("/foo", userLocations)

	req, _ := http.NewRequest("GET", "/foo?group=testdb", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.Contains(resp.Body.String(), "{\"message\":\"Correctly found\",\"success\":true,"), true)
}

func TestGetUserLocations(t *testing.T) {
	router := gin.New()
	router.GET("/foo", getUserLocations)

	req, _ := http.NewRequest("GET", "/foo", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, strings.TrimSpace(resp.Body.String()), "{\"message\":\"Error parsing request\",\"success\":false}")
}

func TestGetUserLocations2(t *testing.T) {
	router := gin.New()
	router.GET("/foo", getUserLocations)

	req, _ := http.NewRequest("GET", "/foo?group=testdb&user=zack", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, strings.Contains(resp.Body.String(), "{\"message\":\"Correctly found locations.\""), true)
}

func TestPutMixinOverrideBad(t *testing.T) {
	router := gin.New()
	router.PUT("/foo", putMixinOverride)

	req, _ := http.NewRequest("PUT", "/foo?group=testdb&mixin=100", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := `{"message":"mixin must be between 0 and 1","success":false}`
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}

func TestPutMixinOverrideGood(t *testing.T) {
	router := gin.New()
	router.PUT("/foo", putMixinOverride)

	req, _ := http.NewRequest("PUT", "/foo?group=testdb&mixin=0", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := `{"message":"Overriding mixin for testdb, now set to 0","success":true}`
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}

func TestEditNetworkName(t *testing.T) {
	router := gin.New()
	router.GET("/foo", editNetworkName)

	req, _ := http.NewRequest("GET", "/foo?group=testdb&oldname=0&newname=home", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Finished\",\"success\":true}"
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}

func TestEditName(t *testing.T) {
	router := gin.New()
	router.GET("/foo", editName)

	req, _ := http.NewRequest("GET", "/foo?group=testdb&location=zakhome%20floor%202%20office&newname=office", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Changed name of 175 things\",\"success\":true}"
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}

func TestEditUserName(t *testing.T) {
	router := gin.New()
	router.GET("/foo", editUserName)

	req, _ := http.NewRequest("GET", "/foo?group=testdb&user=zack&newname=zack2", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Changed name of 344 things\",\"success\":true}"
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}

func TestDeleteUser(t *testing.T) {
	router := gin.New()
	router.DELETE("/foo", deleteUser)

	req, _ := http.NewRequest("DELETE", "/foo?group=testdb&user=zack2@gmail.com", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	response := "{\"message\":\"Deletes 0 things  with user zack2@gmail.com\",\"success\":true}"
	assert.Equal(t, strings.TrimSpace(resp.Body.String()), response)
}
