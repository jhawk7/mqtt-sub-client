package dataparser

// All DataParsers implement IDataParser Interface
type IDataParser interface {
	ParseData([]byte) error
	GetDataMap() map[string]float64
	GetMeterName() string
	GetActionInfo() (action string, alertMsg string)
}
