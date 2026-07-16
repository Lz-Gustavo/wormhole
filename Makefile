.PHONY: all
all: build

.PHONY: build
build:
	go build -o bin/wormhole

.PHONY: test
test:
	go test ./...
