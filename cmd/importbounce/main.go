package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var envConfigURL = os.Getenv("IMPORTBOUNCE_CONFIG_URL")

var (
	flagHTTPAddr  = flag.String("http", "", "Serve HTTP on the provided address instead of AWS Lambda")
	flagConfigURL = flag.String("config", envConfigURL, "Location of the config file to read on each request")
)

func init() {
	http.DefaultClient.Timeout = 2500 * time.Millisecond
}

func main() {
	flag.Parse()

	fetchConfig, err := FetchConfigFuncFromURL(*flagConfigURL)
	if err != nil {
		log.Fatal(err)
	}

	bouncer := &Bouncer{FetchConfig: fetchConfig}

	if *flagHTTPAddr != "" {
		log.Printf("starting HTTP server on %s", *flagHTTPAddr)
		http.ListenAndServe(*flagHTTPAddr, bouncer)
	} else {
		log.Printf("starting AWS Lambda listener")
		lambda.Start(httpadapter.NewV2(bouncer).ProxyWithContext)
	}
}
