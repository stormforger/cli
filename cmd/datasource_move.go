package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceMoveCmd = &cobra.Command{
		Use:              "mv",
		Aliases:          []string{"move", "rename"},
		Short:            "Rename a fixture",
		Run:              runDataSourceMove,
		PersistentPreRun: ensureDatasourceMoveOptions,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceMoveCmd)
}

func ensureDatasourceMoveOptions(cmd *cobra.Command, args []string) {
	if len(args) != 2 {
		log.Fatal("Expecting exactly two arguments: name of source and destination")
	}
}

func runDataSourceMove(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]
	newFileName := args[1]

	fileFixtureListResponse, err := client.ListFileFixture(datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	fileFixtures, err := filefixture.UnmarshalFileFixtures(bytes.NewReader(fileFixtureListResponse))
	if err != nil {
		log.Fatal(err)
	}

	fileFixture := fileFixtures.FindByName(fileName)
	// TODO how to make this better?
	if fileFixture.ID == "" {
		log.Fatal(fmt.Printf("Filefixture %s not found!", fileName))
	}

	result, err := client.MoveFileFixture(datasourceOpts.Organisation, fileFixture.ID, newFileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
