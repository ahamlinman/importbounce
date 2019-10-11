lambda-dist/importbounce: go.mod go.sum $(shell find . -name '*.go')
	mkdir -p lambda-dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
			 go build -v -mod=vendor -ldflags="-s -w" -o "$@" ./cmd/importbounce
