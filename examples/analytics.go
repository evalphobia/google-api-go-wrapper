/*
  example usage:
	$ GOOGLE_API_GO_PRIVATEKEY="`cat secret.pem`" GOOGLE_API_GO_EMAIL=xxx@xxx.gserviceaccount.com go run analytics.go -viewid=xxx
*/
package main

import (
	"flag"
	"fmt"

	"github.com/evalphobia/google-api-go-wrapper/analytics"
	"github.com/evalphobia/google-api-go-wrapper/config"
)

func main() {
	var viewID string
	flag.StringVar(&viewID, "viewid", "00000000", "set google analytics view id")
	flag.Parse()

	cli, err := analytics.New(config.Config{})
	if err != nil {
		panic(err)
	}

	count, err := cli.GetRealtimeActiveUser(viewID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("view=%s, activeUser=%d\n", viewID, count)
}
