build:
	go build -o bin/gendiff ./cmd/gendiff

help:
	go run ./cmd/gendiff -h

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./...