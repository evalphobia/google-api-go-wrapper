/*
  example usage:
	$ GOOGLE_API_GO_PRIVATEKEY="`cat secret.pem`" GOOGLE_API_GO_EMAIL=xxx@xxx.gserviceaccount.com go run stackdriver_monitoring.go -project=xxx -metric=foo_sales -value=2
*/
package main

import (
	"flag"
	"fmt"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/stackdriver/monitoring"
)

func main() {
	var projectID, metric string
	var value int64
	flag.StringVar(&projectID, "project", "", "set google project id")
	flag.StringVar(&metric, "metric", "", "metric name")
	flag.Int64Var(&value, "value", 0, "metric value")
	flag.Parse()

	monitor, err := monitoring.NewMonitor(config.Config{}, projectID)
	if err != nil {
		panic(err)
	}

	err = monitor.Create(monitoring.Data{
		Data:       value,
		MetricType: metric,
		MetricKind: monitoring.MetricKindDefault,
		Resource: &monitoring.Resource{
			Type: "global",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("success to create timeseries: %d\n", value)
}
