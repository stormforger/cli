package misc

import (
	"fmt"
)

var BuildInfo BuildInfos

type BuildInfos struct {
	Version string
	Time    string
	Commit  string
}

func (buildInfo BuildInfos) String() string {
	return fmt.Sprintf("%v %v (%v - %v) - https://stormforger.com", "forge", BuildInfo.Version, BuildInfo.Time, BuildInfo.Commit)
}

func (buildInfo BuildInfos) ShortString() string {
	return BuildInfo.Version
}

func InitInfo(version string, buildCommit string, buildTime string) {
	BuildInfo = BuildInfos{version, buildCommit, buildTime}
}
