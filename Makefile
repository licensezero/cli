.PHONY: licensezero test

licensezero:
	go build -o licensezero

test: licensezero
	go test

build:
	gox -os="linux darwin windows freebsd" -arch="i386 amd64" -output="licensezero-{{.OS}}-{{.Arch}}" -ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose
