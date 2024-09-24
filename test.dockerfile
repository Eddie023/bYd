FROM golang:1.22 as build

WORKDIR /go/src/github.com/eddie023/byd

RUN apt-get update \
    && apt-get install -y -q --no-install-recommends 

ENTRYPOINT [ "/bin/sh", "-c", "go test $(go list --buildvcs=false ./...)" ]