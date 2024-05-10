package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var connHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	LogInfo("successfully connected to mqtt server")
}

var lostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	connErr := fmt.Errorf("lost connection to mqtt server; %v", err)
	LogError(connErr, false)

}

var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	info := fmt.Sprintf("messaged received [id: %d] [topic: %v], [payload: %s]", msg.MessageID(), msg.Topic(), msg.Payload())
	LogInfo(info)

	TriageMsg(msg)
}

func main() {
	//set client options
	broker := os.Getenv("MQTT_SERVER")
	port := os.Getenv("MQTT_PORT")
	user := os.Getenv("MQTT_USER")
	pass := os.Getenv("MQTT_PASS")
	topics := strings.Split(os.Getenv("MQTT_TOPICS"), ":")

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", broker, port))
	opts.SetClientID("test_client")
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetKeepAlive(time.Second * 10)
	opts.OnConnect = connHandler
	opts.OnConnectionLost = lostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err := fmt.Errorf("mqtt connection failed; %v", token.Error())
		fmt.Println(err)
	}

	sub(&client, topics)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	<-ch //blocks until signal from ch is received
	client.Disconnect(5000)
}

func sub(client *mqtt.Client, topics []string) {
	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = '1' //atleast once
	}

	token := (*client).SubscribeMultiple(filters, msgHandler)
	token.Wait()
	LogInfo(fmt.Sprintf("subscribed to topics %v", topics))
}

func TriageMsg(msg mqtt.Message) {
	// send alert
	// send to prometheus
}

func LogError(err error, fatal bool) {
	if err != nil {
		if fatal {
			log.Fatalf("fatal error: %v", err)
		} else {
			log.Errorf("error: %v", err)
		}
	}
}

func LogInfo(msg string) {
	log.Info(msg)
}
