FROM golangci/golangci-lint:v1.59.1

RUN apt-get update \
    && apt-get install -y -q --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

ENV GOFLAGS=-buildvcs=false 

WORKDIR /src

ENTRYPOINT [ "golangci-lint", "run" ]