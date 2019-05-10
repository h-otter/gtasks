FROM golang:1.11 AS BUILD_GO

ENV GO111MODULE=on

COPY . /go/src/github.com/h-otter/gtasks/api
WORKDIR /go/src/github.com/h-otter/gtasks/api

RUN go build -o gtasks-api .

FROM debian:jessie

# RUN apt update \
#  && apt install -y openssh-client \
#  && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*

COPY --from=BUILD_GO /go/src/github.com/h-otter/gtasks/api/gtasks-api /usr/local/bin/

WORKDIR /
CMD /bin/bash
