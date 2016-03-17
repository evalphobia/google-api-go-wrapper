package bigquery

import (
	"errors"
	"fmt"

	bigquery "google.golang.org/api/bigquery/v2"
)

var (
	errOperationCreateTable = errors.New("error occured on bigquery.Insert")
	errOperationInsertAll   = errors.New("error occured on bigquery.InsertAll")
	errDataType             = errors.New("error data type")
)

// Dataset repesents bigquery dataset
type Dataset struct {
	service   *bigquery.Service
	projectID string
	datasetID string
}

// NewDataset returns initialized Dataset
func NewDataset(cli Client, projectID, datasetID string) (*Dataset, error) {
	svc, err := bigquery.New(cli.httpClient)
	if err != nil {
		return nil, err
	}

	ds := &Dataset{
		service:   svc,
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
	tbl := &bigquery.Table{
		Schema: schema,
		TableReference: &bigquery.TableReference{
			ProjectId: ds.projectID,
			DatasetId: ds.datasetID,
			TableId:   tableID,
		},
	}
	resp, err := ds.service.Tables.Insert(ds.projectID, ds.datasetID, tbl).Do()
	if err != nil {
		fmt.Printf("err: %+v\n", err.Error())
		return err
	}
	_ = resp

	return nil
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
		return err
	case len(resp.InsertErrors) != 0:
		return errOperationInsertAll
	}

	return nil
}

func buildTableDataInsertAllRequest(data interface{}) (*bigquery.TableDataInsertAllRequest, error) {
	switch v := data.(type) {
	case []map[string]interface{}:
		return buildRowsFromMaps(v), nil
	case map[string]interface{}:
		return &bigquery.TableDataInsertAllRequest{
			Rows: []*bigquery.TableDataInsertAllRequestRows{buildRowFromMap(v)},
		}, nil
	}

	if list, ok := getSliceData(data); ok {
		rows := make([]*bigquery.TableDataInsertAllRequestRows, len(list))
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
		return &bigquery.TableDataInsertAllRequest{
			Rows: rows,
		}, nil
	}

	if isStruct(data) {
		row, err := buildRowsFromStruct(data)
		return &bigquery.TableDataInsertAllRequest{
			Rows: []*bigquery.TableDataInsertAllRequestRows{row},
		}, err
	}

	return nil, errDataType
}

func buildRowsFromMaps(list []map[string]interface{}) *bigquery.TableDataInsertAllRequest {
	rows := make([]*bigquery.TableDataInsertAllRequestRows, len(list))
	for i, row := range list {
		rows[i] = buildRowFromMap(row)
	}

	return &bigquery.TableDataInsertAllRequest{
		Rows: rows,
	}
}

func buildRowFromMap(row map[string]interface{}) *bigquery.TableDataInsertAllRequestRows {
	return &bigquery.TableDataInsertAllRequestRows{
		Json: buildJSONValue(row),
	}
}

func buildJSONValue(row map[string]interface{}) map[string]bigquery.JsonValue {
	jsonValue := make(map[string]bigquery.JsonValue)
	for k, v := range row {
		jsonValue[k] = v
	}
	return jsonValue
}

func buildRowsFromStruct(data interface{}) (*bigquery.TableDataInsertAllRequestRows, error) {
	row, err := convertStructToMap(data)
	if err != nil {
		return nil, err
	}
	return buildRowFromMap(row), nil
}
