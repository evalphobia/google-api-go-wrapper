package bigquery

import (
	SDK "google.golang.org/api/bigquery/v2"
)

// Table represents SDK.Table.
type Table struct {
	*SDK.Table
}

// Dataset represents SDK.Dataset.
type Dataset struct {
	*SDK.Dataset
}
