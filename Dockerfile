FROM ubuntu:trusty

COPY forge /bin

RUN apt-get -y update && apt-get -y install ca-certificates

ENTRYPOINT [ "forge" ]
