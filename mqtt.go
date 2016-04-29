package main

import (
	"fmt"
	"strings"

	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

func sendMQTTMessage(message string, group string, user string) error {
	server := "tcp://ml.internalpositioning.com:1883"
	room := group
	name := user
	// subTopic := strings.Join([]string{"/find/", room, "/+"}, "")
	pubTopic := strings.Join([]string{"/find/", room, "/", name}, "")

	opts := MQTT.NewClientOptions().AddBroker(server).SetClientID(name).SetCleanSession(true)

	// opts.OnConnect = func(c *MQTT.Client) {
	// 	if token := c.Subscribe(subTopic, 1, messageReceived); token.Wait() && token.Error() != nil {
	// 		panic(token.Error())
	// 	}
	// }

	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	// Debug.Println("Connected ", name, server)

	if token := client.Publish(pubTopic, 1, false, message); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Failed to send message")
	}
	return nil
}
