package calendar

import (
	"time"

	SDK "google.golang.org/api/calendar/v3"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log"
)

const (
	serviceName = "calendar"
)

// Calendar repesents GSUITE Calendar API client.
type Calendar struct {
	service *SDK.Service
	logger  log.Logger
}

// New returns initialized *Calendar.
func New(conf config.Config) (*Calendar, error) {
	if len(conf.Scopes) == 0 {
		conf.Scopes = append(conf.Scopes, SDK.CalendarScope)
	}
	cli, err := conf.Client()
	if err != nil {
		return nil, err
	}

	svc, err := SDK.New(cli)
	if err != nil {
		return nil, err
	}

	Calendar := &Calendar{
		service: svc,
		logger:  log.DefaultLogger,
	}
	return Calendar, nil
}

// SetLogger sets internal API logger.
func (c *Calendar) SetLogger(logger log.Logger) {
	c.logger = logger
}

// EventList gets calendarID's events after current time.
func (c *Calendar) EventList(calendarID string) (*EventList, error) {
	resp, err := c.eventList(calendarID, EventListOption{
		TimeMin:      time.Now(),
		SingleEvents: true,
		OrderBy:      OrderByStartTime,
	})
	if err != nil {
		return nil, err
	}
	return &EventList{
		List: NewEvents(resp.Items),
	}, nil
}

// eventList executes events.list operation.
func (c *Calendar) eventList(calendarID string, opt EventListOption) (*EventListResponse, error) {
	listCall := c.service.Events.List(calendarID)
	if opt.hasMaxResults() {
		listCall.MaxResults(opt.MaxResults)
	}
	if opt.hasOrderBy() {
		listCall.OrderBy(opt.OrderBy.String())
	}
	if opt.hasPageToken() {
		listCall.PageToken(opt.PageToken)
	}
	if opt.ShowDeleted {
		listCall.ShowDeleted(opt.ShowDeleted)
	}
	if opt.ShowHiddenInvitations {
		listCall.ShowHiddenInvitations(opt.ShowHiddenInvitations)
	}
	if opt.SingleEvents {
		listCall.SingleEvents(opt.SingleEvents)
	}
	if opt.hasTimeMax() {
		listCall.TimeMax(opt.GetTimeMax())
	}
	if opt.hasTimeMin() {
		listCall.TimeMin(opt.GetTimeMin())
	}

	resp, err := listCall.Do()
	if err != nil {
		c.Errorf("error on `Events.List` operation;  error=%s", err.Error())
	}
	return &EventListResponse{resp}, err
}

// Errorf logging error information.
func (c *Calendar) Errorf(format string, vv ...interface{}) {
	c.logger.Errorf(serviceName, format, vv...)
}
