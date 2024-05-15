package promexporter

import (
	"github.com/jhawk7/go-opentel/opentel"
	"go.opentelemetry.io/otel/metric"
)

type IExporter interface {
	Export([]byte)
}

type exporter struct {
	mp metric.MeterProvider
}

func InitExporter() IExporter {
	opentel.InitOpentelProviders()
	mProvider := opentel.GetMeterProvider()

	return &exporter{
		mp: mProvider,
	}
}

func (exp *exporter) Export(data []byte) {
	//export to prom
}
