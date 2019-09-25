package bigquery

import (
	"errors"
	"strconv"
	"time"

	SDK "google.golang.org/api/bigquery/v2"
)

type QueryResponse struct {
	*SDK.QueryResponse
}

func (r *QueryResponse) ToMap() ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, 0, len(r.Rows))
	if r.Schema == nil {
		return nil, errors.New("Schema is nil")
	}

	columnTypes := make([]columnType, len(r.Schema.Fields))
	for i, f := range r.Schema.Fields {
		columnTypes[i] = columnType{
			Name: f.Name,
			Type: f.Type,
		}
	}

	for _, r := range r.Rows {
		row := make(map[string]interface{})
		for i, col := range r.F {
			typ := columnTypes[i]
			typ.AssignData(row, col.V)
		}
		results = append(results, row)
	}

	return results, nil
}

type columnType struct {
	Name string
	Type string
}

func (c columnType) AssignData(row map[string]interface{}, value interface{}) {
	v := value.(string)
	switch {
	case c.IsString():
		row[c.Name] = v
	case c.IsInt():
		row[c.Name], _ = strconv.ParseInt(v, 10, 64)
	case c.IsDate():
		row[c.Name], _ = time.Parse("2006-01-02", v)
	}
}

func (c columnType) IsString() bool {
	return c.Type == "STRING"
}

func (c columnType) IsInt() bool {
	return c.Type == "INTEGER"
}

func (c columnType) IsDate() bool {
	return c.Type == "DATE"
}
