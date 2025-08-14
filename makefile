build: 
	go build .

test:
	go test -json ./typechecker | tparse -all

all: build test

.PHONY: test all build
.DEFAULT_GOAL := all
