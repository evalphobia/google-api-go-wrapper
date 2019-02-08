package bigquery

import (
	"fmt"
	"strings"

	SDK "google.golang.org/api/bigquery/v2"
)

// see API documents: https://cloud.google.com/bigquery/docs/reference/rest/v2

// ==========
// Datasets
// ==========

// CreateDataset performes Datasets.Insert operation.
// Creates a new empty dataset.
func (b *BigQuery) CreateDataset(dataset *SDK.Dataset) (*Dataset, error) {
	ds, err := b.service.Datasets.Insert(b.projectID, dataset).Do()
	b.logAPIError("Datasets.Insert", err, logArgs("datasetID", dataset.Id))
	return &Dataset{ds}, err
}

// PatchDataset performes Datasets.Patch operation.
// Updates information in an existing dataset. The update method replaces the entire dataset resource, whereas the patch method only replaces fields that are provided in the submitted dataset resource. This method supports patch semantics.
func (b *BigQuery) PatchDataset(datasetID string, dataset *SDK.Dataset) (*Dataset, error) {
	ds, err := b.service.Datasets.Patch(b.projectID, datasetID, dataset).Do()
	b.logAPIError("Datasets.Patch", err, logArgs("datasetID", datasetID))
	return &Dataset{ds}, err
}

// UpdateDataset performes Datasets.Update operation.
// Updates information in an existing dataset. The update method replaces the entire dataset resource, whereas the patch method only replaces fields that are provided in the submitted dataset resource.
func (b *BigQuery) UpdateDataset(datasetID string, dataset *SDK.Dataset) (*Dataset, error) {
	ds, err := b.service.Datasets.Update(b.projectID, datasetID, dataset).Do()
	b.logAPIError("Datasets.Update", err, logArgs("datasetID", datasetID))
	return &Dataset{ds}, err
}

// DeleteDataset performes Datasets.Delete operation.
// Deletes the dataset specified by the datasetId value. Before you can delete a dataset, you must delete all its tables, either manually or by specifying deleteContents. Immediately after deletion, you can create another dataset with the same name.
func (b *BigQuery) DeleteDataset(datasetID string) error {
	err := b.service.Datasets.Delete(b.projectID, datasetID).Do()
	b.logAPIError("Datasets.Delete", err, logArgs("datasetID", datasetID))
	return err
}

// GetDataset performes Datasets.Get operation.
// Returns the dataset specified by datasetID.
func (b *BigQuery) GetDataset(datasetID string) (*Dataset, error) {
	ds, err := b.service.Datasets.Get(b.projectID, datasetID).Do()
	b.logAPIError("Datasets.Get", err, logArgs("datasetID", datasetID))
	return &Dataset{ds}, err
}

// ListDatasets performes Datasets.List operation.
// Lists all datasets in the specified project to which you have been granted the READER dataset role.
func (b *BigQuery) ListDatasets() (*SDK.DatasetList, error) {
	list, err := b.service.Datasets.List(b.projectID).Do()
	b.logAPIError("Datasets.List", err)
	return list, err
}

// ==========
// Jobs
// ==========

// RunJob performes Jobs.Insert operation.
// Starts a new asynchronous job. Requires the Can View project role.
func (b *BigQuery) RunJob(job *SDK.Job) (*SDK.Job, error) {
	j, err := b.service.Jobs.Insert(b.projectID, job).Do()
	b.logAPIError("Jobs.Insert", err)
	return j, err
}

// RunQuery performes Jobs.Query operation.
// Runs a BigQuery SQL query and returns results if the query completes within a specified timeout.
func (b *BigQuery) RunQuery(query *SDK.QueryRequest) (*SDK.QueryResponse, error) {
	resp, err := b.service.Jobs.Query(b.projectID, query).Do()
	b.logAPIError("Jobs.Query", err)
	return resp, err
}

// CancelJob performes Jobs.Cancel operation.
// Requests that a job be cancelled. This call will return immediately, and the client will need to poll for the job status to see if the cancel completed successfully. Cancelled jobs may still incur costs. For more information, see pricing.
func (b *BigQuery) CancelJob(jobID string) (*SDK.JobCancelResponse, error) {
	resp, err := b.service.Jobs.Cancel(b.projectID, jobID).Do()
	b.logAPIError("Jobs.Cancel", err, logArgs("jobID", jobID))
	return resp, err
}

// GetJob performes Jobs.Get operation.
// Returns information about a specific job. Job information is available for a six month period after creation. Requires that you're the person who ran the job, or have the Is Owner project role.
func (b *BigQuery) GetJob(jobID string) (*SDK.Job, error) {
	j, err := b.service.Jobs.Get(b.projectID, jobID).Do()
	b.logAPIError("Jobs.Get", err, logArgs("jobID", jobID))
	return j, err
}

// ListJobs performes Jobs.List operation.
// Lists all jobs that you started in the specified project. Job information is available for a six month period after creation. The job list is sorted in reverse chronological order, by job creation time. Requires the Can View project role, or the Is Owner project role if you set the allUsers property.
func (b *BigQuery) ListJobs() (*SDK.JobList, error) {
	list, err := b.service.Jobs.List(b.projectID).Do()
	b.logAPIError("Jobs.List", err)
	return list, err
}

// GetQueryResults performes Jobs.GetQueryResults operation.
// Retrieves the results of a query job.
func (b *BigQuery) GetQueryResults(jobID string) (*SDK.GetQueryResultsResponse, error) {
	resp, err := b.service.Jobs.GetQueryResults(b.projectID, jobID).Do()
	b.logAPIError("Jobs.GetQueryResults", err, logArgs("jobID", jobID))
	return resp, err
}

// ==========
// Tabledata
// ==========

// InsertAll performes Tabledata.InsertAll operation.
// Streams data into BigQuery one record at a time without needing to run a load job. For more information, see streaming data into BigQuery.
func (b *BigQuery) InsertAll(datasetID string, tableID string, rows *SDK.TableDataInsertAllRequest) (*SDK.TableDataInsertAllResponse, error) {
	resp, err := b.service.Tabledata.InsertAll(b.projectID, datasetID, tableID, rows).Do()
	b.logAPIError("Tabledata.InsertAll", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return resp, err
}

// GetTableData performes Tabledata.List operation.
// Retrieves table data from a specified set of rows. Requires the READER dataset role.
func (b *BigQuery) GetTableData(datasetID string, tableID string) (*SDK.TableDataList, error) {
	list, err := b.service.Tabledata.List(b.projectID, datasetID, tableID).Do()
	b.logAPIError("Tabledata.List", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return list, err
}

// ==========
// Table
// ==========

// CreateTable performes Table.Insert operation.
// Creates a new, empty table in the dataset.
func (b *BigQuery) CreateTable(datasetID string, tbl *SDK.Table) (*Table, error) {
	t, err := b.service.Tables.Insert(b.projectID, datasetID, tbl).Do()
	b.logAPIError("Table.Insert", err, logArgs("datasetID", datasetID))
	return &Table{t}, err
}

// PatchTable performes Tables.Patch operation.
// Updates information in an existing table. The update method replaces the entire table resource, whereas the patch method only replaces fields that are provided in the submitted table resource. This method supports patch semantics.
func (b *BigQuery) PatchTable(datasetID string, tableID string, tbl *SDK.Table) (*Table, error) {
	t, err := b.service.Tables.Patch(b.projectID, datasetID, tableID, tbl).Do()
	b.logAPIError("Table.Patch", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return &Table{t}, err
}

// UpdateTable performes Tables.Update operation.
// Updates information in an existing table. The update method replaces the entire table resource, whereas the patch method only replaces fields that are provided in the submitted table resource.
func (b *BigQuery) UpdateTable(datasetID string, tableID string, tbl *SDK.Table) (*Table, error) {
	t, err := b.service.Tables.Update(b.projectID, datasetID, tableID, tbl).Do()
	b.logAPIError("Table.Update", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return &Table{t}, err
}

// DropTable performes Tables.Delete operation.
// Deletes the table specified by tableId from the dataset. If the table contains data, all the data will be deleted.
func (b *BigQuery) DropTable(datasetID string, tableID string) error {
	err := b.service.Tables.Delete(b.projectID, datasetID, tableID).Do()
	b.logAPIError("Table.Delete", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return err
}

// GetTable performes Tables.Get operation.
// Gets the specified table resource by table ID. This method does not return the data in the table, it only returns the table resource, which describes the structure of this table.
func (b *BigQuery) GetTable(datasetID string, tableID string) (*Table, error) {
	t, err := b.service.Tables.Get(b.projectID, datasetID, tableID).Do()
	b.logAPIError("Table.Get", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return &Table{t}, err
}

// ListTables performes Tables.List operation.
// Lists all tables in the specified dataset. Requires the READER dataset role.
func (b *BigQuery) ListTables(datasetID string, tableID string) (*SDK.TableList, error) {
	list, err := b.service.Tables.List(b.projectID, datasetID).Do()
	b.logAPIError("Table.List", err, logArgs("datasetID", datasetID), logArgs("tableID", tableID))
	return list, err
}

func (b *BigQuery) logAPIError(apiName string, err error, opts ...string) {
	if err == nil {
		return
	}

	msg := fmt.Sprintf("error on `%s` operation; error=[%s] projectID=[%s],", apiName, err.Error(), b.projectID)
	if len(opts) != 0 {
		msg = fmt.Sprintf("%s %s", msg, strings.Join(opts, " "))
	}
	b.Errorf(msg)
}

func logArgs(key, value string) string {
	return fmt.Sprintf("%s=[%s]", key, value)
}
