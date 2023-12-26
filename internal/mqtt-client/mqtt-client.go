package mqtt_client

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func GetMqttClient(clientId string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1985")
	opts.SetClientID(clientId)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func Publish(client mqtt.Client, msg []byte) {
	token := client.Publish("sensor", 0, false, msg)
	token.Wait()
	fmt.Println("Message published successfully!")
}

func Subscribe(client mqtt.Client, subHandler func(client mqtt.Client, msg mqtt.Message)) {
	topic := "sensor"
	token := client.Subscribe(topic, 1, subHandler)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
