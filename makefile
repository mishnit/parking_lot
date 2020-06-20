#!/bin/bash
run: stop up

grpc: #need to update grpc code when there is change in proto definitions
	$(shell go generate parking/grpc_server.go)

mod:
	$(shell GO111MODULE=on go mod tidy)
	$(shell GO111MODULE=on go mod vendor)

stop:
	docker-compose -f docker-compose.yaml stop

up:
	docker-compose -f docker-compose.yaml up -d --build

test:
	docker-compose -f docker-compose-test.yaml up -d --build
	docker system prune -f --volumes

cli:
	@mkdir -p build && \
	cd build && \
	go build -o parking_lot ../parking/commands && \
	chmod -R 777 . && \
	cd ..
