package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var pubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("message published to mqtt; [topic: %v]", msg.Topic())
}

var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("successfully connected to mqtt server")
}

var lostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("lost connection to mqtt server; %v", err)
}

var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("messaged received [id: %d] [topic: %v], [payload: %s]", msg.MessageID(), msg.Topic(), msg.Payload())
}

func test() {
	//set client options
	broker := os.Getenv("MQTT_SERVER")
	port := os.Getenv("MQTT_PORT")
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", broker, port))
	opts.SetClientID("test_client")
	opts.SetUsername("test")
	opts.SetDefaultPublishHandler(pubHandler)
	opts.OnConnect = connHandler
	opts.OnConnectionLost = lostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err := fmt.Errorf("mqtt connection failed; %v", token.Error())
		fmt.Println(err)
	}

	publish(&client)
	sub(&client)

	client.Disconnect(5000)
}

func publish(client *mqtt.Client) {
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := (*client).Publish("topic/test", 1, true, text)
		token.Wait()
		time.Sleep(time.Second * 2)
	}
}

func sub(client *mqtt.Client) {
	topic := "topic/test"
	token := (*client).Subscribe(topic, 1, msgHandler)
	token.Wait()
	fmt.Println("Subscribed to topic: %s", topic)
}
