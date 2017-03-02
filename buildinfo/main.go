package buildinfo

import (
	"fmt"
)

// Build infos are set during build
var (
	Version     = "0.0.42"
	BuildTime   = "build-time"
	BuildCommit = "build-commit-sha"
)

// BuildInfo holds build information like version, build time and commit
var BuildInfo BuildInfos

// BuildInfos struct for holds build information like version, build time and commit
type BuildInfos struct {
	Version string
	Time    string
	Commit  string
}

// String returns the version, build time and commit
func (buildInfo BuildInfos) String() string {
	return fmt.Sprintf("%v %v (%v - %v) - https://stormforger.com", "forge", BuildInfo.Version, BuildInfo.Time, BuildInfo.Commit)
}

// ShortString only returns the build version
func (buildInfo BuildInfos) ShortString() string {
	return BuildInfo.Version
}

func init() {
	BuildInfo = BuildInfos{Version, BuildCommit, BuildTime}
}
