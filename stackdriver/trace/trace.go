package trace

import (
	"golang.org/x/net/context"

	GCP "cloud.google.com/go/trace"
	SDK "google.golang.org/api/cloudtrace/v1"
	"google.golang.org/api/option"

	"github.com/evalphobia/google-api-go-wrapper/config"
)

const (
	serviceName     = "trace"
	LabelStatusCode = `trace.cloud.google.com/http/status_code`
)

// Trace is wrapper struct of GCP.Client.
type Trace struct {
	*GCP.Client
}

// NewTrace returns initialized *Trace.
func NewTrace(ctx context.Context, conf config.Config, projectID string) (*Trace, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = []string{SDK.TraceAppendScope}
	}

	httpClient, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := GCP.NewClient(ctx, projectID, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &Trace{
		Client: svc,
	}, nil
}

// NewLimitedSampler creates wrapper struct of GCP.SamplingPolicy.
func NewLimitedSampler(fraction, maxqps float64) (*SamplingPolicy, error) {
	s, err := GCP.NewLimitedSampler(fraction, maxqps)
	return &SamplingPolicy{
		SamplingPolicy: s,
	}, err
}

// SamplingPolicy is wrapper struct of GCP.SamplingPolicy.
type SamplingPolicy struct {
	GCP.SamplingPolicy
}
