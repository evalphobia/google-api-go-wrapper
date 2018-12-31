package analytics

import (
	"strconv"

	SDK "google.golang.org/api/analytics/v3"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	scope            = SDK.AnalyticsReadonlyScope
	metricActiveUser = "rt:activeUsers"
)

// Analytics is wrapper struct for Google Analytics.
type Analytics struct {
	service *SDK.Service
	logger  log.Logger
}

// New returns initialized Analytics.
func New(conf config.Config) (*Analytics, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = append(conf.Scopes, scope)
	}
	cli, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := SDK.New(cli)
	if err != nil {
		return nil, err
	}

	ds := &Analytics{
		service: svc,
		logger:  log.DefaultLogger,
	}
	return ds, nil
}

// SetLogger sets internal API logger.
func (a *Analytics) SetLogger(logger log.Logger) {
	a.logger = logger
}

// GetRealtime gets *SDK.RealtimeData.
func (a *Analytics) GetRealtime(id string) (*SDK.RealtimeData, error) {
	rs := SDK.NewDataRealtimeService(a.service)
	return rs.Get("ga:"+id, metricActiveUser).Do()
}

// GetRealtimeActiveUser gets active user.
func (a *Analytics) GetRealtimeActiveUser(id string) (int64, error) {
	data, err := a.GetRealtime(id)
	if err != nil {
		return 0, err
	}

	v, ok := data.TotalsForAllResults[metricActiveUser]
	if !ok {
		return 0, nil
	}

	result, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return int64(result), nil
}
