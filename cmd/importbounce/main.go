package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"strings"
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
		FetchConfig: getConfigFetcher(),
	}

	if *httpAddr != "" {
		log.Printf("starting HTTP server on %s", *httpAddr)
		http.ListenAndServe(*httpAddr, bouncer)
		return
	}

	log.Printf("starting AWS Lambda listener")
	lambda.Start(httpadapter.New(bouncer).Proxy)
}

type bouncer struct {
	FetchConfig configFetcher
}

var responseTmpl = template.Must(template.New("").Parse(`<html>
<head>
<meta name="go-import" content="{{.Prefix}} {{.Import}}">
<meta http-equiv="refresh" content="0; url={{.Redirect}}">
</head>
<body>
Redirecting to <a href="{{.Redirect}}">{{.Redirect}}</a>...
</body>
</html>`))

func (b *bouncer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config, err := b.loadConfig()
	if err != nil {
		log.Printf("failed to load config: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// The hostname provided by the Lambda proxy library is a hardcoded default,
	// rather than the one requested by the user. This can be overridden with an
	// environment variable, but reading the Host header manually should be more
	// robust.
	path := r.Header.Get("Host") + r.URL.Path
	pkgConf := config.FindPackage(path)

	if pkgConf == (packageConfig{}) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Package not found\n"))
		return
	}

	if r.URL.Query().Get("go-get") == "" {
		http.Redirect(w, r, pkgConf.Redirect, http.StatusFound)
		return
	}

	err = responseTmpl.Execute(w, pkgConf)
	if err != nil {
		// This is going to be best-effort
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (b *bouncer) loadConfig() (*config, error) {
	r, err := b.FetchConfig()
	if err != nil {
		return nil, xerrors.Errorf("fetching config: %w", err)
	}
	defer r.Close()

	var c config
	_, err = toml.DecodeReader(r, &c)
	if err != nil {
		err = xerrors.Errorf("decoding config: %w", err)
	}
	return &c, err
}

type config struct {
	Packages []packageConfig `toml:"packages"`
}

type packageConfig struct {
	Prefix   string `toml:"prefix"`
	Import   string `toml:"import"`
	Redirect string `toml:"redirect"`
}

func (c *config) FindPackage(path string) packageConfig {
	for _, pkgConf := range c.Packages {
		prefix := strings.TrimSuffix(pkgConf.Prefix, "/")

		if !strings.HasPrefix(path, prefix) {
			continue
		}

		rest := path[len(prefix):]
		if len(rest) != 0 && !strings.HasPrefix(rest, "/") {
			continue
		}

		return pkgConf
	}

	return packageConfig{}
}
