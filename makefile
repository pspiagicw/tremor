build: 
	go build .

test:
	go test ./...

all: build test

.PHONY: test all build
.DEFAULT_GOAL := all
