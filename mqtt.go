package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var mqttClients map[string]*MQTT.Client

func setupMqtt() {
	go clearMqttPool()

	server := "tcp://ml.internalpositioning.com:1883"
	name := "user" + strconv.Itoa(rand.Intn(1000))

	opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(name).SetCleanSession(true)

	opts.OnConnect = func(c *MQTT.Client) {
		if token := c.Subscribe("/fingerprint/track/+", 1, messageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := c.Subscribe("/fingerprint/learn/+", 1, messageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Connected to ", name, server)
}

func clearMqttPool() {
	for {
		mqttClients = make(map[string]*MQTT.Client)
		time.Sleep(time.Minute * 30)
	}
}

func sendMQTTMessage(message string, group string, user string) error {
	server := "tcp://ml.internalpositioning.com:1883"
	room := group
	name := user
	// subTopic := strings.Join([]string{"/find/", room, "/+"}, "")
	pubTopic := strings.Join([]string{"/find/", room, "/", name}, "")

	if _, ok := mqttClients[group]; !ok {
		opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(name).SetCleanSession(true)
		mqttClients[group] = MQTT.NewClient(opts)
		if token := mqttClients[group].Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	if token := mqttClients[group].Publish(pubTopic, 1, false, message); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Failed to send message")
	}
	return nil
}

func messageReceived(client *MQTT.Client, msg MQTT.Message) {
	topics := strings.Split(msg.Topic(), "/")
	Debug.Println(topics)
	if len(topics) != 4 {
		return
	}
	route := strings.TrimSpace(topics[2])
	if (route != "track" && route != "learn") || strings.TrimSpace(topics[1]) != "fingerprint" {
		return
	}

	url := "http://127.0.0.1" + RuntimeArgs.Port + "/" + route

	var jsonStr = msg.Payload()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client2 := &http.Client{}
	resp, err := client2.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
