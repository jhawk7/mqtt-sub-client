package main

import (
	"mqtt-sub/internal/common"
	"mqtt-sub/internal/mqttclient"
	"mqtt-sub/internal/notify"
	"mqtt-sub/internal/promexporter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := common.InitConfig()
	promExp := promexporter.InitExporter(c.PushUrl)
	notifier := notify.InitNotifier(c)
	cMqtt := mqttclient.InitClient(promExp, notifier, c)
	go cMqtt.Listen()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch //blocks until signal from ch is received
	cMqtt.Disconnect()
}
