package cmd

import (
	"bytes"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceListCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List fixtures",
		Run:     runDataSourceList,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceListCmd)
}

func runDataSourceList(cmd *cobra.Command, args []string) {
	client := NewClient()

	result, err := client.ListFileFixture(datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	filefixture.ShowName(bytes.NewReader(result))
}
