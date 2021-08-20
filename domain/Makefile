#!make

test:
	go test ./... -race -failfast -covermode=atomic -coverprofile=coverage.out \
	-coverpkg=./cryptography/... ./...

lint:
	golangci-lint run --config .golangci.yml --timeout=5m

ci:
	make test
	make lint