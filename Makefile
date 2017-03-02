VERSION="0.0.1"

BINARY=forge

BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=`git rev-parse HEAD | cut -c1-10`

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: fmt
	go build -o ${BINARY}

release: fmt vet
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME} -X main.BUILD_COMMIT=${BUILD_COMMIT}"

fmt:
	gofmt -w -s ${GOFILES_NOVENDOR}

vet:
	script/gorun go vet

test: all
	script/gorun go test -v
