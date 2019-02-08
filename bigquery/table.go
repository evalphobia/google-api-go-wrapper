package bigquery

import (
	"net/http"

	SDK "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"

	"github.com/evalphobia/google-api-go-wrapper/config"
)

// TableAPI is bigquery table client.
type TableAPI struct {
	dataset *DatasetAPI
	tableID string
}

// NewTableAPI returns initialized TableAPI
func NewTableAPI(conf config.Config, projectID, datasetID, tableID string) (*TableAPI, error) {
	ds, err := NewDatasetAPI(conf, projectID, datasetID)
	if err != nil {
		return nil, err
	}

	tbl := &TableAPI{
		dataset: ds,
		tableID: tableID,
	}
	return tbl, nil
}

// Create creates the table with schema defined from given struct.
func (t *TableAPI) Create(schemaStruct interface{}) error {
	schema, err := convertToSchema(schemaStruct)
	if err != nil {
		return err
	}

	cli := t.dataset.client
	tbl := &SDK.Table{
		Schema: schema,
		TableReference: &SDK.TableReference{
			ProjectId: cli.projectID,
			DatasetId: t.dataset.datasetID,
			TableId:   t.tableID,
		},
	}

	_, err = cli.CreateTable(t.tableID, tbl)
	return err
}

// Get gets the table.
func (t *TableAPI) Get() (*Table, error) {
	cli := t.dataset.client
	return cli.GetTable(t.dataset.datasetID, t.tableID)
}

// IsExist checks if the table is exists in BQ.
func (t *TableAPI) IsExist() (bool, error) {
	_, err := t.Get()
	if err == nil {
		// table exists.
		return true, nil
	}

	apiErr, ok := err.(*googleapi.Error)
	switch {
	case !ok:
		// unknown error
		return false, err
	case apiErr.Code == http.StatusNotFound:
		// table does not exist.
		return false, nil
	default:
		// other API error
		return false, err
	}
}

// Drop deletes the table.
func (t *TableAPI) Drop() error {
	cli := t.dataset.client
	return cli.DropTable(t.dataset.datasetID, t.tableID)
}

// InsertAll appends all of map data by using InsertAll api
func (t *TableAPI) InsertAll(data interface{}) error {
	rows, err := buildTableDataInsertAllRequest(data)
	if err != nil {
		return err
	}

	cli := t.dataset.client
	resp, err := cli.InsertAll(t.dataset.datasetID, t.tableID, rows)
	switch {
	case err != nil:
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
