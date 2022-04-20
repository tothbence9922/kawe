# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
# Based on https://github.com/golang/go/issues/31997
RUN go env -w GO111MODULE=auto
RUN go mod download

COPY . .
RUN ls -la

RUN go build -o /bin/main ./cmd/kawe/main.go

EXPOSE 8080

CMD [ "/bin/main" ]
