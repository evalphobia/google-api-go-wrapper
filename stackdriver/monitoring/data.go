package monitoring

import (
	"fmt"
	"time"

	SDK "google.golang.org/api/monitoring/v3"
)

// Data is data struct for SDK.TimeSeries data.
type Data struct {
	// for SDK.Metric
	MetricGroup string
	MetricType  string
	Labels      map[string]string

	// for SDK.Point
	Data      interface{} // required
	EndTime   time.Time   // required
	StartTime time.Time

	MetricKind      MetricKind
	ValueType       ValueType
	Resource        *Resource
	ForceSendFields []string
	NullFields      []string
}

// TimeSeriesList converts Data to TimeSeries.
func (d *Data) TimeSeriesList(commonResource ...*Resource) ([]*SDK.TimeSeries, error) {
	switch v := d.Data.(type) {
	case []*SDK.TimeSeries:
		return v, nil
	case *SDK.TimeSeries:
		return []*SDK.TimeSeries{v}, nil
	case SDK.TimeSeries:
		return []*SDK.TimeSeries{&v}, nil
	}

	point, err := d.createPoint()
	if err != nil {
		return nil, err
	}

	ent := &SDK.TimeSeries{
		Metric: &SDK.Metric{
			Type:   d.formatMetricType(),
			Labels: d.Labels,
		},
		Points:          []*SDK.Point{point},
		MetricKind:      string(d.MetricKind),
		ValueType:       string(d.ValueType),
		ForceSendFields: d.ForceSendFields,
		NullFields:      d.NullFields,
	}

	switch {
	case d.Resource != nil:
		ent.Resource = d.Resource.toMonitorResource()
	case len(commonResource) != 0:
		ent.Resource = commonResource[0].toMonitorResource()
	}

	return []*SDK.TimeSeries{ent}, nil
}

func (d *Data) createPoint() (*SDK.Point, error) {
	// create TypedValue
	value := &SDK.TypedValue{}
	switch v := d.Data.(type) {
	case bool:
		value.BoolValue = &v
	case string:
		value.StringValue = &v
	case float64:
		value.DoubleValue = &v
	case float32:
		vv := float64(v)
		value.DoubleValue = &vv
	case int64:
		value.Int64Value = &v
	case int, int32, uint, uint8, uint16, uint32, uint64:
		vv := toInt64(v)
		value.Int64Value = &vv
	case *SDK.Distribution:
		value.DistributionValue = v
	case SDK.Distribution:
		value.DistributionValue = &v
	default:
		return nil, fmt.Errorf("unknown type value; type=%T", d.Data)
	}

	// create Interval
	if d.EndTime.IsZero() {
		d.EndTime = time.Now()
	}
	interval := &SDK.TimeInterval{
		EndTime: d.EndTime.Format(time.RFC3339),
	}
	if !d.StartTime.IsZero() {
		interval.StartTime = d.StartTime.Format(time.RFC3339)
	}

	return &SDK.Point{
		Interval: interval,
		Value:    value,
	}, nil
}

func toInt64(v interface{}) int64 {
	switch vv := v.(type) {
	case int:
		return int64(vv)
	case int32:
		return int64(vv)
	case uint:
		return int64(vv)
	case uint32:
		return int64(vv)
	case uint64:
		return int64(vv)
	case uint8:
		return int64(vv)
	case uint16:
		return int64(vv)
	default:
		return 0
	}
}

func (d *Data) formatMetricType() string {
	if d.MetricGroup != "" {
		return fmt.Sprintf("custom.googleapis.com/%s/%s", d.MetricGroup, d.MetricType)
	}
	return fmt.Sprintf("custom.googleapis.com/%s", d.MetricType)
}
