package dataparser

import (
	"encoding/json"
	"fmt"
)

type TempParser struct {
	rawData   []byte
	datamap   map[string]float64 //metric to value
	data      tempData
	metername string
}

type tempData struct {
	tempF    float64 `json:"tempF"`
	humidity float64 `json:"humidity"`
}

func InitTempParser(name string) *TempParser {
	return &TempParser{
		datamap:   make(map[string]float64),
		metername: name,
	}
}

func (parser *TempParser) ParseData(data []byte) error {
	parser.rawData = data
	if jErr := json.Unmarshal(parser.rawData, &parser.data); jErr != nil {
		return fmt.Errorf("failed to unmarshal temperature data; %v", jErr)
	}

	if len(parser.datamap) == 0 {
		parser.datamap["temp_farenheight"] = parser.data.tempF
		parser.datamap["relative_humidity"] = parser.data.humidity
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
	action = ""
	alertMsg = ""
	return
}
