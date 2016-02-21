package bigquery

import (
	"errors"
	"fmt"

	bigquery "google.golang.org/api/bigquery/v2"
)

var (
	errOperationCreateTable = errors.New("error occured on bigquery.Insert")
	errOperationInsertAll   = errors.New("error occured on bigquery.InsertAll")
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
	fmt.Printf("schema: %#v\n", schema.Fields)

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
func (ds *Dataset) InsertAll(tableID string, rows []map[string]interface{}) error {
	resp, err := ds.service.Tabledata.InsertAll(ds.projectID, ds.datasetID, tableID, buildRowsFromMaps(rows)).Do()
	switch {
	case err != nil:
		return err
	case len(resp.InsertErrors) != 0:
		return errOperationInsertAll
	}

	return nil
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
