#!/bin/bash
run: stop up

grpc: #need to update protobuf code when there is change in grpc code
	$(shell go generate parking/grpc_server.go)

mod:
	$(shell GO111MODULE=on go mod tidy)
	$(shell GO111MODULE=on go mod vendor)

stop:
	docker-compose -f docker-compose.yaml stop

up:
	docker-compose -f docker-compose.yaml up -d --build --remove-orphans

test:
	docker-compose -f docker-compose-test.yaml up -d --build --remove-orphans

cli:
	@mkdir -p build && \
	cd build && \
	go build -o parking_lot ../commando/cmd/commando && \
	chmod -R 777 . && \
	cd ..
