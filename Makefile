BINARY=forge

VERSION=0.1.0
BUILD_TIME=`date +%FT%T%z`

all:
	go build -o ${BINARY}
