package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
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

func getCommand() (command string, err error) {
	err = nil
	command = ""
	if runtime.GOOS == "darwin" {
		command = scanCommandOSX()
	} else if runtime.GOOS == "linux" {
		command = scanCommandLinux()
	} else if runtime.GOOS == "windows" {
		command = scanCommandWindows()
	} else {
		err = fmt.Errorf(runtime.GOOS + " system has no known WiFi scanning command")
	}
	return
}

func scanWifi() (string, error) {
	command, err := getCommand()
	if err != nil {
		return "", err
	}
	log.Println("Gathering fingerprint with '" + command + "'")
	out, err := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...).Output()
	return string(out), err
}

func main() {
	var f Fingerprint
	var times int
	var address string
	app := cli.NewApp()
	app.Name = "fingerprint"
	app.Usage = "client for sending WiFi fingerprints to a FIND server"
	app.Version = "0.1"
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
	}
	app.Action = func(c *cli.Context) {
		times = c.Int("continue")
		// set group
		f.Group = strings.ToLower(c.String("group"))
		if f.Group == "group" {
			f.Group = getInput("Enter group: ")
			fmt.Println("Next time use './fingerprint -g " + f.Group + "'")
		}
		// make sure to get a location if learning
		f.Location = strings.ToLower(c.String("location"))
		if c.Bool("learn") && c.String("location") == "location" {
			f.Location = getInput("Enter location: ")
			fmt.Println("Next time use './fingerprint -g " + f.Group + " -l " + f.Location + "'")
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
	}
	app.Run(os.Args)

	for i := 0; i < times; i++ {

		log.Println("Scanning Wifi")
		out, err := scanWifi()
		if err != nil {
			if strings.Contains(err.Error(), "255") {
				fmt.Println("\nNeed to run with sudo: \n\nsudo ./fingerprint")
			}
			log.Fatal(string(out), err)
		}

		log.Println("Processing ", len(out), " lines out output")
		f.WifiFingerprint, err = processOutput(out, runtime.GOOS)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Sending fingerprint to " + address)
		sendFingerprint(address, f)
		if err != nil {
			log.Fatal(err)
		}

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

	return string(body), nil
}
