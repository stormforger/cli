package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// DatasourceCmd is the cobra definition
var DatasourceCmd = &cobra.Command{
	Use:   "datasource",
	Short: "Work with and manage data sources",
	Long: `Work with and manage data sources.

  Currently only a rough validation is implemented.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Fatal("Cannot be run without subcommand, like validate!")
	},
}

func init() {
	RootCmd.AddCommand(DatasourceCmd)
}
