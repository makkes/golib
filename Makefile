DEFAULT_TARGET: test

CURRENT_DIR=$(shell pwd)
SRCS := $(shell find . -type f -name '*.go')

.PHONY: lint
lint:
	pre-commit run --all-files

.PHONY: test
test: lint
	go test ./...
