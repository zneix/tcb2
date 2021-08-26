lint:
	@staticcheck ./...

build:
	@cd cmd/bot && go build

check: lint
