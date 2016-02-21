/*
  example usage:
	$ go run analytics.go -viewid=xxx -email=xxx@xxx.gserviceaccount.com -secret="`cat secret.pem`"
*/
package main

import (
	"flag"
	"fmt"

	"github.com/evalphobia/google-api-go-wrapper/analytics"
)

var email, secret, viewID string

func init() {
	flag.StringVar(&secret, "secret", "", "set service account private key")
	flag.StringVar(&email, "email", "xxxxx@xxxxx.gserviceaccount.com", "set service account email")
	flag.StringVar(&viewID, "viewid", "00000000", "set google analytics view id")
}

func main() {
	flag.Parse()

	cli := analytics.NewWithParams(secret, email)
	count, err := cli.GetRealtimeActiveUser(viewID)
	if err != nil {
		panic(err)
		return
	}

	fmt.Printf("view=%s, activeUser=%d\n", viewID, count)
}
