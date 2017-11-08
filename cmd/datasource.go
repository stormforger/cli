package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceCmd = &cobra.Command{
		Use:     "datasource",
		Aliases: []string{"ds"},
		Short:   "Work with and manage data sources",
		Long: `Work with and manage data sources.

  Currently only a rough validation is implemented.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal("Cannot be run without subcommand, like validate!")
		},

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if datasourceOpts.Organisation == "" {
				datasourceOpts.Organisation = readOrganisationUIDFromFile()
				if datasourceOpts.Organisation == "" {
					log.Fatal("Missing organization flag")
				}
			}
		},
	}

	datasourceOpts struct {
		Organisation string
	}
)

func init() {
	RootCmd.AddCommand(datasourceCmd)

	datasourceCmd.PersistentFlags().StringVarP(&datasourceOpts.Organisation, "organization", "o", "", "Name of the organization")
}
