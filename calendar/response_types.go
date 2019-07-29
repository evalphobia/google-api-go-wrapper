package calendar

import (
	"time"

	SDK "google.golang.org/api/calendar/v3"
)

type EventListResponse struct {
	*SDK.Events
}

// EventList contains multiple calendar events.
type EventList struct {
	List []Event
}

// Event contains event information on Google Calendar.
type Event struct {
	ID          string
	Status      string
	Summary     string
	Description string
	Location    string
	HTMLLink    string
	HangoutLink string
	ICalUID     string

	StartDate time.Time
	EndDate   time.Time

	Attendees []User
	Creator   User
	Organizer User

	Transparent bool
	Visibility  string
	Locked      bool
	PrivateCopy bool
	Sequence    int64
	Kind        string

	Source    EventSource
	Reminders EventReminders

	Created time.Time
	Updated time.Time
}

func NewEvents(list []*SDK.Event) []Event {
	results := make([]Event, len(list))
	for i, e := range list {
		results[i] = NewEvent(e)
	}
	return results
}

func NewEvent(e *SDK.Event) Event {
	return Event{
		ID:          e.Id,
		Status:      e.Status,
		Summary:     e.Summary,
		Description: e.Description,
		Location:    e.Location,
		HTMLLink:    e.HtmlLink,
		HangoutLink: e.HangoutLink,
		ICalUID:     e.ICalUID,
		StartDate:   mustTimeFromDateTime(e.Start),
		EndDate:     mustTimeFromDateTime(e.End),
		Creator:     newUserFromCreator(e.Creator),
		Organizer:   newUserFromOrganizer(e.Organizer),
		Attendees:   newUserFromAttendees(e.Attendees),
		Transparent: e.Transparency == "tranparent",
		Visibility:  e.Visibility,
		Locked:      e.Locked,
		PrivateCopy: e.PrivateCopy,
		Sequence:    e.Sequence,
		Kind:        e.Kind,
		Reminders:   newEventReminders(e.Reminders),
		Source:      newEventSource(e.Source),
		Created:     mustTimeFromString(e.Created),
		Updated:     mustTimeFromString(e.Updated),
	}
}

func mustTimeFromString(str string) time.Time {
	dt, _ := time.Parse(str, time.RFC3339)
	return dt
}

func mustTimeFromDateTime(dt *SDK.EventDateTime) time.Time {
	if dt.Date != "" {
		t, _ := time.Parse(dt.Date, "2006-01-02")
		return t
	}

	t, _ := time.Parse(dt.DateTime, time.RFC3339)
	return t
}

type User struct {
	ID             string
	Email          string
	DisplayName    string
	ResponseStatus string
	Self           bool
	Organizer      bool
}

func newUserFromAttendees(list []*SDK.EventAttendee) []User {
	results := make([]User, len(list))
	for i, e := range list {
		results[i] = newUserFromAttendee(e)
	}
	return results
}

func newUserFromAttendee(e *SDK.EventAttendee) User {
	return User{
		ID:             e.Id,
		Email:          e.Email,
		DisplayName:    e.DisplayName,
		ResponseStatus: e.ResponseStatus,
		Self:           e.Self,
		Organizer:      e.Organizer,
	}
}

func newUserFromCreator(e *SDK.EventCreator) User {
	return User{
		ID:          e.Id,
		Email:       e.Email,
		DisplayName: e.DisplayName,
		Self:        e.Self,
	}
}

func newUserFromOrganizer(e *SDK.EventOrganizer) User {
	return User{
		ID:          e.Id,
		Email:       e.Email,
		DisplayName: e.DisplayName,
		Self:        e.Self,
	}
}

type EventSource struct {
	Title string
	URL   string
}

func newEventSource(s *SDK.EventSource) EventSource {
	if s == nil {
		return EventSource{}
	}

	return EventSource{
		Title: s.Title,
		URL:   s.Url,
	}
}

type EventReminders struct {
	Overrides  []EventReminder
	UseDefault bool
}

func newEventReminders(r *SDK.EventReminders) EventReminders {
	list := make([]EventReminder, len(r.Overrides))
	for i, e := range r.Overrides {
		list[i] = newEventReminder(e)
	}
	return EventReminders{
		Overrides:  list,
		UseDefault: r.UseDefault,
	}
}

type EventReminder struct {
	Method  string
	Minutes int64
}

func newEventReminder(e *SDK.EventReminder) EventReminder {
	return EventReminder{
		Method:  e.Method,
		Minutes: e.Minutes,
	}
}
