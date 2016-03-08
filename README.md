google-api-go-wrapper
====

golang wrapper library of [Google APIs Client Library for Go](https://github.com/google/google-api-go-client)

# Supported API

- [analytics](https://godoc.org/google.golang.org/api/analytics/v3)
    - Realtime.Get
- [bigquery](https://godoc.org/google.golang.org/api/bigquery/v2)
    - create table
    - InsertAll

# Requirements

Depends on the google's each libraries.


# Installation

Install google-api-go-wrapper and required packages using `go get` command:

```bash
$ go get github.com/evalphobia/google-api-go-wrapper/...
```

# Usage

## analytics

Install google's library:

```bash
$ go get google.golang.org/api/analytics/v3
```

### Reatime ActiveUser

```go
import "github.com/evalphobia/google-api-go-wrapper/analytics"


client := analytics.NewWithParams(secret, email)

viewID := "00000000"
result, err := cli.GetRealtimeActiveUser(viewID)
if err != nil {
    fmt.Printf("[ERROR] %s\n", err.Error())
    return
}

fmt.Printf("activeUser=%d\n", result)
```
