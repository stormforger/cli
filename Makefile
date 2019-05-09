BINARY=forge

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all test build release local_release fmt vet setup lint

all: build

test: vet
	go test -v

build:
	go build -o ${BINARY}

release:
	goreleaser

local_release:
	goreleaser \
	--skip-publish \
	--skip-validate \
	--rm-dist

fmt:
	gofmt -w -s ${GOFILES_NOVENDOR}

vet:
	go vet

lint: vet
	errcheck
	golangci-lint run

setup:
	go get -u github.com/kisielk/errcheck
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/goreleaser/goreleaser
