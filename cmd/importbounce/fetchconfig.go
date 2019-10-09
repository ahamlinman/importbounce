package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	"s3":   getS3ConfigFetcher,
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

func getS3ConfigFetcher(u *url.URL) configFetcher {
	sess := session.Must(session.NewSession())
	s3Client := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(strings.TrimPrefix(u.Path, "/")),
	}

	return func() (io.ReadCloser, error) {
		output, err := s3Client.GetObject(input)
		if err != nil {
			err = xerrors.Errorf("fetching config: %w", err)
		}
		return output.Body, err
	}
}
