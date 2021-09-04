#!make
include ../local.env
export

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