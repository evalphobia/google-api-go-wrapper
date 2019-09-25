package opencensus

import (
	"go.opencensus.io/stats/view"
)

func RegisterView(v *view.View) error {
	return view.Register(v)
}

func (e *Exporter) StartMetricsExporter() error {
	return e.Exporter.StartMetricsExporter()
}

func (e *Exporter) StopMetricsExporter() {
	e.Exporter.StopMetricsExporter()
}
