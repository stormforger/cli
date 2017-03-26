package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceDeleteCmd = &cobra.Command{
		Use:     "rm <file-uid>",
		Aliases: []string{"delete", "remove"},
		Short:   "Delete a fixture",
		Run:     runDatasourceDelete,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceDeleteCmd)
}

func runDatasourceDelete(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File UID to delete")
	}

	fileUID := args[0]
	client := NewClient()

	result, err := client.DeleteFileFixture(fileUID, datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
