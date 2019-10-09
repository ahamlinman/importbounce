package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	httpAddr := flag.String("http", "", "Serve HTTP on the provided address instead of AWS Lambda")
	flag.Parse()

	bouncer := &bouncer{
		ConfigURL: requireEnv("IMPORTBOUNCE_CONFIG_URL"),
	}

	if *httpAddr != "" {
		log.Printf("starting HTTP server on %s", *httpAddr)
		http.ListenAndServe(*httpAddr, bouncer)
		return
	}

	log.Printf("starting AWS Lambda listener")
	lambda.Start(httpadapter.New(bouncer).Proxy)
}

func requireEnv(name string) string {
	result, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("must have %q in environment", name)
	}
	return result
}

type bouncer struct {
	ConfigURL string
}

func (b *bouncer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
