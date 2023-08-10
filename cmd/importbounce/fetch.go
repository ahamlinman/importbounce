package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	xrayawsv2 "github.com/aws/aws-xray-sdk-go/instrumentation/awsv2"
)

// FetchConfigFunc is a type for functions that can load TOML configuration
// files for a Bouncer.
type FetchConfigFunc func(context.Context) (io.ReadCloser, error)

// FetchConfigFuncFromURL returns a FetchConfigFunc based on the value of the
// provided URL string. The following schemes are supported:
//
//	http://{path...}           Retrieve via HTTP request
//	https://{path...}          Retrieve via HTTPS request
//	file://{path...}           Retrieve from the local filesystem
//	s3://{bucket}/{path...}    Retrieve from Amazon S3
func FetchConfigFuncFromURL(urlString string) (FetchConfigFunc, error) {
	if urlString == "" {
		return nil, errors.New("config URL not provided")
	}

	u, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid config URL %q: %w", urlString, err)
	}

	if factory, ok := fetcherFactories[u.Scheme]; ok {
		return factory(u), nil
	}

	return nil, fmt.Errorf("unknown config URL scheme %q", u.Scheme)
}

var fetcherFactories = map[string]func(*url.URL) FetchConfigFunc{
	"http":     getHTTPConfigFetcher,
	"https":    getHTTPConfigFetcher,
	"file":     getFileConfigFetcher,
	"s3":       getS3ConfigFetcher,
	"s3+nossl": getS3ConfigFetcher,
}

func getHTTPConfigFetcher(u *url.URL) FetchConfigFunc {
	return func(ctx context.Context) (io.ReadCloser, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("fetching config: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("fetching config: %w", err)
		}

		return resp.Body, nil
	}
}

func getFileConfigFetcher(u *url.URL) FetchConfigFunc {
	return func(_ context.Context) (io.ReadCloser, error) {
		path := filepath.Join(u.Host, u.Path)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening config: %w", err)
		}
		return f, nil
	}
}

func getS3ConfigFetcher(u *url.URL) FetchConfigFunc {
	input := &s3.GetObjectInput{
		Bucket: aws.String(u.Host),
		Key:    aws.String(strings.TrimPrefix(u.Path, "/")),
	}

	cfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(fmt.Errorf("loading AWS config: %w", err))
	}
	xrayawsv2.AWSV2Instrumentor(&cfg.APIOptions)

	disableSSL := strings.HasSuffix(u.Scheme, "+nossl")
	s3Client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.EndpointOptions.DisableHTTPS = disableSSL
	})

	return func(ctx context.Context) (io.ReadCloser, error) {
		output, err := s3Client.GetObject(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("fetching config: %w", err)
		}
		return output.Body, nil
	}
}
