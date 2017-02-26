package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stormforger/cli/api"
)

var (
	RootCmd = &cobra.Command{
		Use:   "forge",
		Short: "Command line client to StormForger (https://stormforger.com)",
		Long: `The command line client "forge" to StormForger offers a interface
to the StormForger API and several convinience methods
to work on load and performance tests.

Happy Load Testing :)`,
	}

	rootOpts struct {
		APIEndpoint string
		JWTToken    string
	}
)

const (
	CONFIG_FILENAME = ".stormforger"
	ENV_PREFIX      = "stormforger"
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

// Execute
// searchs for JWT-token,
// adds all child commands to the root command
// and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string, buildTime string, buildCommit string) {
	buildInfo = BuildInfo{version, buildTime, buildCommit}

	// most important thing is the jwt
	//viper.SetDefault("jwt", "")
	findJwt()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if viper.GetString("jwt") == "" {
		color.Yellow("\nNo JWT token in config file, environment or via command line flag!\nPlease provide JWT token to talk to the StormForger API. See README.md.\n")
	} else {
		color.Green("\nUsing token: %s\n", viper.GetString("jwt"))
	}
}

func NewClient() *api.Client {
	return api.NewClient(rootOpts.APIEndpoint, viper.GetString("jwt"))
}

// JWT token must be either in one of these locations
// .stormforger.toml
// $HOMR/.stormforger.toml
// $ENV["STORMFORGER_JWT"]
// via cli flag --jwt-token
func findJwt() bool {
	// ENV overrides config
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.BindEnv("jwt")

	// config
	viper.SetConfigName(CONFIG_FILENAME) // name of config file (without extension)
	viper.AddConfigPath("$HOME")         // call multiple times to add many search paths
	viper.AddConfigPath(".")             // optionally look for config in the working directory
	err := viper.ReadInConfig()          // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		// no config
		// TODO I don't think that we inform the user at this stage about a missing config file
		// log.Println("No config file: %s \n", err)
	}

	// command line flag overrides config and ENV
	viper.BindPFlag("jwt", RootCmd.Flags().Lookup("jwt-token"))

	return true
}

func init() {
	RootCmd.PersistentFlags().StringVar(&rootOpts.APIEndpoint, "api-endpoint", "https://api.stormforger.com", "API Endpoint")
	RootCmd.Flags().StringVar(&rootOpts.JWTToken, "jwt-token", "", "JWT Token")
}
