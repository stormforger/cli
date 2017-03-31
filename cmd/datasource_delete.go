package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceDeleteCmd = &cobra.Command{
		Use:              "rm <file-uid>",
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
		log.Fatal("Expecting exactly one argument: File UID to delete")
	}
}

func runDatasourceDelete(cmd *cobra.Command, args []string) {
	fileUID := args[0]
	client := NewClient()

	result, err := client.DeleteFileFixture(fileUID, datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
