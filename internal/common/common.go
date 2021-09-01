// common is a Singleton package that contains values filled in while building with Makefile
// values are initialized at the application startup and can be shared across all packages
package common

import "fmt"

var (
	// Values filled in with ./build.sh (ldflags)

	// BuildTime time when the binary was built
	BuildTime string

	// BuildVersion version of the bot itself, as described by most recent git tag
	BuildVersion string

	// BuildHash short Git commit hash
	BuildHash string

	// BuildBranch Git branch
	BuildBranch string
)

func Version() string {
	return fmt.Sprintf("%s %s@%s", BuildVersion, BuildBranch, BuildHash)
}
