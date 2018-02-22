package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceShowCmd = &cobra.Command{
		Use:              "show <file-name>",
		Aliases:          []string{},
		Short:            "Show details of fixture",
		Run:              runDatasourceShow,
		PersistentPreRun: ensureDatasourceShowOptions,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceShowCmd)
}

func ensureDatasourceShowOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File name of fixture to show details about")
	}

	datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

	if datasourceOpts.Organisation == "" {
		log.Fatal("Missing organization")
	}
}

func runDatasourceShow(cmd *cobra.Command, args []string) {
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
	if fileFixture.ID == "" {
		log.Fatal(fmt.Printf("Filefixture %s not found!\n", fileName))
	}

	fileFixtureUID := fileFixture.ID

	fileFixtureResponse, err := client.GetFileFixture(datasourceOpts.Organisation, fileFixtureUID)
	if err != nil {
		log.Fatal(err)
	}

	fileFixtureFoo, err := filefixture.UnmarshalFileFixture(bytes.NewReader(fileFixtureResponse))
	if err != nil {
		log.Fatal(err)
	}

	filefixture.ShowDetails(fileFixtureFoo)
}
