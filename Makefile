lint:
	@golangci-lint run

# BuildTime time when the binary was built
# BuildVersion version of the bot itself, as described by most recent git tag
# TODO: Add this one later after a stable release is made
# -X \"main.buildVersion=$$(git describe --abbrev=0 2>/dev/null || echo -n)\" \
# BuildHash short Git commit hash
# BuildBranch Git branch
build:
	@cd cmd/bot && go build -ldflags " \
	-X \"main.buildTime=$$(date +%Y-%m-%dT%H:%M:%S%:z)\" \
	-X \"main.buildHash=$$(git rev-parse --short HEAD)\" \
	-X \"main.buildBranch=$$(git rev-parse --abbrev-ref HEAD)\" \
	" \
	"${@:2}"

check: lint
