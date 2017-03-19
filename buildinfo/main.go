package buildinfo

import (
	"fmt"
	"log"
)

// Build infos are set during build
var (
	version     string
	buildTime   string
	buildCommit string
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
	log.Println(version)

	BuildInfo = BuildInfos{version, buildCommit, buildTime}
}
