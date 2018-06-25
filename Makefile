.PHONY: licensezero test

LDFLAGS=-X main.Rev=$(shell git rev-parse --short HEAD)

licensezero:
	go build -o licensezero -ldflags "$(LDFLAGS)"

test: licensezero
	go test

build:
	gox -os="linux darwin windows freebsd" -arch="386 amd64" -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS)" -verbose
