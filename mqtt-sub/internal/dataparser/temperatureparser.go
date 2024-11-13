package dataparser

import (
	"encoding/json"
	"fmt"
	"time"
)

type TempParser struct {
	rawData    []byte
	datamap    map[string]float64 //metric to value
	data       tempData
	metername  string
	notifyRate time.Duration
}

type tempData struct {
	TempF    float64 `json:"tempF"`
	Humidity float64 `json:"humidity"`
	Action   string  `json:"action"`
}

func InitTempParser(name string) IDataParser {
	return &TempParser{
		datamap:    make(map[string]float64),
		metername:  name,
		notifyRate: time.Hour * 0, //we wont notify
	}
}

func (parser *TempParser) ParseData(data []byte) error {
	parser.rawData = data
	if jErr := json.Unmarshal(parser.rawData, &parser.data); jErr != nil {
		return fmt.Errorf("failed to unmarshal temperature data; %v", jErr)
	}

	if len(parser.datamap) == 0 {
		parser.datamap["temp_farenheight"] = parser.data.TempF
		parser.datamap["relative_humidity"] = parser.data.Humidity
	}
	return nil
}

func (parser *TempParser) GetDataMap() map[string]float64 {
	return parser.datamap
}

func (parser *TempParser) GetMeterName() string {
	return parser.metername
}

func (parser *TempParser) GetActionInfo() (action string, alertMsg string) {
	action = parser.data.Action
	alertMsg = ""
	return
}

func (parser *TempParser) NotificationRate() time.Duration {
	return parser.notifyRate
}
