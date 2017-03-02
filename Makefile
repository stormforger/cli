VERSION="0.0.2"

BINARY=forge

BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=`git rev-parse HEAD | cut -c1-10`

GOFILES_NOVENDOR = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: fmt build

test: vet
	script/gorun go test -v

build:
	go build -o ${BINARY}

build_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME} -X main.BUILD_COMMIT=${BUILD_COMMIT}"

local_release:
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME} -X main.BUILD_COMMIT=${BUILD_COMMIT}" \
		-tasks+="publish-github"


fmt:
	gofmt -w -s ${GOFILES_NOVENDOR}

vet:
	script/gorun go vet
