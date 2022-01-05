// Package buildinfo provides high-level build information injected during
// build.
package buildinfo

var (
	// BuildID is the unique build identifier.
	BuildID string = "unknown"

	// BuildTag is the git tag from which this build was created.
	BuildTag string = "unknown"
)
