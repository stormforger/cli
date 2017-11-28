FROM golang:alpine

COPY . /go/src/github.com/stormforger/cli
WORKDIR /go/src/github.com/stormforger/cli

ENV binary forge

RUN apk --update add ca-certificates make g++ git \
    && go build \
      -v -o ${binary} \
      -ldflags '-s -w -X github.com/stormforger/cli/buildinfo.version={{.Version}} -X github.com/stormforger/cli/buildinfo.commit={{.Commit}} -X github.com/stormforger/cli/buildinfo.date={{.Date}}' \
      . \
    && cp ./${binary} /go/bin \
    && rm -rf .git \
    && apk del make g++ git

ENTRYPOINT [ "forge" ]
