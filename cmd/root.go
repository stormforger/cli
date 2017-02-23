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

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
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
