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

	organisationListOpts struct {
		JSON bool
	}
)

func init() {
	organisationCmd.AddCommand(organisationListCmd)

	organisationListCmd.Flags().BoolVarP(&organisationListOpts.JSON, "json", "", false, "Output machine-readable JSON")
}

func runOrganisationList(cmd *cobra.Command, args []string) {
	client := NewClient()

	result, err := client.ListOrganisations()
	if err != nil {
		log.Fatal(err)
	}

	if organisationListOpts.JSON {
		fmt.Println(string(result))
	} else {
		organisation.ShowNames(bytes.NewReader(result))
	}
}
