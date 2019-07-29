package calendar

import "time"

// EventListOption contains option for event.list operation.
type EventListOption struct {
	MaxResults            int64
	OrderBy               OrderBy
	PageToken             string
	ShowDeleted           bool
	ShowHiddenInvitations bool
	SingleEvents          bool
	TimeMax               time.Time
	TimeMin               time.Time
}

func (o EventListOption) GetTimeMax() string {
	return getTime(o.TimeMax)
}

func (o EventListOption) GetTimeMin() string {
	return getTime(o.TimeMin)
}

func getTime(dt time.Time) string {
	return dt.Format(time.RFC3339)
}

func (o EventListOption) hasMaxResults() bool {
	return o.MaxResults > 0
}

func (o EventListOption) hasOrderBy() bool {
	return o.OrderBy != ""
}

func (o EventListOption) hasPageToken() bool {
	return o.PageToken != ""
}

func (o EventListOption) hasTimeMax() bool {
	return !o.TimeMax.IsZero()
}

func (o EventListOption) hasTimeMin() bool {
	return !o.TimeMin.IsZero()
}

type OrderBy string

// OrderBy list
const (
	OrderByStartTime OrderBy = "startTime" // used only SingleEvent=true
	OrderByUpdated   OrderBy = "updated"
)

func (f OrderBy) String() string {
	return string(f)
}
