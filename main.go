package main

import (
	"github.com/stormforger/cli/cmd"
	"github.com/stormforger/cli/misc"
)

// Build infos are set during build
var (
	VERSION      = "0.0.0"
	BUILD_TIME   = "build-time"
	BUILD_COMMIT = "build-commit-sha"
)

func main() {
	misc.InitInfo(VERSION, BUILD_COMMIT, BUILD_TIME)

	cmd.Execute()
}
