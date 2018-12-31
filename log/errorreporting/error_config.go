package errorreporting

import (
	"cloud.google.com/go/errorreporting"
)

// ErrorConfig is config for Google Stackdriver error repoting client.
type ErrorConfig struct {
	ProjectID      string
	ServiceName    string
	ServiceVersion string
	OnError        func(err error)
	UseSync        bool
}

func (c ErrorConfig) Config() errorreporting.Config {
	return errorreporting.Config{
		ServiceName:    c.ServiceName,
		ServiceVersion: c.ServiceVersion,
		OnError:        c.OnError,
	}
}
