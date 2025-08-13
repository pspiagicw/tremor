build: 
	go build .

test:
	go test -json ./... | tparse -all

all: build test

.PHONY: test all build
.DEFAULT_GOAL := all
