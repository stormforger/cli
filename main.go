package main

import "github.com/stormforger/cli/cmd"

// Build infos are set during build
var (
	VERSION      = "0.0.0"
	BUILD_TIME   = "build-time"
	BUILD_COMMIT = "build-commit-sha"
)

func main() {
	cmd.Execute(VERSION, BUILD_TIME, BUILD_COMMIT)
}
