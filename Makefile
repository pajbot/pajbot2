swag:
	@swag init --generalInfo cmd/bot/main.go

lint:
	@staticcheck ./...

build:
	@go build ./cmd/bot/

build-csharp:
	@go build -tags csharp ./cmd/bot/

build-all: build build-csharp

check: lint build-all
