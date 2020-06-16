FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/parking_lot
COPY go.mod go.sum ./
COPY vendor vendor
RUN GO111MODULE=on go test ./...
