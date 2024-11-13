package common

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	MQTTServer string
	MQTTPort   string
	MQTTUser   string
	MQTTPass   string
	MQTTTopics []string
	PushUrl    string
	ClientFrom string
	SMTPPass   string
	SMTPServer string
	SMTPPort   string
	ClientsTo  []string
	RedisHost  string
	RedisPass  string
}

func InitConfig() *Config {
	//var c Config
	server, found := os.LookupEnv("MQTT_SERVER")
	if !found {
		LogError(fmt.Errorf("MQTT_SERVER not set"), true)
	}

	port, found := os.LookupEnv("MQTT_PORT")
	if !found {
		LogError(fmt.Errorf("MQTT_PORT not set"), true)
	}

	user, found := os.LookupEnv("MQTT_USER")
	if !found {
		LogError(fmt.Errorf("MQTT_USER not set"), true)
	}

	pass, found := os.LookupEnv("MQTT_PASS")
	if !found {
		LogError(fmt.Errorf("MQTT_PASS not set"), true)
	}

	topics, found := os.LookupEnv("MQTT_TOPICS")
	if !found {
		LogError(fmt.Errorf("MQTT_TOPICS not set"), true)
	}

	pushUrl, found := os.LookupEnv("PUSH_URL")
	if !found {
		LogError(fmt.Errorf("PUSH_URL not set"), true)
	}

	clientFrom, found := os.LookupEnv("SMTP_EMAIL")
	if !found {
		LogError(fmt.Errorf("SMTP_EMAIL not set"), true)
	}

	epass, found := os.LookupEnv("SMTP_PASS")
	if !found {
		LogError(fmt.Errorf("SMTP_PASS not set"), true)
	}

	smtpServer, found := os.LookupEnv("SMTP_SERVER")
	if !found {
		LogError(fmt.Errorf("SMTP_SERVER not set"), true)
	}

	smtpPort, found := os.LookupEnv("SMTP_PORT")
	if !found {
		LogError(fmt.Errorf("SMTP_PORT not set"), true)
	}

	clientsTo, found := os.LookupEnv("SMTP_CLIENTS")
	if !found {
		LogError(fmt.Errorf("SMTP_CLIENT not set"), true)
	}

	redisHost, found := os.LookupEnv("REDIS_HOST")
	if !found {
		LogError(fmt.Errorf("REDIS_HOST not set"), true)
	}

	redisPass, found := os.LookupEnv("REDIS_PASS")
	if !found {
		LogError(fmt.Errorf("REDIS_PASS not set"), true)
	}

	return &Config{
		MQTTServer: server,
		MQTTPort:   port,
		MQTTUser:   user,
		MQTTPass:   pass,
		MQTTTopics: strings.Split(topics, ":"),
		PushUrl:    pushUrl,
		ClientFrom: clientFrom,
		ClientsTo:  strings.Split(clientsTo, ":"),
		SMTPPass:   epass,
		SMTPServer: smtpServer,
		SMTPPort:   smtpPort,
		RedisHost:  redisHost,
		RedisPass:  redisPass,
	}
}
