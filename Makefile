.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/wormhole

.PHONY: test
test:
	go test ./... -race

.PHONY: test-cov
test-cov:
	go test ./... -race -coverprofile=coverage.out
	go tool cover -html=coverage.out
