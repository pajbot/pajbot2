swag:
	@swag init --generalInfo cmd/bot/main.go

lint:
	@golangci-lint run --new-from-rev=2baeaf2880~3
