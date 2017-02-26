package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/stormforger/cli/api"
)

var (
	RootCmd = &cobra.Command{
		Use:   "forge",
		Short: "This is the StormForger command line client!",
		Long:  `This is the StormForger command line client!`,
	}

	rootOpts struct {
		APIEndpoint string
	}
)

type BuildInfo struct {
	Version string
	Time    string
	Commit  string
}

func (buildInfo BuildInfo) String() string {
	return fmt.Sprintf("%v %v (%v - %v) - https://stormforger.com", RootCmd.Use, buildInfo.Version, buildInfo.Time, buildInfo.Commit)
}

func (buildInfo BuildInfo) ShortString() string {
	return buildInfo.Version
}

var buildInfo BuildInfo

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, buildTime string, buildCommit string) {
	buildInfo = BuildInfo{version, buildTime, buildCommit}

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func NewClient() *api.Client {
	return api.NewClient(rootOpts.APIEndpoint, readJwt())
}

func readJwt() string {
	jwt, _ := ioutil.ReadFile("./.stormforger_jwt")

	return strings.TrimSpace(string(jwt))
}

func init() {
	RootCmd.PersistentFlags().StringVar(&rootOpts.APIEndpoint, "api-endpoint", "https://api.stormforger.com", "API Endpoint")
}
