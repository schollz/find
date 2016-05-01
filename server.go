// Copyright 2015-2016 Zack Scholl. All rights reserved.
// Use of this source code is governed by a AGPL
// license that can be found in the LICENSE file.

// server.go handles Flag parsing and starts the Gin-Tonic webserver.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RuntimeArgs contains all runtime
// arguments available
var RuntimeArgs struct {
	ExternalIP        string
	Port              string
	ServerCRT         string
	ServerKey         string
	SourcePath        string
	Socket            string
	Cwd               string
	MqttServer        string
	MqttAdmin         string
	MosquittoPID      string
	MqttAdminPassword string
	Mqtt              bool
}

// VersionNum keeps track of the version
var VersionNum string

// init initiates the paths in RuntimeArgs
func init() {
	cwd, _ := os.Getwd()
	RuntimeArgs.SourcePath = path.Join(cwd, "data")
	RuntimeArgs.Cwd = cwd
}

func main() {
	VersionNum = "2.0"
	// _, executableFile, _, _ := runtime.Caller(0) // get full path of this file

	// Bing flags for changing parameters of FIND
	flag.StringVar(&RuntimeArgs.Port, "p", ":8003", "port to bind")
	flag.StringVar(&RuntimeArgs.Socket, "s", "", "unix socket")
	flag.StringVar(&RuntimeArgs.ServerCRT, "crt", "", "location of ssl crt")
	flag.StringVar(&RuntimeArgs.ServerKey, "key", "", "location of ssl key")
	flag.StringVar(&RuntimeArgs.MqttServer, "mqtt", "", "turn on MQTT message passing")
	flag.StringVar(&RuntimeArgs.MqttAdmin, "mqttadmin", "", "admin to read all messages")
	flag.StringVar(&RuntimeArgs.MqttAdminPassword, "mqttadminpass", "", "admin to read all messages")
	flag.StringVar(&RuntimeArgs.MosquittoPID, "mosquitto", "", "mosquitto PID")
	flag.CommandLine.Usage = func() {
		fmt.Println(`find (version ` + VersionNum + `)
run this to start the server and then visit localhost at the port you specify
(see parameters).
Example: 'find yourserver.com'
Example: 'find -p :8080 localhost:8080'
Example: 'find -s /var/run/find.sock'
Example: 'find -db /var/lib/find/db.bolt localhost:8003'
Example: 'find -p :8080 -crt ssl/server.crt -key ssl/server.key localhost:8080'
Options:`)
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()
	RuntimeArgs.ExternalIP = flag.Arg(0)
	if RuntimeArgs.ExternalIP == "" {
		RuntimeArgs.ExternalIP = GetLocalIP() + RuntimeArgs.Port
	}

	// Check whether all the MQTT variables are passed to initiate the MQTT routines
	if len(RuntimeArgs.MqttServer) > 0 && len(RuntimeArgs.MqttAdmin) > 0 && len(RuntimeArgs.MosquittoPID) > 0 {
		RuntimeArgs.Mqtt = true
		setupMqtt()
	} else {
		RuntimeArgs.Mqtt = false
	}

	// Setup Gin-Gonic
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Load templates
	r.LoadHTMLGlob(path.Join(RuntimeArgs.Cwd, "templates/*"))

	// Load static files (if they are not hosted by external service)
	r.Static("static/", path.Join(RuntimeArgs.Cwd, "static/"))

	// Create cookie store to keep track of logged in user
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 404-page redirects to login
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"ErrorMessage": "Please login first.",
		})
	})

	// Routes for logging in and viewing dashboards (routes.go)
	r.GET("/", slash)
	r.GET("/login", slashLogin)
	r.POST("/login", slashLoginPOST)
	r.GET("/logout", slashLogout)
	r.GET("/dashboard/:group", slashDashboard)
	r.GET("/explore/:group/:network/:location", slashExplore2)
	r.GET("/pie/:group/:network/:location", slashPie)

	// Routes for performing fingerprinting (fingerprint.go)
	r.POST("/learn", learnFingerprintPOST)
	r.POST("/track", trackFingerprintPOST)

	// Routes for MQTT (mqtt.go)
	r.PUT("/mqtt", putMQTT)

	// Routes for API access (api.go)
	r.GET("/location", getUserLocations)
	r.GET("/editname", editName)
	r.GET("/editusername", editUserName)
	r.GET("/editnetworkname", editNetworkName)
	r.DELETE("/location", deleteLocation)
	r.DELETE("/locations", deleteLocations)
	r.DELETE("/user", deleteUser)
	r.GET("/calculate", calculate)
	r.GET("/status", getStatus)
	r.GET("/userlocs", userLocations) // to be deprecated
	r.GET("/whereami", whereAmI)      // to be deprecated
	r.PUT("/mixin", putMixinOverride)

	// Load and display the logo
	dat, _ := ioutil.ReadFile("./static/logo.txt")
	fmt.Println(string(dat))

	// Check whether user is providing certificates
	if RuntimeArgs.Socket != "" {
		r.RunUnix(RuntimeArgs.Socket)
	} else if RuntimeArgs.ServerCRT != "" && RuntimeArgs.ServerKey != "" {
		fmt.Println("(version " + VersionNum + ") is up and running on https://" + RuntimeArgs.ExternalIP)
		fmt.Println("-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----")
		r.RunTLS(RuntimeArgs.Port, RuntimeArgs.ServerCRT, RuntimeArgs.ServerKey)
	} else {
		fmt.Println("(version " + VersionNum + ") is up and running on http://" + RuntimeArgs.ExternalIP)
		fmt.Println("-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----")
		r.Run(RuntimeArgs.Port)
	}
}
