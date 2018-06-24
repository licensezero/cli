licensezero:
	go build -o licensezero

.PHONY: test

test: licensezero
	go test
