package promexporter

import (
	"fmt"
	"mqtt-sub/internal/common"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

type IExporter interface {
	ExportMetrics(string, map[string]float64)
}

type exporter struct {
	url string
}

func InitExporter(pushUrl string) IExporter {
	return &exporter{url: pushUrl}
}

func (exp *exporter) ExportMetrics(metername string, datamap map[string]float64) {
	common.LogInfo(fmt.Sprintf("exporting prom data for %v", metername))
	var collectors []prometheus.Collector
	registry := prometheus.NewRegistry()

	for metric, datapoint := range datamap {
		opts := prometheus.GaugeOpts{Name: metric}
		guage := prometheus.NewGauge(opts)
		guage.Set(datapoint)
		collectors = append(collectors, guage)
	}

	registry.MustRegister(collectors...)

	if pushErr := push.New(exp.url, metername).
		Gatherer(registry).
		Grouping("topic", metername).
		Push(); pushErr != nil {
		err := fmt.Errorf("failed to push data to prometheus; %v", pushErr)
		common.LogError(err, false)
		return
	}
}
