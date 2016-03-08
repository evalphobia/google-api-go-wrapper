/*
  example usage:
	$ go run bigquery.go -project=xxx -dataset=yyy -table=zzz -email=xxx@xxx.gserviceaccount.com -secret="`cat secret.pem`"
*/
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/evalphobia/google-api-go-wrapper/bigquery"
)

var email, secret string
var projectID, datasetID, tableID string

func init() {
	flag.StringVar(&secret, "secret", "", "set service account private key")
	flag.StringVar(&email, "email", "xxxxx@xxxxx.gserviceaccount.com", "set service account email")
	flag.StringVar(&projectID, "project", "", "set google project id")
	flag.StringVar(&datasetID, "dataset", "", "set bigquery dataset id")
	flag.StringVar(&tableID, "table", "", "set bigquery table id")
}

func main() {
	flag.Parse()

	cli := bigquery.NewWithParams(secret, email)
	ds, err := bigquery.NewDataset(cli, projectID, datasetID)
	if err != nil {
		panic(err)
		return
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
