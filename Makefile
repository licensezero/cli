.PHONY: licensezero test

LDFLAGS=-X main.Rev=$(shell git tag -l --points-at HEAD | sed 's/^v//')

licensezero: node_modules
	go get -ldflags "$(LDFLAGS)" ./...
	go generate subcommands/validation.go
	go build -o licensezero -ldflags "$(LDFLAGS)"

test: licensezero
	go test ./...

build:
	go get -ldflags "$(LDFLAGS)" ./...
	go generate subcommands/validation.go
	gox -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS)" -verbose

node_modules:
	npm install
