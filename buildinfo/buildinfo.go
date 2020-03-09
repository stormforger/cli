package buildinfo

import (
	"fmt"
	"runtime"
)

// Build infos are set during build
var (
	version string = "dev"
	date    string
	commit  string
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
	return fmt.Sprintf("%v %v (%v, %v, %v) - https://stormforger.com", "forge", BuildInfo.Version, BuildInfo.Commit, BuildInfo.Time, runtime.Version())
}

// ShortString only returns the build version
func (buildInfo BuildInfos) ShortString() string {
	return BuildInfo.Version
}

func init() {
	BuildInfo = BuildInfos{version, date, commit}
}
