FROM golang:1.22 as build

WORKDIR /go/src/github.com/eddie023/byd

RUN apt-get update \
    && apt-get install -y -q --no-install-recommends 

# Download all go dependencies in the image to make tests runs faster.
COPY go.mod .
RUN go mod download -x

ENTRYPOINT [ "/bin/sh", "-c", "go test $(go list --buildvcs=false ./...)" ]
