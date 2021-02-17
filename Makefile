DEFAULT_TARGET: test

CURRENT_DIR=$(shell pwd)
SRCS := $(shell find . -type f -name '*.go')

.PHONY: test
test:
	go test ./...

ci:
	pre-commit run -a
	$(MAKE) test
