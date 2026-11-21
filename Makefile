swag:
	@swag init --generalInfo cmd/bot/main.go

lint:
	@go tool staticcheck ./...

build:
	@go build ./cmd/bot/

build-csharp:
	@dotnet publish 3rdParty/MessageHeightTwitch/MessageHeightTwitch.csproj -c Release -r linux-x64 -o cmd/bot/MessageHeightTwitch
	@go build -tags csharp ./cmd/bot/

build-all: build build-csharp

check: lint build-all
