lint:
	@golangci-lint run

# BuildTime time when the binary was built
# BuildVersion version of the bot itself, as described by most recent git tag
# BuildHash short Git commit hash
# BuildBranch Git branch
build:
	@cd cmd/bot && go build -ldflags " \
	-X \"main.buildTime=$$(date +%Y-%m-%dT%H:%M:%S%:z)\" \
	-X \"main.buildVersion=$$(git describe --abbrev=0 2>/dev/null || echo -n)\" \
	-X \"main.buildHash=$$(git rev-parse --short HEAD)\" \
	-X \"main.buildBranch=$$(git rev-parse --abbrev-ref HEAD)\" \
	" \
	"${@:2}"

check: lint
