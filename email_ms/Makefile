#!make
include ../local.env
export

start_mongo:
	docker run --rm -d -p ${MONGO_PORT}:27017 --name letterme_mongo mongo:5.0

test:
	go test ./... -failfast

coverage:
	go test ./... -failfast -covermode=atomic -coverprofile=coverage.out \
	-coverpkg=./... ./...
	go tool cover -func coverage.out

lint:
	golangci-lint run --config .golangci.yml --timeout=5m

gen:
	rm -rf ./mocks/*.go
	mockery --dir ./smtp/ --all --case underscore --disable-version-string

vendor:
	go mod vendor