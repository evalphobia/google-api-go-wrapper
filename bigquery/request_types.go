package bigquery

import (
	SDK "google.golang.org/api/bigquery/v2"
)

type QueryOption struct {
	// required
	SQL string

	// optional
	ProjectID     string
	DatasetID     string
	DryRun        bool
	Location      string
	MaxResults    int64
	ParameterMode string
	TimeoutMs     int64
	UseLegacySql  bool
	NoQueryCache  bool
}

func (o QueryOption) ToRequest() *SDK.QueryRequest {
	in := &SDK.QueryRequest{
		Query:         o.SQL,
		Location:      o.Location,
		MaxResults:    o.MaxResults,
		ParameterMode: o.ParameterMode,
		TimeoutMs:     o.TimeoutMs,
		UseLegacySql:  &o.UseLegacySql,
	}

	if o.ProjectID != "" && o.DatasetID != "" {
		in.DefaultDataset = &SDK.DatasetReference{
			ProjectId: o.ProjectID,
			DatasetId: o.DatasetID,
		}
	}

	if o.NoQueryCache {
		in.UseQueryCache = &o.NoQueryCache
	}
	return in
}
