swag:
	@swag init --generalInfo cmd/bot/main.go

lint:
	@golangci-lint run --new-from-rev=84427cb7eb19ed8edb89b3bfd2962b219691443b

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
