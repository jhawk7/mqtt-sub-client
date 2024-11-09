package mqttclient

import (
	"fmt"
	"mqtt-sub/internal/common"
	"mqtt-sub/internal/dataparser"
	"mqtt-sub/internal/notify"
	"mqtt-sub/internal/promexporter"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	promExp  promexporter.IExporter
	notifier notify.INotifier
	topics   []string
	mChan    chan mqtt.Message
)

const (
	LOG_ACTION    = "log"
	NOTIFY_ACTION = "alert"
)

type IClient interface {
	sub()
	Disconnect()
	Listen()
}

type client struct {
	mqttClient mqtt.Client
}

var connHandler mqtt.OnConnectHandler = func(mclient mqtt.Client) {
	common.LogInfo("successfully connected to mqtt server")
}

var lostHandler mqtt.ConnectionLostHandler = func(mclient mqtt.Client, err error) {
	connErr := fmt.Errorf("lost connection to mqtt server; %v", err)
	common.LogError(connErr, false)
}

var msgHandler mqtt.MessageHandler = func(mclient mqtt.Client, msg mqtt.Message) {
	info := fmt.Sprintf("messaged received [id: %d] [topic: %v], [payload: %s]", msg.MessageID(), msg.Topic(), msg.Payload())
	common.LogInfo(info)

	mChan <- msg
}

func InitClient(p promexporter.IExporter, n notify.INotifier, config *common.Config) IClient {
	//set client options
	broker := config.MQTTServer
	port := config.MQTTPort
	user := config.MQTTUser
	pass := config.MQTTPass
	topics = config.MQTTTopics

	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%v:%v", broker, port))
	opts.SetClientID("mqtt-sub-client")
	opts.SetUsername(user)
	opts.SetPassword(pass)
	opts.SetKeepAlive(time.Second * 10)
	opts.SetCleanSession(false) //disabling clean session on client reconnect so that messages will resume on reconnect
	opts.OnConnect = connHandler
	opts.OnConnectionLost = lostHandler
	mclient := mqtt.NewClient(opts)
	if token := mclient.Connect(); token.Wait() && token.Error() != nil {
		err := fmt.Errorf("mqtt connection failed; %v", token.Error())
		common.LogError(err, true)
	}

	mChan = make(chan mqtt.Message, 1)
	c := &client{mqttClient: mclient}
	promExp = p
	notifier = n
	return c
}

func (c *client) sub() {
	filters := make(map[string]byte)
	for _, topic := range topics {
		common.LogInfo(fmt.Sprintf("subscribing to topic [%v]", topic))
		filters[topic] = 2 //qos: 0 - no standard, 1 - "atleast once", 2 - exactly once
	}

	if token := c.mqttClient.SubscribeMultiple(filters, msgHandler); token.Wait() && token.Error() != nil {
		err := fmt.Errorf("failed to subscribe to topics %v; error %v", topics, token.Error())
		common.LogError(err, true)
	}
	common.LogInfo(fmt.Sprintf("subscribed to topics %v", topics))
}

func (c *client) Disconnect() {
	common.LogInfo("disconnecting from mqtt server..")
	c.mqttClient.Disconnect(5000)
}

func (c *client) Listen() {
	c.sub()
	for msg := range mChan {
		common.LogInfo(fmt.Sprintf("parsing and exporting message from topic %v", msg.Topic()))
		parser, pErr := parseMsg(msg)
		if pErr != nil {
			common.LogError(pErr, false)
			continue
		}

		//send to prom
		promExp.ExportMetrics(parser.GetMeterName(), parser.GetDataMap())
		action, alertmsg := parser.GetActionInfo()

		if action == NOTIFY_ACTION {
			if nErr := notifier.Notify(alertmsg); nErr != nil {
				common.LogError(nErr, false)
			}
		}
	}
}

func parseMsg(msg mqtt.Message) (parser dataparser.IDataParser, err error) {
	parser, dpErr := dataparser.InitDataParser(msg.Topic())
	if dpErr != nil {
		err = dpErr
		return
	}

	if pErr := parser.ParseData(msg.Payload()); pErr != nil {
		err = pErr
		return
	}

	return
}
