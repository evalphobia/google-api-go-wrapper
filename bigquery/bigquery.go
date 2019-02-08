package bigquery

import (
	"errors"

	SDK "google.golang.org/api/bigquery/v2"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "BigQuery"
	scope       = SDK.BigqueryScope
)

var (
	errOperationInsertAll = errors.New("error occured on bigquery.InsertAll")
	errDataType           = errors.New("error data type")
)

// BigQuery is BigQuery API client..
type BigQuery struct {
	service   *SDK.Service
	logger    log.Logger
	projectID string
}

// New returns initialized BigQuery.
func New(conf config.Config, projectID string) (*BigQuery, error) {
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

	b := &BigQuery{
		service:   svc,
		logger:    log.DefaultLogger,
		projectID: projectID,
	}
	return b, nil
}

// SetLogger sets internal API logger.
func (b *BigQuery) SetLogger(logger log.Logger) {
	b.logger = logger
}

// Errorf logging error information.
func (b *BigQuery) Errorf(format string, v ...interface{}) {
	b.logger.Errorf(serviceName, format, v...)
}

// DatasetAPI returns initialized DatasetAPI.
func (b *BigQuery) DatasetAPI(datasetID string) *DatasetAPI {
	return &DatasetAPI{
		client:    b,
		datasetID: datasetID,
	}
}
