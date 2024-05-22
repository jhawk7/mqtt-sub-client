package promexporter

import (
	"context"
	"encoding/json"
	"fmt"
	"mqtt-sub/internal/handlers"
	"strings"

	"go.opentelemetry.io/otel/metric"
)

type IExporter interface {
	Export(string, []byte)
}

type exporter struct {
	MeterName string
	MeterCh   chan int
	mp        metric.MeterProvider
}

func InitExporter(meterProvider *metric.MeterProvider) IExporter {
	return &exporter{
		mp: *meterProvider,
	}
}

func (exp *exporter) Export(name string, data []byte) {
	//export to prom
	parsedData := make(map[string]interface{})
	if uErr := json.Unmarshal(data, &parsedData); uErr != nil {
		err := fmt.Errorf("failed to unmarshal export data; %v", uErr)
		handlers.LogError(err, false)
		return
	}

	meterName := strings.ReplaceAll(name, "/", ".")
	handlers.LogInfo(fmt.Sprintf("exporting metrics for %v", meterName))

}

// func (exp *exporter) Export2(name string, data []byte) {
// 	//export to prom
// 	parsedData := make(map[string]interface{})
// 	if uErr := json.Unmarshal(data, &parsedData); uErr != nil {
// 		err := fmt.Errorf("failed to unmarshal export data; %v", uErr)
// 		handlers.LogError(err, false)
// 		return
// 	}

// 	meterName := strings.ReplaceAll(name, "/", ".")
// 	handlers.LogInfo(fmt.Sprintf("exporting metrics for %v", meterName))

// 	callback := func(ctx context.Context, result metric.Float64ObserverResult) {

// 	}

// 	meter := exp.mp.Meter(meterName)
// 	if _, gaugeErr := meter.NewFloat64GaugeObserver(meterName+".reading", callback); gaugeErr != nil {
// 		err := fmt.Errorf("guage observer error: %v", gaugeErr)
// 		handlers.LogError(err, false)
// 	}
// }
