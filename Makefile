BINARY=forge

all: fmt
	go build -o ${BINARY}

fmt:
	gofmt -w .
	goimports -w .
