#!make
include ./local.env
export

start_mongo:
	docker rm -f letterme_mongo || true
	docker run --rm -d -p ${MONGO_PORT}:27017 --name letterme_mongo mongo:5.0

test:
	go test ./... -race -failfast -covermode=atomic -coverprofile=coverage.out \
	-coverpkg=./controller/...,./persistence...,./utils/... ./...

lint:
	golangci-lint run --config .golangci.yml --timeout=5m

coverage:
	go tool cover -func coverage.out

ci:
	make lint
	make test
	make coverage

gen:
	rm -rf ./mocks/*.go
	mockery --dir ./interfaces/ --all --case underscore --disable-version-string