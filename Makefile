.PHONY: licensezero test

LDFLAGS=-X main.Rev=$(shell git tag -l --points-at HEAD | sed 's/^v//')

licensezero: prebuild
	go build -o licensezero -ldflags "$(LDFLAGS)"

test: licensezero prebuild
	go test ./...

build: prebuild
	gox -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "$(LDFLAGS)" -verbose

.PHONY: prebuild

prebuild: node_modules
	go get -ldflags "$(LDFLAGS)" ./...
	go generate subcommands/validation.go

node_modules:
	npm install
