run: stop up

mod:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

stop:
	docker-compose -f docker-compose.yaml stop

up:
	docker-compose -f docker-compose.yaml up -d --build

down:
	docker-compose -f docker-compose.yaml down

test:
	docker-compose -f docker-compose-test.yaml up --build --abort-on-container-exit
	docker-compose -f docker-compose-test.yaml down --volumes

test-db-up:
	docker-compose -f docker-compose-test.yaml up --build db

test-db-down:
	docker-compose -f docker-compose-test.yaml down --volumes db
