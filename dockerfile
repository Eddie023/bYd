FROM golang:1.22 as build
ARG VERISON 
WORKDIR /go/src/apiserver

COPY go.mod go.sum ./
RUN go mod download 

ADD . .
RUN go install ./cmd/service

FROM debian:bookworm
ARG VERSION 
RUN apt-get update && apt-get install -y -q --no-install-recommends 

COPY --from=build /go/bin/service /bin/apiserver
ENTRYPOINT [ "apiserver" ]

LABEL image.authors="Manish Chaulagain" image.version=${VERSION}