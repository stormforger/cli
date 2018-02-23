package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceMoveCmd = &cobra.Command{
		Use:              "mv",
		Aliases:          []string{"move", "rename"},
		Short:            "Rename a fixture",
		Run:              runDataSourceMove,
		PersistentPreRun: ensureDatasourceMoveOptions,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceMoveCmd)
}

func ensureDatasourceMoveOptions(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("Expecting exactly two arguments: name of source and destination")
	}

	datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

	if datasourceOpts.Organisation == "" {
		log.Fatal("Missing organization")
	}
}

func runDataSourceMove(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]
	newFileName := args[1]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	result, err := client.MoveFileFixture(datasourceOpts.Organisation, fileFixture.ID, newFileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
