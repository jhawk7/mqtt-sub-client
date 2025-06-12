# MQTT SUB Client ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
This repo consist of two sub directories mqtt-sub and mqtt-test

## MQTT-SUB
MQTT-SUB subscribes to MQTT Topics and exports messages as metrics to a Prometheus server via Pushgateway.
- Topics can be added to the MQTT_TOPICS env var using ':' as a delimeter. 
- A `dataparser` interface is used to define how the service should parse the metrics of a topic's message. The dataparser implementation should store the metrics as k,v pairs in a ,map that the is passed to the `promexporter` to be pushed to prometheus via `pushgateway`. A dataparser implementation exists for each topic.
- The service uses the `notify` package to send alerts via SMTP to all SMTP clients (specified in the **SMTP_CLIENTS** env var).
A client will be notified with the given alert message (`alert-msg`) if the `action` field of the MQTT message is set to **alert**.

- List of env vars
```
MQTT_SERVER=<mqtt_server_uri>
MQTT_PORT=<mqtt_server_port>
MQTT_USER=<mqtt_user_name>
MQTT_PASS=<mqtt_user_pass>
MQTT_TOPICS=<list of mqtt topics separated by ':'>
PUSH_URL=<pushgateway_url>
SMTP_EMAIL=<SMTP email sender>
SMTP_PASS=<SMTP email pass>
SMTP_SERVER=<SMTP server>
SMTP_CLIENTS=<list of SMTP receipients/clients separated by ':'>
```

## MQTT-TEST
Contains code to test pub/sub against an MQTT service given the env vars to connect to the service.

