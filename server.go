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
	ExternalIP string
	Port       string
	ServerCRT  string
	ServerKey  string
	SourcePath string
	Socket     string
	Cwd        string
}

// VersionNum keeps track of the version
var VersionNum string

func init() {
	cwd, _ := os.Getwd()
	RuntimeArgs.SourcePath = path.Join(cwd, "data")
	RuntimeArgs.Cwd = cwd
}

func main() {
	VersionNum = "2.0"
	// _, executableFile, _, _ := runtime.Caller(0) // get full path of this file
	flag.StringVar(&RuntimeArgs.Port, "p", ":8003", "port to bind")
	flag.StringVar(&RuntimeArgs.Socket, "s", "", "unix socket")
	flag.StringVar(&RuntimeArgs.ServerCRT, "crt", "", "location of ssl crt")
	flag.StringVar(&RuntimeArgs.ServerKey, "key", "", "location of ssl key")
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

	// var ps FullParameters = *NewFullParameters()
	// getParameters("findtest2", &ps)
	// calculatePriors("findtest2", &ps)
	// saveParameters("findtest2", ps)
	// // fmt.Println(string(dumpParameters(ps)))
	// saveParameters("findtest2", ps)
	// fmt.Println(ps.MacVariability)
	// fmt.Println(ps.NetworkLocs)
	// optimizePriors("findtest2")
	// ps, _ = openParameters("findtest2")
	//
	// getPositionBreakdown("findtest2", "zack")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob(path.Join(RuntimeArgs.Cwd, "templates/*"))
	r.Static("static/", path.Join(RuntimeArgs.Cwd, "static/"))
	store := sessions.NewCookieStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// 404 page
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"ErrorMessage": "Please login first.",
		})
	})

	// webpages (routes.go)
	r.GET("/pie/:group/:network/:location", slashPie)
	r.GET("/", slash)
	r.GET("/login", slashLogin)
	r.POST("/login", slashLoginPOST)
	r.GET("/logout", slashLogout)
	r.GET("/dashboard/:group", slashDashboard)
	r.GET("/location/:group/:user", slashLocation)
	r.GET("/explore/:group/:network/:location", slashExplore2)

	// fingerprinting (fingerprint.go)
	r.POST("/fingerprint", handleFingerprint)
	r.POST("/learn", handleFingerprint)
	r.POST("/track", trackFingerprint)

	// API routes (api.go)
	r.GET("/whereami", whereAmI)
	r.GET("/editname", editName)
	r.GET("/editusername", editUserName)
	r.GET("/editnetworkname", editNetworkName)
	r.DELETE("/location", deleteName)
	r.DELETE("/user", deleteUser)
	r.GET("/calculate", calculate)
	r.GET("/userlocs", userLocations)
	r.GET("/status", getStatus)
	dat, _ := ioutil.ReadFile("./static/logo.txt")
	fmt.Println(string(dat))
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
