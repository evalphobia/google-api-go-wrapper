package logging

import (
	"fmt"
	"sync"

	SDK "google.golang.org/api/logging/v2"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "logging"
)

// Logger repesents stackdriver logger.
type Logger struct {
	service   *SDK.Service
	logger    log.Logger
	projectID string

	commonLogName        string
	commonResource       *SDK.MonitoredResource
	commonLabels         map[string]string
	commonForceFields    []string
	commonNullFields     []string
	commonPartialSuccess bool

	writeMu    sync.Mutex
	writeSpool map[string][]*SDK.LogEntry
}

// NewLogger returns initialized *Logger
func NewLogger(conf config.Config, projectID string) (*Logger, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = append(conf.Scopes, SDK.LoggingWriteScope)
	}
	cli, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := SDK.New(cli)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		service:    svc,
		logger:     log.DefaultLogger,
		projectID:  projectID,
		writeSpool: make(map[string][]*SDK.LogEntry),
	}
	return logger, nil
}

// SetCommonLogName sets common Resource.
func (l *Logger) SetCommonLogName(logName string) {
	l.commonLogName = formatLogName(l.projectID, logName)
}

// SetCommonResource sets common Resource.
func (l *Logger) SetCommonResource(r *Resource) {
	l.commonResource = r.toMonitorResource()
}

// SetCommonLables sets common labels.
func (l *Logger) SetCommonLables(labels map[string]string) {
	l.commonLabels = labels
}

// SetCommonForceFields sets common FourceSendsFields.
func (l *Logger) SetCommonForceFields(fields []string) {
	l.commonForceFields = fields
}

// SetCommonNullFields sets common NullFields.
func (l *Logger) SetCommonNullFields(fields []string) {
	l.commonNullFields = fields
}

// SetCommonPartialSuccess sets common PartialSuccess.
func (l *Logger) SetCommonPartialSuccess(partialSuccess bool) {
	l.commonPartialSuccess = partialSuccess
}

// SetLogger sets internal API logger.
func (l *Logger) SetLogger(logger log.Logger) {
	l.logger = logger
}

// Write sends log data to stackdriver's log.
func (l *Logger) Write(data WriteData) error {
	entryList, err := l.buildLogEntryList(data)
	if err != nil {
		return err
	}

	req := l.CreateWriteRequest(entryList, data.LogName)
	_, err = l.service.Entries.Write(req).Do()
	if err != nil {
		l.Errorf("error on `Write` operation; projectID=%s, error=%s", l.projectID, err.Error())
	}
	return err
}

// Add adds the log entry to write spool.
func (l *Logger) Add(data WriteData) error {
	entryList, err := l.buildLogEntryList(data)
	if err != nil {
		return err
	}

	logName := data.LogName
	l.writeMu.Lock()
	l.writeSpool[logName] = append(l.writeSpool[logName], entryList...)
	l.writeMu.Unlock()
	return nil
}

// FlushAll executes Write operation from the write spool.
func (l *Logger) FlushAll() error {
	l.writeMu.Lock()
	defer l.writeMu.Unlock()
	for logName, entryList := range l.writeSpool {
		req := l.CreateWriteRequest(entryList, logName)
		_, err := l.service.Entries.Write(req).Do()
		if err != nil {
			l.Errorf("error on `Write` operation; projectID=%s, logName=%s, error=%s", l.projectID, logName, err.Error())
			return err
		}

		l.writeSpool[logName] = nil
	}
	return nil
}

// CreateWriteRequest creates *SDK.WriteLogEntriesRequest from WriteData.
func (l *Logger) CreateWriteRequest(rows []*SDK.LogEntry, logName string) *SDK.WriteLogEntriesRequest {
	return &SDK.WriteLogEntriesRequest{
		Entries:         rows,
		Labels:          l.commonLabels,
		LogName:         l.commonLogName,
		PartialSuccess:  l.commonPartialSuccess,
		Resource:        l.commonResource,
		ForceSendFields: l.commonForceFields,
		NullFields:      l.commonNullFields,
	}
}

func (l *Logger) buildLogEntryList(data WriteData) ([]*SDK.LogEntry, error) {
	switch {
	case data.Data == nil:
		return nil, fmt.Errorf("error: data is required")
	case data.Resource == nil && l.commonResource == nil:
		return nil, fmt.Errorf("error: resource is required")
	}

	return data.LogEntryList(l.projectID)
}

// Errorf logging error information.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Errorf(serviceName, format, v...)
}
