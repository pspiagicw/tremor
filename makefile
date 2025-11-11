GO_BINARY=/usr/bin/go
build: 
	${GO_BINARY} build .

test:
	${GO_BINARY} test -json ./... | tparse -all

all: build test

.PHONY: test all build
.DEFAULT_GOAL := all
