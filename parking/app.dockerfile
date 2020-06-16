FROM golang:1.13-alpine3.11 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/parking_lot/parking
COPY vendor ../vendor
COPY parking ./
RUN go build -o /go/bin/app ./cmd/parking

FROM alpine:3.11
WORKDIR /usr/bin
ADD parking/static ./static
COPY --from=build /go/bin .
EXPOSE 5566
CMD ["app"]
