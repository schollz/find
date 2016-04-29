package main

import (
	"fmt"
	"strings"
	"time"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var mqttClients map[string]*MQTT.Client

func init() {
	go clearMqttPool()
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
