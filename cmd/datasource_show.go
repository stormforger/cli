package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceShowCmd = &cobra.Command{
		Use:     "show <organisation-ref> <name>",
		Aliases: []string{},
		Short:   "Show details of fixture",
		Run:     runDatasourceShow,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 2 {
				log.Fatal("Missing organisation or datasource")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceShowCmd)
}

func runDatasourceShow(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[1]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	filefixture.ShowDetails(fileFixture)
}
