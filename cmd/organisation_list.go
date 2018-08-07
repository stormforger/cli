package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/organisation"
)

var (
	organisationListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{},
		Short:   "List organisations you have access to",
		Run:     runOrganisationList,
	}
)

func init() {
	organisationCmd.AddCommand(organisationListCmd)
}

func runOrganisationList(cmd *cobra.Command, args []string) {
	client := NewClient()

	success, result, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
	}
	if !success {
		log.Fatal(string(result))
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	items, err := organisation.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items.Organisations {
		if rootOpts.OutputFormat == "human" {
			fmt.Printf("%s (ID: %s)\n", item.Name, item.ID)
		} else {
			fmt.Printf("%s\n", item.Name)
		}
	}
}
