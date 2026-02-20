GO_BINARY=/usr/bin/go
EXAMPLES_DIR := examples
FILES := $(wildcard $(EXAMPLES_DIR)/*)
build: 
	${GO_BINARY} build .

test:
	${GO_BINARY} test -json ./... | tparse -all

run-tremor:
	@for file in $(FILES); do \
		echo "Running tremor on $$file ..."; \
		./tremor "$$file" || { echo "Error: tremor failed on $$file"; exit 1; }; \
		done
	@echo "All files processed successfully."


all: build test


.PHONY: test all build
.DEFAULT_GOAL := all
