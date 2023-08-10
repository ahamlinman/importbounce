package bouncer

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

// fetcherFunc is a type for functions that can load TOML configuration files
// for a Bouncer.
type fetcherFunc func(context.Context) (io.ReadCloser, error)

func getFetcherFromURL(configURL string) (fetcherFunc, error) {
	if configURL == "" {
		return nil, errors.New("config URL not provided")
	}

	u, err := url.Parse(configURL)
	if err != nil {
		return nil, fmt.Errorf("invalid config URL %q: %w", configURL, err)
	}

	factory, ok := fetcherFactories[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("unknown config URL scheme %q", u.Scheme)
	}
	return factory(u), nil
}

var fetcherFactories = map[string]func(*url.URL) fetcherFunc{
	"http":     getHTTPConfigFetcher,
	"https":    getHTTPConfigFetcher,
	"file":     getFileConfigFetcher,
	"s3":       getS3ConfigFetcher,
	"s3+nossl": getS3ConfigFetcher,
}

func getHTTPConfigFetcher(u *url.URL) fetcherFunc {
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

func getFileConfigFetcher(u *url.URL) fetcherFunc {
	return func(_ context.Context) (io.ReadCloser, error) {
		path := filepath.Join(u.Host, u.Path)
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("opening config: %w", err)
		}
		return f, nil
	}
}

func getS3ConfigFetcher(u *url.URL) fetcherFunc {
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
