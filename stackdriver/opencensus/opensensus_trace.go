package opencensus

import (
	"context"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"google.golang.org/api/option"

	"github.com/evalphobia/google-api-go-wrapper/config"
)

type Exporter struct {
	*stackdriver.Exporter
}

// NewExporter creates new exporter stackdriver exporter to opencensus.
func NewExporter(ctx context.Context, conf config.Config, projectID string) (*Exporter, error) {
	ts, err := conf.TokenSource(ctx)
	if err != nil {
		return nil, err
	}

	opts := []option.ClientOption{option.WithTokenSource(ts)}
	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:               projectID,
		TraceClientOptions:      opts,
		MonitoringClientOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	return &Exporter{
		Exporter: exporter,
	}, nil
}

func (e *Exporter) RegisterTrace() {
	trace.RegisterExporter(e.Exporter)
}
