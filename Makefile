BINARY=forge

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all test build release local_release fmt vet dep setup

all: build

test: vet
	script/gorun go test -v

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
	script/gorun go vet

errcheck:
	script/gorun errcheck

dep:
	dep ensure && dep prune

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/kisielk/errcheck
	go get -u gopkg.in/alecthomas/gometalinter.v1
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/goreleaser/goreleaser
	gometalinter.v1 --install
