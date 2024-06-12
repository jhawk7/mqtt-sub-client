package dataparser

import (
	"fmt"
	"strings"
)

// All DataParsers implement IDataParser Interface
type IDataParser interface {
	ParseData([]byte) error
	GetDataMap() map[string]float64
	GetMeterName() string
	GetActionInfo() (action string, alertMsg string)
}

func InitDataParser(topic string) (parser IDataParser, err error) {
	meter := strings.ReplaceAll(topic, "/", ".")
	switch meter {
	case "picow.house.plant-moisture":
		parser = InitMoistureParser(meter)
	case "picow.tempF":
		parser = InitTempParser(meter)
	default:
		err = fmt.Errorf("unknown mqtt topic")
		return
	}

	return
}
