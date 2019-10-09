package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"golang.org/x/xerrors"
)

func main() {
	http.DefaultClient.Timeout = 10 * time.Second

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
	config, err := b.loadConfig()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: not this
	json.NewEncoder(w).Encode(config)
}

type config struct {
	Packages []packageConfig `toml:"packages"`
}

type packageConfig struct {
	Prefix   string `toml:"prefix"`
	Import   string `toml:"import"`
	Redirect string `toml:"redirect"`
}

func (b *bouncer) loadConfig() (config, error) {
	resp, err := http.Get(b.ConfigURL)
	if err != nil {
		return config{}, xerrors.Errorf("fetching config: %w", err)
	}
	defer resp.Body.Close()

	var c config
	_, err = toml.DecodeReader(resp.Body, &c)
	if err != nil {
		err = xerrors.Errorf("decoding config: %w", err)
	}
	return c, err
}
