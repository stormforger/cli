package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceShowCmd = &cobra.Command{
		Use:              "show <file-name>",
		Aliases:          []string{},
		Short:            "Show details of fixture",
		Run:              runDatasourceShow,
		PersistentPreRun: ensureDatasourceShowOptions,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceShowCmd)
}

func ensureDatasourceShowOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File name of fixture to show details about")
	}

	datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

	if datasourceOpts.Organisation == "" {
		log.Fatal("Missing organization")
	}
}

func runDatasourceShow(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	filefixture.ShowDetails(fileFixture)
}
