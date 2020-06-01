# FROM golang:1.13.5-alpine3.10 AS build
FROM golang:1.14.3-stretch


ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV TZ=Asia/Shanghai

WORKDIR /root/

VOLUME /root/input

RUN mkdir /root/chier/
COPY ./src /root/chier/
WORKDIR /root/chier/src
RUN go build -o /root/app/main src
WORKDIR /root/app
RUN chmod +x main

ENTRYPOINT ["/root/app/main", "daemon"]