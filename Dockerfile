# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY . .
# Based on https://github.com/golang/go/issues/31997
RUN go env -w GO111MODULE=auto
RUN go mod download
RUN go build -o /bin/main ./main.go

FROM alpine:latest
RUN adduser -u 2222 -D kawe

WORKDIR /

COPY --from=builder /bin/main /bin/main

USER kawe

EXPOSE 80

CMD [ "/bin/main" ]
