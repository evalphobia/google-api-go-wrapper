package monitoring

import (
	"fmt"
	"sync"

	SDK "google.golang.org/api/monitoring/v3"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "monitoring"
)

// Monitor repesents stackdriver monitoring.
type Monitor struct {
	service   *SDK.Service
	logger    log.Logger
	projectID string

	commonResource    *Resource
	commonLabels      map[string]string
	commonForceFields []string
	commonNullFields  []string

	writeMu    sync.Mutex
	writeSpool map[string][]*SDK.TimeSeries
}

// NewMonitor returns initialized *Monitor
func NewMonitor(conf config.Config, projectID string) (*Monitor, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = append(conf.Scopes, SDK.MonitoringScope)
	}
	cli, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := SDK.New(cli)
	if err != nil {
		return nil, err
	}

	monitor := &Monitor{
		service:    svc,
		logger:     log.DefaultLogger,
		projectID:  projectID,
		writeSpool: make(map[string][]*SDK.TimeSeries),
	}
	return monitor, nil
}

// SetCommonResource sets common Resource.
func (m *Monitor) SetCommonResource(r *Resource) {
	m.commonResource = r
}

// SetCommonForceFields sets common FourceSendsFields.
func (m *Monitor) SetCommonForceFields(fields []string) {
	m.commonForceFields = fields
}

// SetCommonNullFields sets common NullFields.
func (m *Monitor) SetCommonNullFields(fields []string) {
	m.commonNullFields = fields
}

// SetLogger sets internal API logger.
func (m *Monitor) SetLogger(l log.Logger) {
	m.logger = l
}

// Create sends custom metric data to stackdriver.
func (m *Monitor) Create(data Data) error {
	tsList, err := m.buildTimeSeriesList(data)
	if err != nil {
		return err
	}

	req := m.CreateTimeSeriesRequest(tsList)
	_, err = m.service.Projects.TimeSeries.Create(m.formatProjectName(), req).Do()
	if err != nil {
		m.Errorf("error on `TimeSeries.Create` operation; projectID=%s, error=%s", m.projectID, err.Error())
	}
	return err
}

// Add adds the timeseries to write spool.
func (m *Monitor) Add(data Data) error {
	tsList, err := m.buildTimeSeriesList(data)
	if err != nil {
		return err
	}

	spool := m.writeSpool
	m.writeMu.Lock()
	defer m.writeMu.Unlock()
	for _, ts := range tsList {
		name := ts.Metric.Type
		spool[name] = append(spool[name], ts)
	}
	return nil
}

// FlushAll executes TimeSeries.Create operation from the write spool.
func (m *Monitor) FlushAll() error {
	m.writeMu.Lock()
	defer m.writeMu.Unlock()
	for name, tsList := range m.writeSpool {
		req := m.CreateTimeSeriesRequest(tsList)
		_, err := m.service.Projects.TimeSeries.Create(m.formatProjectName(), req).Do()
		if err != nil {
			m.Errorf("error on `TimeSeries.Create` operation; projectID=%s, name=%s, error=%s", m.projectID, name, err.Error())
			return err
		}

		delete(m.writeSpool, name)
	}
	return nil
}

// CreateTimeSeriesRequest creates *SDK.TimeSeries from WriteData.
func (m *Monitor) CreateTimeSeriesRequest(tsList []*SDK.TimeSeries) *SDK.CreateTimeSeriesRequest {
	return &SDK.CreateTimeSeriesRequest{
		TimeSeries:      tsList,
		ForceSendFields: m.commonForceFields,
		NullFields:      m.commonNullFields,
	}
}

func (m *Monitor) buildTimeSeriesList(data Data) ([]*SDK.TimeSeries, error) {
	switch {
	case data.Data == nil:
		return nil, fmt.Errorf("error: data is required")
	case data.Resource == nil && m.commonResource == nil:
		return nil, fmt.Errorf("error: resource is required")
	}

	return data.TimeSeriesList(m.commonResource)
}

func (m *Monitor) formatProjectName() string {
	return fmt.Sprintf("projects/%s", m.projectID)
}

// Errorf logging error information.
func (m *Monitor) Errorf(format string, v ...interface{}) {
	m.logger.Errorf(serviceName, format, v...)
}
