.PHONY: build
build:
	go build -o bin/CA-service ./cmd/CA-service

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint run

.PHONY: format
format:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint fmt
