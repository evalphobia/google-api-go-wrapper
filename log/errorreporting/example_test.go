package errorreporting_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/log/errorreporting"
)

func ExampleNew(t *testing.T) {
	ctx := context.Background()
	logger, err := errorreporting.New(ctx, errorreporting.ErrorConfig{
		ProjectID:   "my-project-id",
		ServiceName: "my-service",
		UseSync:     true,
		OnError: func(err error) {
			fmt.Printf("[ERROR] err: %s", err.Error())
		},
	})
	if err != nil {
		panic(err)
	}

	logger.SetOnError(func(err error) {
		fmt.Printf("[ReportSync Error] err: %s", err.Error())
	})

	// report error
	logger.Errorf("test", "id:%d email:%s", 100, "example@example.com")
}

func ExampleNewWithConfig(t *testing.T) {
	ctx := context.Background()
	logger, err := errorreporting.NewWithConfig(ctx, errorreporting.ErrorConfig{
		ProjectID:   "my-project-id",
		ServiceName: "my-service",
		UseSync:     true,
		OnError: func(err error) {
			fmt.Printf("[ERROR] err: %s", err.Error())
		},
	}, config.Config{
		UseTempCredsFile: true,
		CredsJSONBody: `{
			"type":"service_account",
			"private_key":"-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
			"client_email":"example@iam.gserviceaccount.com"
		}`,
		// PrivateKey: "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
		// Email:      "example@iam.gserviceaccount.com",

	})
	if err != nil {
		panic(err)
	}

	logger.SetOnError(func(err error) {
		fmt.Printf("[ReportSync Error] err: %s", err.Error())
	})

	// report error
	logger.Errorf("test", "id:%d email:%s", 100, "example@example.com")
}
