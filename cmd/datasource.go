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

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organization")
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
