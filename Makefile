NAME := zip4win
VERSION := $(shell git describe --tags --abbrev=0)
REVISION := $(shell git rev-parse --short HEAD)
LDFLAGS := -X 'main.version=$(VERSION)' \
	   -X 'main.revision=$(REVISION)'

# Default target
default: bin/$(NAME)

## Setup
setup:
	go get github.com/pkg/errors
	go get golang.org/x/text/encoding/japanese
	go get golang.org/x/text/transform
	go get golang.org/x/text/unicode/norm

## Run tests
test:
	go test

## build binaries ex. make bin/zip4win
bin/%: cmd/%/main.go
	go build -ldflags "$(LDFLAGS)" -o $@ $<

.PHONY: default setup test
