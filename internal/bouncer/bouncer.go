package bouncer

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
)

// Bouncer handles HTTP requests for Go imports by retrieving a configuration
// file, finding a matching package prefix, and serving an appropriate
// redirect.
type Bouncer struct {
	FetchConfig FetchConfigFunc
}

// New creates a new Bouncer using the configuration from the provided URL.
// The following URL schemes are supported:
//
//	https://{path...}               Retrieve via HTTPS request
//	http://{path...}                Retrieve via HTTP request
//	file://{path...}                Retrieve from the local filesystem
//	s3://{bucket}/{path...}         Retrieve from Amazon S3 with HTTPS
//	s3+nossl://{bucket}/{path...}   Retrieve from Amazon S3 with HTTP
func New(configURL string) (*Bouncer, error) {
	fetchConfig, err := FetchConfigFuncFromURL(configURL)
	if err != nil {
		return nil, err
	}
	return &Bouncer{FetchConfig: fetchConfig}, nil
}

func (b *Bouncer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config, err := b.loadConfig(r.Context())
	if err != nil {
		log.Printf("failed to load config: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	path := r.Host + r.URL.Path
	pkgConf := config.FindPackage(path)

	if pkgConf == (packageConfig{}) {
		b.tryDefaultRedirect(w, r, config.DefaultRedirect)
		return
	}

	if r.URL.Query().Get("go-get") == "" {
		http.Redirect(w, r, pkgConf.Redirect, http.StatusFound)
		return
	}

	if err := responseTmpl.Execute(w, pkgConf); err != nil {
		// This is going to be best-effort.
		w.WriteHeader(http.StatusInternalServerError)
	}
}

var responseTmpl = template.Must(template.New("").Parse(`<html>
<head>
<meta name="go-import" content="{{.Prefix}} {{.Import}}">
<meta http-equiv="refresh" content="0; url={{.Redirect}}">
</head>
<body>Redirecting…</body>
</html>`))

func (b *Bouncer) loadConfig(ctx context.Context) (config, error) {
	r, err := b.FetchConfig(ctx)
	if err != nil {
		return config{}, fmt.Errorf("fetching config: %w", err)
	}
	defer r.Close()

	var c config
	_, err = toml.NewDecoder(r).Decode(&c)
	if err != nil {
		return config{}, fmt.Errorf("decoding config: %w", err)
	}
	return c, err
}

func (b *Bouncer) tryDefaultRedirect(w http.ResponseWriter, r *http.Request, url string) {
	if url == "" || r.URL.Query().Get("go-get") != "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Package not found\n"))
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
