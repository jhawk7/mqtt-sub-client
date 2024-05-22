package main

import (
	"fmt"
	"mqtt-sub/internal/handlers"
	"mqtt-sub/internal/mqttclient"
	"mqtt-sub/internal/notify"
	"mqtt-sub/internal/promexporter"
	"os"
	"os/signal"

	"github.com/jhawk7/go-opentel/opentel"
)

func main() {
	if opentelErr := opentel.InitOpentelProviders(); opentelErr != nil {
		err := fmt.Errorf("failed to init opentel providers; %v", opentelErr)
		handlers.LogError(err, true)
	}

	mp := opentel.GetMeterProvider()
	promExp := promexporter.InitExporter(&mp)
	defer func() {
		opentel.ShutdownOpentelProviders()
	}()
	notifier := notify.InitNotifier()
	cMqtt := mqttclient.InitClient(promExp, notifier)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	<-ch //blocks until signal from ch is received
	cMqtt.Disconnect()
}
