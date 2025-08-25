build: 
	go build .

test:
	go test -json ./parser | tparse -all

all: build test

.PHONY: test all build
.DEFAULT_GOAL := all
