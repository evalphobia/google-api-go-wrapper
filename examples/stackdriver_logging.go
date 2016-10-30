/*
  example usage:
	$ GOOGLE_API_GO_PRIVATEKEY="`cat secret.pem`" GOOGLE_API_GO_EMAIL=xxx@xxx.gserviceaccount.com go run stackdriver_logging.go  -project=xxx -data='{"key": "value"}'
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/stackdriver/logging"
)

func main() {
	var projectID, data string
	flag.StringVar(&projectID, "project", "", "set google project id")
	flag.StringVar(&data, "data", "", "logging data(JSON)")
	flag.Parse()

	entryBody := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &entryBody)
	if err != nil {
		panic(err)
	}

	logger, err := logging.NewLogger(config.Config{}, projectID)
	if err != nil {
		panic(err)
	}

	err = logger.Write(logging.WriteData{
		Data:    entryBody,
		LogName: "test_log",
		Resource: &logging.Resource{
			Type: "global",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("success to logging: %+v\n", entryBody)
}
