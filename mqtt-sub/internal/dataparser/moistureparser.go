package dataparser

import (
	"encoding/json"
	"fmt"
)

type MoistureParser struct {
	rawData   []byte
	datamap   map[string]float64 //metric to value
	data      moistureData
	metername string
}

type moistureData struct {
	moistureReading float64 `json:"plant-moisture"`
	rawReading      float64 `json:"raw-reading"`
	status          string  `json:"plant-status"`
	threshold       float64 `json:"plant-threshold"`
	action          string  `json:"action"`
	alsertMsg       string  `json:"alert-msg"`
}

func InitMoistureParser(name string) IDataParser {
	return &MoistureParser{
		datamap:   make(map[string]float64),
		metername: name,
	}
}

func (parser *MoistureParser) ParseData(data []byte) error {
	parser.rawData = data
	if jErr := json.Unmarshal(parser.rawData, &parser.data); jErr != nil {
		return fmt.Errorf("failed to unmarshal moisture data; %v", jErr)
	}

	if len(parser.datamap) == 0 {
		parser.datamap["moisture_percentage"] = parser.data.moistureReading
		parser.datamap["plant_threshold"] = parser.data.threshold
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
	action = parser.data.action
	alertMsg = parser.data.alsertMsg
	return
}
