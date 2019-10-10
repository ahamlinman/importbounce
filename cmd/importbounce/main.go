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

func main() {
	http.DefaultClient.Timeout = 2500 * time.Millisecond

	var (
		httpAddr = flag.String(
			"http", "", "Serve HTTP on the provided address instead of AWS Lambda",
		)

		defaultConfigURL = os.Getenv("IMPORTBOUNCE_CONFIG_URL")
		configURL        = flag.String(
			"config", defaultConfigURL, "Location of the config file to read on each request",
		)
	)

	flag.Parse()

	fetchConfig, err := FetchConfigFuncFromURL(*configURL)
	if err != nil {
		log.Fatal(err)
	}

	bouncer := &Bouncer{FetchConfig: fetchConfig}

	if *httpAddr != "" {
		log.Printf("starting HTTP server on %s", *httpAddr)
		http.ListenAndServe(*httpAddr, bouncer)
		return
	}

	log.Printf("starting AWS Lambda listener")
	lambda.Start(httpadapter.New(bouncer).ProxyWithContext)
}
