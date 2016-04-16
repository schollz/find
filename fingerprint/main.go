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
	"strconv"
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

func processOutput(out []byte) ([]WifiData, error) {
	w := []WifiData{}
	wTemp := WifiData{Mac: "none", Rssi: 0}
	if runtime.GOOS == "linux" {

		for _, line := range strings.Split(string(out), "\n") {
			if len(line) < 3 {
				continue
			}
			if line[0:3] == "BSS" {
				wTemp.Mac = strings.Split(line, " ")[1]
			}
			if strings.Contains(line, "signal") && strings.Contains(line, "dBm") {
				val, _ := strconv.ParseFloat(strings.Split(line, " ")[1], 10)
				wTemp.Rssi = int(val)
				if wTemp.Mac != "none" && wTemp.Rssi != 0 {
					w = append(w, wTemp)
				}
				wTemp = WifiData{Mac: "none", Rssi: 0}
			}
		}
	} else if runtime.GOOS == "darwin" {
		for _, line := range strings.Split(string(out), "\n") {
			// PLEASE HELP
			fmt.Println(line) // <- this should be parsed to fill out []WifiData{}
		}

	} else if runtime.GOOS == "windows" {
		for _, line := range strings.Split(string(out), "\n") {
			// PLEASE HELP
			fmt.Println(line) // <- this should be parsed to fill out []WifiData{}
		}

	}
	if len(w) == 0 {
		return w, fmt.Errorf("This operating system is no supported for processing output")
	} else {
		return w, nil
	}
}

func getCommand() (string, error) {
	log.Println("Detected OS: " + runtime.GOOS)
	if runtime.GOOS == "darwin" {
		return "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I en0", nil
	} else if runtime.GOOS == "linux" {
		return `/sbin/iw dev wlan0 scan -u`, nil
	} else if runtime.GOOS == "windows" {
		return "netsh wlan show network mode=bssid", nil
	}
	return "none", fmt.Errorf("This operating system is not supported for getting WiFi")
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
	command, err := getCommand()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < times; i++ {
		log.Println("Gathering fingerprint with '" + command + "'")
		out, err := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...).Output()
		if err != nil {
			if strings.Contains(err.Error(), "255") {
				fmt.Println("\nNeed to run with sudo: `sudo ./fingerprint`\n")
			}
			log.Fatal(err)
		}
		f.WifiFingerprint, err = processOutput(out)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Sending fingerprint to " + address)
		b, _ := json.Marshal(f)
		req, err := http.NewRequest("POST", address, bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// fmt.Println("response Status:", resp.Status)
		// fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("response:", string(body))
		if strings.Contains(string(body), `"success":true`) == false {
			log.Fatal("Something wrong with server")
		}
	}

}
