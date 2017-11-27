FROM golang:alpine

COPY forge /bin

RUN apk --update add ca-certificates

ENTRYPOINT ["/bin/forge"]
