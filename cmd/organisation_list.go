package cmd

import (
	"bytes"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/organisation"
)

var (
	organisationListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{},
		Short:   "Manage organisations",
		Run:     runOrganisationList,
	}
)

func init() {
	organisationCmd.AddCommand(organisationListCmd)
}

func runOrganisationList(cmd *cobra.Command, args []string) {
	client := NewClient()

	result, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
	}

	organisation.ShowNames(bytes.NewReader(result))
}
