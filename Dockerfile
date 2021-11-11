# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
# Based on https://github.com/golang/go/issues/31997
RUN go env -w GO111MODULE=auto
RUN go mod download

COPY * ./
RUN apt install tree -y
RUN tree /app -L 4
RUN go build ./kawe/cmd/kawe -o /bin/main

EXPOSE 8080

CMD [ "/bin/main" ]
