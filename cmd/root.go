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
	// RootCmd represents the cobra root command
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
	// ConfigFilename is the forge config file without extension
	ConfigFilename = ".stormforger"
	// EnvPrefix is the prefix for environment configuration
	EnvPrefix = "stormforger"
)

// Execute is the entry function for cobra
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

// NewClient initializes a new API Client
func NewClient() *api.Client {
	return api.NewClient(viper.GetString("endpoint"), viper.GetString("jwt"))
}

/*
	Configuration for JWT can come from (in this order)
	* Environment
	* Configuration ~/.stormforger.toml, ./.stormforger.toml
	* Command line flag
*/
func setupConfig() {
	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv("jwt")
	viper.BindEnv("endpoint")

	viper.SetConfigName(ConfigFilename)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.ReadInConfig()

	viper.BindPFlag("jwt", RootCmd.PersistentFlags().Lookup("jwt"))
	viper.BindPFlag("endpoint", RootCmd.PersistentFlags().Lookup("endpoint"))
}

func init() {
	RootCmd.PersistentFlags().StringVar(&rootOpts.APIEndpoint, "endpoint", "https://api.stormforger.com", "API Endpoint")
	RootCmd.PersistentFlags().StringVar(&rootOpts.JWT, "jwt", "", "JWT access token")
}
