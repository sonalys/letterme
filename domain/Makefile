#!make
include ../local.env
export

start_mongo:
	docker run --rm -d -p ${MONGO_PORT}:27017 --name letterme_mongo mongo:5.0

start_rabbit:
	docker run --rm -d -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3.9-management

test:
	go test ./... -failfast

coverage:
	go test ./... -race -failfast -covermode=atomic -coverprofile=coverage.out \
	-coverpkg=./cryptography/... ./...
	go tool cover -func coverage.out

lint:
	golangci-lint run --config .golangci.yml --timeout=5m

gen:
	rm -rf ./messaging/mocks && mockery --dir ./messaging/ --output ./messaging/mocks/ --all --case underscore --disable-version-string
	rm -rf ./persistence/mocks && mockery --dir ./persistence/ --output ./persistence/mocks/ --all --case underscore --disable-version-string