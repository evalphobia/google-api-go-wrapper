package logging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/googleapi"
	SDK "google.golang.org/api/logging/v2"
)

// WriteData is data struct for LogEntry data.
type WriteData struct {
	Data            interface{} // required
	LogName         string      // required if commonLogName is empty
	Severity        Severity
	Labels          map[string]string
	PartialSuccess  bool
	Resource        *Resource
	Timestamp       time.Time
	InsertId        string
	Request         *http.Request
	Response        *http.Response
	Operation       *SDK.LogEntryOperation
	ForceSendFields []string
	NullFields      []string
}

// LogEntryList converts WriteData to *LogEntry.
func (d *WriteData) LogEntryList(projectID string) ([]*SDK.LogEntry, error) {
	switch v := d.Data.(type) {
	case []*SDK.LogEntry:
		return v, nil
	case *SDK.LogEntry:
		return []*SDK.LogEntry{v}, nil
	case SDK.LogEntry:
		return []*SDK.LogEntry{&v}, nil
	}

	ent := &SDK.LogEntry{
		LogName:         formatLogName(projectID, d.LogName),
		Severity:        string(d.Severity),
		Labels:          d.Labels,
		InsertId:        d.InsertId,
		Operation:       d.Operation,
		ForceSendFields: d.ForceSendFields,
		NullFields:      d.NullFields,
		HttpRequest:     toHttpRequest(d.Request, d.Response),
	}

	// set data
	switch v := d.Data.(type) {
	case string:
		ent.TextPayload = v
	case googleapi.RawMessage:
		ent.JsonPayload = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		ent.JsonPayload = b
	}

	if d.Resource != nil {
		ent.Resource = d.Resource.toMonitorResource()
	}
	if !d.Timestamp.IsZero() {
		ent.Timestamp = d.Timestamp.Format(time.RFC3339)
	}

	return []*SDK.LogEntry{ent}, nil
}

func toHttpRequest(req *http.Request, resp *http.Response) *SDK.HttpRequest {
	if req == nil && resp == nil {
		return nil
	}

	httpRequest := &SDK.HttpRequest{}
	if req != nil {
		httpRequest.Referer = req.Referer()
		httpRequest.RemoteIp = req.RemoteAddr
		httpRequest.RequestMethod = req.Method
		httpRequest.RequestSize = req.ContentLength
		httpRequest.RequestUrl = req.URL.String()
		httpRequest.UserAgent = req.UserAgent()
	}

	if resp != nil {
		httpRequest.ResponseSize = resp.ContentLength
		httpRequest.Status = int64(resp.StatusCode)
	}

	return httpRequest
}

func formatLogName(projectID, logID string) string {
	return fmt.Sprintf("projects/%s/logs/%s", projectID, logID)
}
