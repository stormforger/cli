package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	datasourceMoveCmd = &cobra.Command{
		Use:     "mv <organization-ref> <name> <new-name>",
		Aliases: []string{"move", "rename"},
		Short:   "Rename a fixture",
		Run:     runDataSourceMove,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) != 3 {
				log.Fatal("Expecting exactly three arguments: organisation, name of source and destination")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}
		},
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceMoveCmd)
}

func runDataSourceMove(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[1]
	newFileName := args[2]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	success, result, err := client.MoveFileFixture(datasourceOpts.Organisation, fileFixture.ID, newFileName)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		fmt.Fprintln(os.Stderr, "Could not move data source!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}
}
