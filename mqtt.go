package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var mqttClients map[string]*MQTT.Client

func init() {
	if RuntimeArgs.Mqtt {
		go clearMqttPool()

		server := "tcp://ml.internalpositioning.com:1883"
		name := "user" + strconv.Itoa(rand.Intn(1000))

		subTopic := strings.Join([]string{"/find/", "fingerprint", "/+"}, "")
		opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(name).SetCleanSession(true)

		opts.OnConnect = func(c *MQTT.Client) {
			if token := c.Subscribe(subTopic, 1, messageReceived); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
		}

		client := MQTT.NewClient(opts)

		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Printf("Connected to ", name, server)
	}
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
	msgFrom := topics[len(topics)-1]
	fmt.Print(msgFrom + ": " + string(msg.Payload()))
}
