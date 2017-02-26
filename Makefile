BINARY=forge

all: fmt
	go build -o ${BINARY}

release: fmt
	goxc

fmt:
	gofmt -w .
	goimports -w .

test: all
	go test -v ./...
