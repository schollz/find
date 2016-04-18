package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/op/go-logging"
)

// Fingerprint to be sent to the server
type Fingerprint struct {
	Username        string     `json:"username"`
	Location        string     `json:"location"`
	Group           string     `json:"group"`
	Time            int64      `json:"time"`
	WifiFingerprint []WifiData `json:"wifi-fingerprint"`
}

// WifiData collected from the system
type WifiData struct {
	Mac  string `json:"mac"`
	Rssi int    `json:"rssi"`
}

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format1 = logging.MustStringFormatter(
	`%{time:15:04:05.000} %{message}`,
)
var format2 = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} - %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var verbose = true
var errorsInARow = 0

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(text))
}

func processOutput(out string, os string) (data []WifiData, err error) {
	err = nil
	data = []WifiData{}
	if os == "linux" {
		data = processOutputLinux(out)
	} else {
		err = fmt.Errorf(os + " system has no known WiFi scanning parser")
	}
	return
}

func getCommand(i string) (command string, err error) {
	err = nil
	command = ""
	if runtime.GOOS == "darwin" {
		command = scanCommandOSX()
	} else if runtime.GOOS == "linux" {
		command = scanCommandLinux(i)
	} else if runtime.GOOS == "windows" {
		command = scanCommandWindows()
	} else {
		err = fmt.Errorf(runtime.GOOS + " system has no known WiFi scanning command")
	}
	return
}

func scanWifi(i string) (string, error) {
	command, err := getCommand(i)
	if err != nil {
		return "", err
	}
	log.Info("Gathering fingerprint with '" + command + "'")
	out, err := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...).Output()
	return string(out), err
}

func setupLogging() {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stdout, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend1Formatter := logging.NewBackendFormatter(backend1, format1)
	backend2Formatter := logging.NewBackendFormatter(backend2, format2)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1Formatter)
	backend1Leveled.SetLevel(logging.CRITICAL, "")
	// Everything should be sent to backend2
	backend2Leveled := logging.AddModuleLevel(backend2Formatter)
	backend2Leveled.SetLevel(logging.INFO, "")

	// Set the backends to be used.
	if verbose {
		logging.SetBackend(backend2Leveled)
	} else {
		logging.SetBackend(backend1Leveled)
	}

	// log.Info("debug")
	// log.Info("info")
	// log.Notice("notice")
	// log.Warning("warning")
	// log.Error("err")
	// log.Critical("crit")
}

func main() {
	var f Fingerprint
	var times int
	var address string
	var wlan_interface string
	app := cli.NewApp()
	app.Name = "fingerprint"
	app.Usage = "client for sending WiFi fingerprints to a FIND server"
	app.Version = "0.2"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server,s",
			Value: "https://ml.internalpositioning.com",
			Usage: "server to connect",
		},
		cli.StringFlag{
			Name:  "group,g",
			Value: "group",
			Usage: "group name",
		},
		cli.StringFlag{
			Name:  "user,u",
			Value: "user",
			Usage: "user name",
		},
		cli.StringFlag{
			Name:  "location,l",
			Value: "location",
			Usage: "location (needed for '--learn')",
		},
		cli.IntFlag{
			Name:  "continue,c",
			Value: 3,
			Usage: "number of times to run",
		},
		cli.BoolFlag{
			Name:  "learn,e",
			Usage: "need to set if you want to learn location",
		},
		cli.BoolFlag{
			Name:  "nodebug,d",
			Usage: "turns off debugging",
		},
		cli.StringFlag{
			Name:  "interface,i",
			Value: "wlan0",
			Usage: "WiFi interface to use for scaning",
		},
	}
	app.Action = func(c *cli.Context) {
		times = c.Int("continue")
		wlan_interface = c.String("interface")
		// set group
		f.Group = strings.ToLower(c.String("group"))
		if f.Group == "group" {
			f.Group = getInput("Enter group: ")
			fmt.Println("Instead of typing next time, add '-g " + f.Group + "'")
		}
		// make sure to get a location if learning
		f.Location = strings.ToLower(c.String("location"))
		if c.Bool("learn") && c.String("location") == "location" {
			f.Location = getInput("Enter location: ")
			fmt.Println("Instead of typing next time, add '-l " + f.Location + "'")
		}
		// set server
		if c.Bool("learn") {
			address = c.String("server") + "/learn"
		} else {
			address = c.String("server") + "/track"
		}
		// set fingerprint things
		f.Time = time.Now().UnixNano() / 1000000
		f.Username = strings.ToLower(c.String("user"))
		verbose = !c.Bool("nodebug")
		setupLogging()
	}
	app.Run(os.Args)

	// Print the current parameters
	log.Notice("You can see fewer messages by adding --nodebug")
	log.Notice("User: " + f.Username)
	log.Notice("Group: " + f.Group)
	log.Notice("Server: " + address)
	if strings.Contains(address, "/learn") {
		log.Notice("Location: " + f.Location)
	}
	log.Notice("Running", times, "times (you can run more using '-c SOMENUM'). Please wait...")

	for i := 0; i < times; i++ {

		log.Info("Scanning Wifi")
		out, err := scanWifi(wlan_interface)
		if err != nil {
			if strings.Contains(err.Error(), "255") {
				fmt.Println("\nNeed to run with sudo: 'sudo ./fingerprint'")
				fmt.Println("")
				log.Fatal(string(out), err)
			} else {
				errorsInARow++
				log.Warning("Scan failed, will continue after a rest")
				time.Sleep(3000 * time.Millisecond)
				if errorsInARow > 3 {
					log.Critical("Are you sure this computer has WiFi enabled?")
					log.Fatal(string(out), err)
				} else {
					continue
				}
			}
		}
		errorsInARow = 0

		log.Info("Processing", len(strings.Split(out, "\n")), "lines out output")
		f.WifiFingerprint, err = processOutput(out, runtime.GOOS)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Sending fingerprint to " + address)
		response, err := sendFingerprint(address, f)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("Response: " + response)

	}

}

func sendFingerprint(address string, f Fingerprint) (string, error) {
	b, _ := json.Marshal(f)
	req, err := http.NewRequest("POST", address, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(body), `"success":true`) == false && strings.Contains(string(body), `"success"`) == true {
		return "", fmt.Errorf("Something wrong with server")
	}
	log.Info(string(body))

	type Response struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}
	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	if strings.Contains(r.Message, ":") {
		log.Critical(strings.TrimSpace(strings.Split(r.Message, ":")[1]))
	}
	return r.Message, nil
}
