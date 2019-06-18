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
	go generate subcommands/validation.go
	go get -ldflags "$(LDFLAGS)" ./...

node_modules:
	npm install
