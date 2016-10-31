/*
  example usage:
	$ GOOGLE_API_GO_PRIVATEKEY="`cat secret.pem`" GOOGLE_API_GO_EMAIL=xxx@xxx.gserviceaccount.com go run bigquery.go  -project=xxx -dataset=yyy -table=zzz
*/
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/evalphobia/google-api-go-wrapper/bigquery"
	"github.com/evalphobia/google-api-go-wrapper/config"
)

func main() {
	var projectID, datasetID, tableID string
	flag.StringVar(&projectID, "project", "", "set google project id")
	flag.StringVar(&datasetID, "dataset", "", "set bigquery dataset id")
	flag.StringVar(&tableID, "table", "", "set bigquery table id")
	flag.Parse()

	ds, err := bigquery.NewDataset(config.Config{}, projectID, datasetID)
	if err != nil {
		panic(err)
	}

	err = ds.CreateTable(tableID, Schema{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("success to create table: %s.%s\n", datasetID, tableID)
}

type Schema struct {
	Name      string    `bigquery:"username"`
	CreatedAt time.Time `bigquery:"created_at"`
}
