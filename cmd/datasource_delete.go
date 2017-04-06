package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceDeleteCmd = &cobra.Command{
		Use:              "rm <file-name>",
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
		log.Fatal("Expecting exactly one argument: File name to delete")
	}
}

func runDatasourceDelete(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]

	fileFixtureListResponse, err := client.ListFileFixture(datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fileFixtures, err := filefixture.UnmarshalFileFixtures(bytes.NewReader(fileFixtureListResponse))
	if err != nil {
		log.Fatal(err)
	}

	fileFixture := fileFixtures.FindByName(fileName)
	if fileFixture == *new(filefixture.FileFixture) {
		log.Fatal(fmt.Printf("Filefixture %s not found!", fileName))
	}

	result, err := client.DeleteFileFixture(fileFixture.ID, datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
