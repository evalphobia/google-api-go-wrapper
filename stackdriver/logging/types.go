package logging

import SDK "google.golang.org/api/logging/v2"

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

// Severity is wrapper struct for Log Severity.
type Severity string

// severity list
const (
	SeverityDefault   Severity = "DEFAULT"
	SeverityDebug     Severity = "DEBUG"
	SeverityInfo      Severity = "INFO"
	SeverityNotice    Severity = "NOTICE"
	SeverityWarning   Severity = "WARNING"
	SeverityError     Severity = "ERROR"
	SeverityCritical  Severity = "CRITICAL"
	SeverityAlert     Severity = "ALERT"
	SeverityEmergency Severity = "EMERGENCY"
)
