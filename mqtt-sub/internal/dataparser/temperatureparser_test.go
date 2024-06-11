package dataparser

import "testing"

func TestTempParser(t *testing.T) {
	m := InitTempParser("picow/tempF")
	fakedata := []byte(`{"action": "log", "tempF": 71.18883, "humidity": 47.84314}`)
	if parseErr := m.ParseData(fakedata); parseErr != nil {
		t.Logf("unexpedted failure to parse data %v", parseErr)
		t.Fail()
	}

	datamap := m.GetDataMap()
	if datamap["temp_farenheight"] != 71.18883 {
		t.Logf("unexpected temperature value %v", datamap["temp_farenheight"])
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
