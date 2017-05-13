VERSION="0.3.1"

BINARY=forge

BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=`git rev-parse HEAD | cut -c1-10`

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

BUILDINFO_PACKAGE="github.com/stormforger/cli/buildinfo"

LDFLAGS="-X ${BUILDINFO_PACKAGE}.version=${VERSION} -X ${BUILDINFO_PACKAGE}.buildTime=${BUILD_TIME} -X ${BUILDINFO_PACKAGE}.buildCommit=${BUILD_COMMIT}"

.PHONY: all test build build_release local_release fmt vet

all: build

test: vet
	script/gorun go test -v

build:
	go build -o ${BINARY}

build_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags=${LDFLAGS}

local_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags=${LDFLAGS}
		-tasks+="publish-github"

fmt:
	gofmt -w -s ${GOFILES_NOVENDOR}

vet:
	script/gorun go vet
