VERSION="0.0.1"

BINARY=forge

BUILD_TIME=`date +%FT%T%z`
BUILD_COMMIT=`git rev-parse HEAD | cut -c1-10`

all: fmt
	go build -o ${BINARY}

release: fmt
	goxc \
		-pv=${VERSION} \
		-build-ldflags="-X main.VERSION=${VERSION} -X main.BUILD_TIME=${BUILD_TIME} -X main.BUILD_COMMIT=${BUILD_COMMIT}"

fmt:
	gofmt -w .
	goimports -w .

test: all
	go test -v ./...
