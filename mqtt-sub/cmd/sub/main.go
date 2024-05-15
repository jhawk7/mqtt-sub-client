package main

import (
	"mqtt-sub/internal/mqttclient"
	"mqtt-sub/internal/notify"
	"mqtt-sub/internal/promexporter"
	"os"
	"os/signal"
)

func main() {
	promExp := promexporter.InitExporter()
	notifier := notify.InitNotifier()
	cMqtt := mqttclient.InitClient(promExp, notifier)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	<-ch //blocks until signal from ch is received
	cMqtt.Disconnect()
}
