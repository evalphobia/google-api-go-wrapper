google-api-go-wrapper
====

golang wrapper library of [Google APIs Client Library for Go](https://github.com/google/google-api-go-client)

# Supported API

- [Google Analytics](https://godoc.org/google.golang.org/api/analytics/v3)
    - Realtime.Get
- [BigQuery](https://godoc.org/google.golang.org/api/bigquery/v2)
    - create table
    - InsertAll
- [Stackdriver logging](https://godoc.org/google.golang.org/api/logging/v2)
    - Write
- [Stackdriver monitoring](https://godoc.org/google.golang.org/api/monitoring/v3)
    - TimeSeries.Create
- [Cloud Vision](https://godoc.org/google.golang.org/api/vision/v1)
    - Annotate

# Requirements

Depends on the google's each libraries.


# Installation

Install google-api-go-wrapper and required packages using `go get` command:

```bash
$ go get github.com/evalphobia/google-api-go-wrapper/...
```

# Usage

## Config usage

```go
// create client by given paramter.
client = analytics.New(config.Config{
    PrivateKey: `-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----`,
    Email: "xxx@xxx.gserviceaccount.com",
})

// create client by given file path of credential.
client = analytics.New(config.Config{
    Filename: "/path/to/pem.json",
})

// use env parameter when above fields are empty.
// $GOOGLE_API_GO_PRIVATEKEY : used as `PrivateKey` field.
// $GOOGLE_API_GO_EMAIL : used as `Email` field.
client = analytics.New(config.Config{})
```

If no other credentials could be found, `Config` will use https://godoc.org/golang.org/x/oauth2/google#FindDefaultCredentials.

## Google Analytics

Install google's library:

```bash
$ go get google.golang.org/api/analytics/v3
```

### Reatime ActiveUser

```go
import (
    "github.com/evalphobia/google-api-go-wrapper/analytics"
    "github.com/evalphobia/google-api-go-wrapper/config"
)

...


client := analytics.New(config.Config{})

viewID := "00000000"
result, err := cli.GetRealtimeActiveUser(viewID)
if err != nil {
    fmt.Printf("[ERROR] %s\n", err.Error())
    return
}

fmt.Printf("activeUser=%d\n", result)
```


## BigQuery

Install google's library:

```bash
$ go get google.golang.org/api/bigquery/v2
```

### CreateTable

```go
import (
    "github.com/evalphobia/google-api-go-wrapper/bigquery"
    "github.com/evalphobia/google-api-go-wrapper/config"
)

...


ds, err := bigquery.NewDataset(config.Config{}, projectID, datasetID)
if err != nil {
    panic(err)
}

err = ds.CreateTable(tableID, MySchema{})
if err != nil {
    panic(err)
}

...

type MySchema struct {
    Name      string    `bigquery:"username"`
    CreatedAt time.Time `bigquery:"created_at"`
}
```

### InsertAll

```go
import (
    "github.com/evalphobia/google-api-go-wrapper/bigquery"
    "github.com/evalphobia/google-api-go-wrapper/config"
)

...


ds, err := bigquery.NewDataset(config.Config{}, projectID, datasetID)
if err != nil {
    panic(err)
}

err = ds.InserAll(tableID, &MySchema{
    Name:      "foo",
    CreatedAt: time.Now(),
})
if err != nil {
    panic(err)
}

...

type MySchema struct {
    Name      string    `bigquery:"username"`
    CreatedAt time.Time `bigquery:"created_at"`
}
```


## Stackdriver

Install google's library:

```bash
$ go get google.golang.org/api/logging/v2
$ go get google.golang.org/api/monitoring/v3
```

### logging

```go
import (
    "github.com/evalphobia/google-api-go-wrapper/config"
    "github.com/evalphobia/google-api-go-wrapper/stackdriver/logging"
)

...


logger, err := logging.NewLogger(config.Config{}, projectID)
if err != nil {
    panic(err)
}

err = logger.Write(logging.WriteData{
    Data:    map[string]interface{}{"key": "value"},
    LogName: "test_log",
    Resource: &logging.Resource{
        Type: "global",
    },
})
if err != nil {
    panic(err)
}
```

### monitoring

```go
import (
    "github.com/evalphobia/google-api-go-wrapper/config"
    "github.com/evalphobia/google-api-go-wrapper/stackdriver/monitoring"
)

...


monitor, err := monitoring.NewMonitor(config.Config{}, projectID)
if err != nil {
    panic(err)
}

err = monitor.Create(monitoring.Data{
    Data:       100,
    MetricType: "foo_sales",
    MetricKind: monitoring.MetricKindGauge,
    Resource: &monitoring.Resource{
        Type: "global",
    },
})
if err != nil {
    panic(err)
}
```

## Cloud Vision

Install google's library:

```bash
$ go get google.golang.org/api/vision/v1
```

### Annotate

```go
import (
    "fmt"
    "io/ioutil"

    "github.com/evalphobia/google-api-go-wrapper/config"
    "github.com/evalphobia/google-api-go-wrapper/vision"
)

...


client, err := vision.New(config.Config{})
if err != nil {
    panic(err)
}

img, err := ioutil.ReadFile(file)
if err != nil {
    panic(err)
}

faceResult, err := cli.Face(img)
if err != nil {
    panic(err)
}
fmt.Printf("FaceDetect=%+v\n", faceResult.FaceResult())

safeResult, err := cli.Safe(img)
if err != nil {
    panic(err)
}
fmt.Printf("SafeSearch=%+v\n", safeResult.SafeResult())

textResult, err := cli.Text(img)
if err != nil {
    panic(err)
}
fmt.Printf("TextDetect=%+v\n", textResult.TextResult())
```

