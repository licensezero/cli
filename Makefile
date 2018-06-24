.PHONY: licensezero test

licensezero:
	go build -o licensezero

test: licensezero
	go test
