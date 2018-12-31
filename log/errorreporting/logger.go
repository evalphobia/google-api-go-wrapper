package errorreporting

import (
	"context"
	"fmt"

	"cloud.google.com/go/errorreporting"
	"google.golang.org/api/option"

	"github.com/evalphobia/google-api-go-wrapper/config"
)

// Logger use Google Stackdriver error repoting.
type Logger struct {
	errorClient *errorreporting.Client
	useSync     bool

	// Some errors seem that they does not handled by errorreporting's `OnError`. (e.g. creds error)
	// Set `OnError` then it handles, only when `useSync=true`.
	onError func(err error)
}

// New creates initialized Logger.
func New(ctx context.Context, errConfig ErrorConfig) (*Logger, error) {
	errorClient, err := errorreporting.NewClient(ctx, errConfig.ProjectID, errConfig.Config())
	if err != nil {
		return nil, err
	}

	return &Logger{
		errorClient: errorClient,
		useSync:     errConfig.UseSync,
	}, nil
}

// NewWithConfig creates initialized Logger.
func NewWithConfig(ctx context.Context, errConfig ErrorConfig, conf config.Config) (*Logger, error) {
	credsFile, err := conf.CredsFilePath()
	if err != nil {
		return nil, err
	}

	credsOpt := option.WithCredentialsFile(credsFile)
	errorClient, err := errorreporting.NewClient(ctx, errConfig.ProjectID, errConfig.Config(), credsOpt)
	if err != nil {
		return nil, err
	}
	err = conf.DeleteTempCredsFile()
	if err != nil {
		return nil, err
	}

	return &Logger{
		errorClient: errorClient,
		useSync:     errConfig.UseSync,
	}, nil
}

// SetOnError sets OnError.
func (l *Logger) SetOnError(onErr func(error)) {
	l.onError = onErr
}

// Infof logging information.
func (*Logger) Infof(service, format string, v ...interface{}) {
	// do nothing
}

// Errorf logging error information.
func (l *Logger) Errorf(service, format string, v ...interface{}) {
	errMsg := fmt.Errorf("[%s] %s", service, fmt.Sprintf(format, v...))
	switch {
	case l.useSync:
		err := l.errorClient.ReportSync(context.TODO(), errorreporting.Entry{
			Error: errMsg,
		})
		if err != nil && l.onError != nil {
			l.onError(err)
		}
	default:
		l.errorClient.Report(errorreporting.Entry{
			Error: errMsg,
		})
	}
}

// Flush sends buffered logs.
func (l *Logger) Flush() {
	l.errorClient.Flush()
}
