FROM ubuntu:18.04

COPY forge /bin

RUN apt-get -y update && apt-get -y install ca-certificates

ENTRYPOINT [ "forge" ]
