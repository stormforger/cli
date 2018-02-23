package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceDeleteCmd = &cobra.Command{
		Use:              "rm <file-name>",
		Aliases:          []string{"delete", "remove"},
		Short:            "Delete a fixture",
		Run:              runDatasourceDelete,
		PersistentPreRun: ensureDatasourceDeleteOptions,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceDeleteCmd)
}

func ensureDatasourceDeleteOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File name to delete")
	}

	datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

	if datasourceOpts.Organisation == "" {
		log.Fatal("Missing organization")
	}
}

func runDatasourceDelete(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	result, err := client.DeleteFileFixture(fileFixture.ID, datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
