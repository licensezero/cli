.PHONY: licensezero test

LDFLAGS=-X main.Rev=$(shell git tag -l --points-at HEAD | sed 's/^v//')

licensezero:
	go build -o licensezero -ldflags "$(LDFLAGS)"

test: licensezero
	go test

build:
	gox -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS)" -verbose
