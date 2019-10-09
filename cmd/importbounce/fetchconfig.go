package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type configFetcher func() (io.ReadCloser, error)

func getConfigFetcher() configFetcher {
	urlString, ok := os.LookupEnv("IMPORTBOUNCE_CONFIG_URL")
	if !ok {
		log.Fatal("must set IMPORTBOUNCE_CONFIG_URL")
	}

	u, err := url.Parse(urlString)
	if err != nil {
		log.Fatalf("invalid IMPORTBOUNCE_CONFIG_URL: %v", err)
	}

	if factory, ok := fetcherFactories[u.Scheme]; ok {
		return factory(u)
	}

	log.Fatalf("unrecognized IMPORTBOUNCE_CONFIG_URL scheme %v", u.Scheme)
	return nil
}

var fetcherFactories = map[string]func(*url.URL) configFetcher{
	"http": getHTTPConfigFetcher,
	"file": getFileConfigFetcher,
}

func getHTTPConfigFetcher(u *url.URL) configFetcher {
	return func() (io.ReadCloser, error) {
		resp, err := http.Get(u.String())
		if err != nil {
			err = xerrors.Errorf("fetching config: %w", err)
		}
		return resp.Body, err
	}
}

func getFileConfigFetcher(u *url.URL) configFetcher {
	return func() (io.ReadCloser, error) {
		path := filepath.Join(u.Host, u.Path)
		f, err := os.Open(path)
		if err != nil {
			err = xerrors.Errorf("opening config: %w", err)
		}
		return f, err
	}
}
