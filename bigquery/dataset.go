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

// Dataset repesents bigquery dataset.
type Dataset struct {
	service   *SDK.Service
	logger    log.Logger
	projectID string
	datasetID string
}

// NewDataset returns initialized Dataset
func NewDataset(conf config.Config, projectID, datasetID string) (*Dataset, error) {
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

	ds := &Dataset{
		service:   svc,
		logger:    log.DefaultLogger,
		projectID: projectID,
		datasetID: datasetID,
	}
	return ds, nil
}

// CreateTable creates the table with schema defined from given struct
func (ds *Dataset) CreateTable(tableID string, schemaStruct interface{}) error {
	schema, err := convertToSchema(schemaStruct)
	if err != nil {
		return err
	}
	tbl := &SDK.Table{
		Schema: schema,
		TableReference: &SDK.TableReference{
			ProjectId: ds.projectID,
			DatasetId: ds.datasetID,
			TableId:   tableID,
		},
	}

	_, err = ds.service.Tables.Insert(ds.projectID, ds.datasetID, tbl).Do()
	if err != nil {
		ds.Errorf("error on `Insert` operation; projectID=%s, datasetID=%s, error=%s;", ds.projectID, ds.datasetID, err.Error())
	}
	return err
}

// InsertAll appends all of map data by using InsertAll api
func (ds *Dataset) InsertAll(tableID string, data interface{}) error {
	rows, err := buildTableDataInsertAllRequest(data)
	if err != nil {
		return err
	}

	resp, err := ds.service.Tabledata.InsertAll(ds.projectID, ds.datasetID, tableID, rows).Do()
	switch {
	case err != nil:
		ds.Errorf("error on `InsertAll` operation; projectID=%s, datasetID=%s, table=%s, error=%s;", ds.projectID, ds.datasetID, tableID, err.Error())
		return err
	case len(resp.InsertErrors) != 0:
		return errOperationInsertAll
	}

	return nil
}

func buildTableDataInsertAllRequest(data interface{}) (*SDK.TableDataInsertAllRequest, error) {
	switch v := data.(type) {
	case []map[string]interface{}:
		return buildRowsFromMaps(v), nil
	case map[string]interface{}:
		return &SDK.TableDataInsertAllRequest{
			Rows: []*SDK.TableDataInsertAllRequestRows{buildRowFromMap(v)},
		}, nil
	}

	if list, ok := getSliceData(data); ok {
		rows := make([]*SDK.TableDataInsertAllRequestRows, len(list))
		for i, v := range list {
			if !isStruct(v) {
				continue
			}

			row, err := buildRowsFromStruct(v)
			if err != nil {
				return nil, err
			}

			rows[i] = row
		}
		return &SDK.TableDataInsertAllRequest{
			Rows: rows,
		}, nil
	}

	if isStruct(data) {
		row, err := buildRowsFromStruct(data)
		return &SDK.TableDataInsertAllRequest{
			Rows: []*SDK.TableDataInsertAllRequestRows{row},
		}, err
	}

	return nil, errDataType
}

func buildRowsFromMaps(list []map[string]interface{}) *SDK.TableDataInsertAllRequest {
	rows := make([]*SDK.TableDataInsertAllRequestRows, len(list))
	for i, row := range list {
		rows[i] = buildRowFromMap(row)
	}

	return &SDK.TableDataInsertAllRequest{
		Rows: rows,
	}
}

func buildRowFromMap(row map[string]interface{}) *SDK.TableDataInsertAllRequestRows {
	return &SDK.TableDataInsertAllRequestRows{
		Json: buildJSONValue(row),
	}
}

func buildJSONValue(row map[string]interface{}) map[string]SDK.JsonValue {
	jsonValue := make(map[string]SDK.JsonValue)
	for k, v := range row {
		jsonValue[k] = v
	}
	return jsonValue
}

func buildRowsFromStruct(data interface{}) (*SDK.TableDataInsertAllRequestRows, error) {
	row, err := convertStructToMap(data)
	if err != nil {
		return nil, err
	}
	return buildRowFromMap(row), nil
}

// Errorf logging error information.
func (ds *Dataset) Errorf(format string, v ...interface{}) {
	ds.logger.Errorf(serviceName, format, v...)
}
