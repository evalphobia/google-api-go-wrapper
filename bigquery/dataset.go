package bigquery

import (
	"github.com/evalphobia/google-api-go-wrapper/config"
)

// DatasetAPI is bigquery dataset client.
type DatasetAPI struct {
	client    *BigQuery
	datasetID string
}

// NewDatasetAPI returns initialized DatasetAPI
func NewDatasetAPI(conf config.Config, projectID, datasetID string) (*DatasetAPI, error) {
	cli, err := New(conf, projectID)
	if err != nil {
		return nil, err
	}

	ds := &DatasetAPI{
		client:    cli,
		datasetID: datasetID,
	}
	return ds, nil
}

// TableAPI returns initialized TableAPI.
func (ds *DatasetAPI) TableAPI(tableID string) *TableAPI {
	return &TableAPI{
		dataset: ds,
		tableID: tableID,
	}
}

// Get gets the dataset.
func (ds *DatasetAPI) Get() (*Dataset, error) {
	cli := ds.client
	return cli.GetDataset(ds.datasetID)
}

// Delete deletes the dataset.
func (ds *DatasetAPI) Delete() error {
	cli := ds.client
	return cli.DeleteDataset(ds.datasetID)
}

// CreateTable creates the table with schema defined from given struct
// (*Deprecated)
func (ds *DatasetAPI) CreateTable(tableID string, schemaStruct interface{}) error {
	return ds.TableAPI(tableID).Create(schemaStruct)
}

// InsertAll appends all of map data by using InsertAll api
// (*Deprecated)
func (ds *DatasetAPI) InsertAll(tableID string, data interface{}) error {
	return ds.TableAPI(tableID).InsertAll(data)
}
