lint:
	@golangci-lint run

build:
	@cd cmd/bot && go build

check: lint
