swag:
	@swag init --generalInfo cmd/bot/main.go

lint:
	@golangci-lint run --new-from-rev=2baeaf2880~13

lint-all:
	@golangci-lint run

lint-current:
	@golangci-lint run --new

build:
	@go build ./cmd/bot/

build-csharp:
	@go build -tags csharp ./cmd/bot/

build-all: build build-csharp

check: lint build-all
