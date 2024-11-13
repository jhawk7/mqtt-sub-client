package dataparser

import (
	"encoding/json"
	"fmt"
	"time"
)

type MoistureParser struct {
	rawData    []byte
	datamap    map[string]float64 //metric to value
	data       moistureData
	metername  string
	notifyRate time.Duration
}

type moistureData struct {
	MoistureReading float64 `json:"plant-moisture"`
	RawReading      float64 `json:"raw-reading"`
	Status          string  `json:"plant-status"`
	Threshold       float64 `json:"plant-threshold"`
	Action          string  `json:"action"`
	AlsertMsg       string  `json:"alert-msg,omitempty"`
}

func InitMoistureParser(name string) IDataParser {
	return &MoistureParser{
		datamap:    make(map[string]float64),
		metername:  name,
		notifyRate: time.Hour * 24, // notify once every 24hrs
	}
}

func (parser *MoistureParser) ParseData(data []byte) error {
	parser.rawData = data
	if jErr := json.Unmarshal(parser.rawData, &parser.data); jErr != nil {
		return fmt.Errorf("failed to unmarshal moisture data; %v", jErr)
	}

	if len(parser.datamap) == 0 {
		parser.datamap["moisture_percentage"] = parser.data.MoistureReading
		parser.datamap["plant_threshold"] = parser.data.Threshold
	}
	return nil
}

func (parser *MoistureParser) GetDataMap() map[string]float64 {
	return parser.datamap
}

func (parser *MoistureParser) GetMeterName() string {
	return parser.metername
}

func (parser *MoistureParser) GetActionInfo() (action string, alertMsg string) {
	action = parser.data.Action
	alertMsg = parser.data.AlsertMsg
	return
}

func (parser *MoistureParser) NotificationRate() time.Duration {
	return parser.notifyRate
}
