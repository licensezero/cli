.PHONY: licensezero test

LDFLAGS=-X main.Rev=$(shell git tag -l --points-at HEAD | sed 's/^v//')

licensezero:
	go build -o licensezero -ldflags "$(LDFLAGS)"

test: licensezero
	go test

build:
	gox -os="linux darwin windows freebsd" -arch="386 amd64 arm" -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS)" -verbose
