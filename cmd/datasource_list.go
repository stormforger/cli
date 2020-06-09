package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceListCmd = &cobra.Command{
		Use:     "list <organisation-ref>",
		Aliases: []string{"ls"},
		Short:   "List fixtures",
		Run:     runDataSourceList,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
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
	datasourceCmd.AddCommand(datasourceListCmd)
}

func runDataSourceList(cmd *cobra.Command, args []string) {
	client := NewClient()

	success, result, err := client.ListFileFixture(datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		fmt.Fprintln(os.Stderr, "Could not list data sources!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	items, err := filefixture.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items.Fixtures {
		if rootOpts.OutputFormat == "human" {
			fmt.Printf("%s (ID: %s)\n", item.Name, item.ID)
		} else {
			fmt.Printf("%s\n", item.Name)
		}
	}
}
