package monitoring

import SDK "google.golang.org/api/monitoring/v3"

// Resource is wrapper struct for SDK.MonitoredResource.
type Resource struct {
	Labels          map[string]string `json:"labels,omitempty"`
	Type            string            `json:"type,omitempty"`
	ForceSendFields []string          `json:"-"`
	NullFields      []string          `json:"-"`
}

func (r *Resource) toMonitorResource() *SDK.MonitoredResource {
	return &SDK.MonitoredResource{
		Labels:          r.Labels,
		Type:            r.Type,
		ForceSendFields: r.ForceSendFields,
		NullFields:      r.NullFields,
	}
}

// MetricKind is wrapper struct for Monitoring MetricKind.
type MetricKind string

// MetricKind list
const (
	MetricKindDefault    MetricKind = "METRIC_KIND_UNSPECIFIED"
	MetricKindGauge      MetricKind = "GAUGE"
	MetricKindDelta      MetricKind = "DELTA"
	MetricKindCumulative MetricKind = "CUMULATIVE"
)

// ValueType is wrapper struct for Monitoring ValueType.
type ValueType string

// ValueType list
const (
	ValueTypeDefault      ValueType = "VALUE_TYPE_UNSPECIFIED"
	ValueTypeBool         ValueType = "BOOL" // for GAUGE
	ValueTypeInt64        ValueType = "INT64"
	ValueTypeDouble       ValueType = "DOUBLE"
	ValueTypeString       ValueType = "STRING" // for GAUGE
	ValueTypeDistribution ValueType = "DISTRIBUTION"
	ValueTypeMoney        ValueType = "MONEY"
)
