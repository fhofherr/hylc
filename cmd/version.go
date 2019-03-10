package cmd

var (
	// BuildTime is the time this version of hylc was built.
	BuildTime string
	// GitHash is the git commit hash of this version of hylc.
	GitHash string
	// Version is the git tag of this version of hylc. It may be empty
	// if hylc was built from a non-tagged commit.
	Version string
)
