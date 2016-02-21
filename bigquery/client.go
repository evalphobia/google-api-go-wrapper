package bigquery

import (
	"net/http"
	"time"

	bigquery "google.golang.org/api/bigquery/v2"

	"github.com/evalphobia/google-api-go-wrapper/library"
)

const (
	defaultTimeout = time.Second * 60

	scope = bigquery.BigqueryScope
)

var timeout = defaultTimeout

// SetTimeout sets HTTP timeout
func SetTimeout(t time.Duration) {
	timeout = t
}

// Client is wrapper of http.Client
type Client struct {
	httpClient *http.Client
}

// New creates Client from json key data
func New(jsonKey []byte) (Client, error) {
	conf, err := library.NewAPIConfig(jsonKey, scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}

// NewWithPem creates Client from pem key filepath
func NewWithPem(pemPath string) (Client, error) {
	conf, err := library.NewAPIConfigWithPem(pemPath, scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}

// NewWithParams creates Client from key and email
func NewWithParams(key, email string) Client {
	conf := library.NewAPIConfigWithParams(key, email, scope)
	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}
}

// NewWithConfig creates Client from Config
func NewWithConfig(c library.Config) Client {
	c.Scopes = []string{scope}
	conf := library.NewAPIConfigWithConfig(c)
	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}
}

// NewWithEnv creates Client from environment variables
func NewWithEnv() (Client, error) {
	conf, err := library.NewAPIConfigWithEnv(scope)
	if err != nil {
		return Client{}, err
	}

	ctx := library.NewContext(timeout)
	return Client{conf.Client(ctx)}, nil
}
