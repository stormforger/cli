VERSION="0.0.2"

BINARY=forge

BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=`git rev-parse HEAD | cut -c1-10`

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all test build build_release local_release fmt vet

all: build

test: vet
	script/gorun go test -v

build:
	go build -o ${BINARY}

build_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X buildinfo.Version=${VERSION} -X buildinfo.BuildTime=${BUILD_TIME} -X buildinfo.BuildCommit=${BUILD_COMMIT}"

local_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X buildinfo.Version=${VERSION} -X buildinfo.BuildTime=${BUILD_TIME} -X buildinfo.BuildCommit=${BUILD_COMMIT}" \
		-tasks+="publish-github"

fmt:
	gofmt -w -s ${GOFILES_NOVENDOR}

vet:
	script/gorun go vet
