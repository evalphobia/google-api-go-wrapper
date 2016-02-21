package analytics

import (
	"net/http"
	"strconv"
	"time"

	analytics "google.golang.org/api/analytics/v3"

	"github.com/evalphobia/google-api-go-wrapper/library"
)

const (
	defaultTimeout = time.Second * 10

	scope            = analytics.AnalyticsReadonlyScope
	metricActiveUser = "rt:activeUsers"
)

var timeout = defaultTimeout

func SetTimeout(t time.Duration) {
	timeout = t
}

type Client struct {
	httpClient *http.Client
}

func New(jsonKey []byte) (Client, error) {
	conf, err := library.NewAPIConfig(jsonKey, scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}

func NewWithPem(pemPath string) (Client, error) {
	conf, err := library.NewAPIConfigWithPem(pemPath, scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}

func NewWithParams(key, email string) Client {
	conf := library.NewAPIConfigWithParams(key, email, scope)
	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}
}

func NewWithConfig(c library.Config) Client {
	c.Scopes = []string{scope}
	conf := library.NewAPIConfigWithConfig(c)
	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}
}

func NewWithEnv() (Client, error) {
	conf, err := library.NewAPIConfigWithEnv(scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}

func (c *Client) GetRealtime(id string) (*analytics.RealtimeData, error) {
	service, err := analytics.New(c.httpClient)
	if err != nil {
		return nil, err
	}

	rs := analytics.NewDataRealtimeService(service)
	return rs.Get("ga:"+id, metricActiveUser).Do()
}

func (c *Client) GetRealtimeActiveUser(id string) (int64, error) {
	data, err := c.GetRealtime(id)
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
