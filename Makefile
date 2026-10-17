.PHONY: build
build:
	go build -o bin/gendiff ./cmd/gendiff

.PHONY: help
help:
	go run ./cmd/gendiff -h

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out ./...
