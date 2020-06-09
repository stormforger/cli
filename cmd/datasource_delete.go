package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceDeleteCmd = &cobra.Command{
		Use:     "rm <organisation-ref> <name>",
		Aliases: []string{"delete", "remove"},
		Short:   "Delete a fixture",
		Run:     runDatasourceDelete,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 1 {
				log.Fatal("Missing organisation")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceDeleteCmd)
}

func runDatasourceDelete(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[1]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	success, result, err := client.DeleteFileFixture(fileFixture.ID, datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Println(result)
	}
}
