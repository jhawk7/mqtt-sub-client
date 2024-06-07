package mqttclient

import (
	"encoding/json"
	"fmt"
	"mqtt-sub/internal/handlers"
	"mqtt-sub/internal/notify"
	"mqtt-sub/internal/promexporter"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	promExp  promexporter.IExporter
	notifier notify.INotifier
	topics   []string
)

const (
	LOG_ACTION    = "log"
	NOTIFY_ACTION = "notify"
)

type IClient interface {
	sub()
	Disconnect()
}

type client struct {
	mqttClient mqtt.Client
}

type parsedMessage struct {
	data     interface{} `json:"data"`
	action   string      `json:"action"`
	alertMsg string      `json:"alert-msg,omitempty"`
}

var connHandler mqtt.OnConnectHandler = func(mclient mqtt.Client) {
	handlers.LogInfo("successfully connected to mqtt server")
}

var lostHandler mqtt.ConnectionLostHandler = func(mclient mqtt.Client, err error) {
	connErr := fmt.Errorf("lost connection to mqtt server; %v", err)
	handlers.LogError(connErr, false)

}

var msgHandler mqtt.MessageHandler = func(mclient mqtt.Client, msg mqtt.Message) {
	info := fmt.Sprintf("messaged received [id: %d] [topic: %v], [payload: %s]", msg.MessageID(), msg.Topic(), msg.Payload())
	handlers.LogInfo(info)

	triageMsg(msg)
}

func InitClient(p promexporter.IExporter, n notify.INotifier) IClient {
	//set client options
	broker := os.Getenv("MQTT_SERVER")
	port := os.Getenv("MQTT_PORT")
	user := os.Getenv("MQTT_USER")
	pass := os.Getenv("MQTT_PASS")
	topics = strings.Split(os.Getenv("MQTT_TOPICS"), ":")

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", broker, port))
	opts.SetClientID("mqtt-sub-client")
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetKeepAlive(time.Second * 10)
	opts.OnConnect = connHandler
	opts.OnConnectionLost = lostHandler
	mclient := mqtt.NewClient(opts)
	if token := mclient.Connect(); token.Wait() && token.Error() != nil {
		err := fmt.Errorf("mqtt connection failed; %v", token.Error())
		handlers.LogError(err, true)
	}

	c := &client{mqttClient: mclient}
	promExp = p
	notifier = n

	c.sub()
	return c
}

func (c *client) sub() {
	filters := make(map[string]byte)
	for _, topic := range topics {
		filters[topic] = '1' //qos - atleast once
	}

	token := c.mqttClient.SubscribeMultiple(filters, msgHandler)
	token.Wait()
	handlers.LogInfo(fmt.Sprintf("subscribed to topics %v", topics))
}

func (c *client) Disconnect() {
	c.mqttClient.Disconnect(5000)
}

func triageMsg(msg mqtt.Message) {
	var pmsg parsedMessage
	if mErr := json.Unmarshal(msg.Payload(), &pmsg); mErr != nil {
		err := fmt.Errorf("failed to parse incoming mqtt message; %v", mErr)
		handlers.LogError(err, false)
		return
	}

	//send to prom
	promExp.Export(msg.Topic(), pmsg.data)

	if pmsg.action == NOTIFY_ACTION {
		//notify
		//notifier.Notify(pmsg.data)
	}
}
