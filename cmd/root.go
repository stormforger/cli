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
to the StormForger API and several convenience methods
to handle load and performance tests.

Happy Load Testing :)`,
	}

	rootOpts struct {
		APIEndpoint string
		JWT         string
	}
)

const (
	CONFIG_FILENAME = ".stormforger"
	ENV_PREFIX      = "stormforger"
)

// Execute
// searchs for JWT-token,
// adds all child commands to the root command
// and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	setupConfig()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	if viper.GetString("jwt") == "" {
		color.Yellow("\nNo JWT token in config file, environment or via command line flag!\n")
	}
}

func NewClient() *api.Client {
	return api.NewClient(rootOpts.APIEndpoint, viper.GetString("jwt"))
}

/*
	Configuration for JWT can come from (in this order)
	* Environment
	* Configuration ~/.stormforger.toml, ./.stormforger.toml
	* Command line flag
*/
func setupConfig() {
	viper.SetEnvPrefix(ENV_PREFIX)
	viper.BindEnv("jwt")

	viper.SetConfigName(CONFIG_FILENAME)
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	viper.BindPFlag("jwt", RootCmd.Flags().Lookup("jwt"))
}

func init() {
	RootCmd.PersistentFlags().StringVar(&rootOpts.APIEndpoint, "endpoint", "https://api.stormforger.com", "API Endpoint")
	RootCmd.Flags().StringVar(&rootOpts.JWT, "jwt", "", "JWT access token")
}
