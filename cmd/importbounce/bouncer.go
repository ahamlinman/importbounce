package main

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

// Bouncer handles HTTP requests for Go imports by retrieving a configuration
// file, finding a matching package prefix, and serving an appropriate
// redirect.
type Bouncer struct {
	FetchConfig FetchConfigFunc
}

var responseTmpl = template.Must(template.New("").Parse(`<html>
<head>
<meta name="go-import" content="{{.Prefix}} {{.Import}}">
<meta http-equiv="refresh" content="0; url={{.Redirect}}">
</head>
<body>Redirecting…</body>
</html>`))

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

func (b *Bouncer) loadConfig(ctx context.Context) (config, error) {
	r, err := b.FetchConfig(ctx)
	if err != nil {
		return config{}, xerrors.Errorf("fetching config: %w", err)
	}
	defer r.Close()

	var c config
	_, err = toml.DecodeReader(r, &c)
	if err != nil {
		err = xerrors.Errorf("decoding config: %w", err)
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
