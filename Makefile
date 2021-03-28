# We make a small effort to reduce the size of the binary to help Lambda
# download it more quickly, but try to produce a binary no matter what.
#
# (In my brief testing "strip" removes slightly more from the binary than Go's
# -ldflags="-s -w" alone.)
lambda-dist/importbounce: go.mod go.sum $(shell find . -name '*.go')
	mkdir -p lambda-dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
			 go build -v \
				 -mod=vendor \
				 -trimpath -ldflags="-s -w" \
				 -o "$@" \
				 ./cmd/importbounce
	if type strip; then strip "$@"; fi
