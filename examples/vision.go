/*
  example usage:
	$ GOOGLE_API_GO_PRIVATEKEY="`cat secret.pem`" GOOGLE_API_GO_EMAIL=xxx@xxx.gserviceaccount.com go run vision.go -file=xxx
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/evalphobia/google-api-go-wrapper/config"
	"github.com/evalphobia/google-api-go-wrapper/vision"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "file", "", "set image file path")
	flag.Parse()

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	cli, err := vision.New(config.Config{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.Safe(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", resp.SafeResult())
}
