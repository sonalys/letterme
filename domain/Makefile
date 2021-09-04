#!make
include ../local.env
export

start_mongo:
	docker run --rm -d -p ${MONGO_PORT}:27017 --name letterme_mongo mongo:5.0

start_rabbit:
	docker run --rm -d -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3.9-management

test:
	go test ./... -race -failfast -covermode=atomic -coverprofile=coverage.out \
	-coverpkg=./cryptography/... ./...

lint:
	golangci-lint run --config .golangci.yml --timeout=5m

ci:
	make test
	make lint

gen:
	rm -rf ./mocks/*.go
	mockery --dir ./interfaces/ --all --case underscore --disable-version-string