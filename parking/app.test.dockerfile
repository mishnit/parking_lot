FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/parking_lot/parking
COPY vendor ../vendor
COPY parking ./
RUN go test ./...
