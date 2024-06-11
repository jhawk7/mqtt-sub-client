package dataparser

import "testing"

func TestMoistureParser(t *testing.T) {
	m := InitMoistureParser("picow/house/plant-moisture")
	fakedata := []byte(`{"action": "log", "plant-moisture": 33.75575, "raw-reading": 40601, "plant-status": "ok"}`)
	if parseErr := m.ParseData(fakedata); parseErr != nil {
		t.Logf("unexpedted failure to parse data %v", parseErr)
		t.Fail()
	}

	datamap := m.GetDataMap()
	if datamap["moisture_percentage"] != 33.75575 {
		t.Logf("unexpected moisture_percentage value %v", datamap["moisture_percentage"])
		t.Fail()
	}

	if action, _ := m.GetActionInfo(); action != "log" {
		t.Logf("unexpected action %v", action)
		t.Fail()
	}

	if m.GetMeterName() == "" {
		t.Logf("unexpected metername %v", m.GetMeterName())
		t.Fail()
	}
}
