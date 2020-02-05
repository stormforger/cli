package cmd

import (
	"fmt"
	"log"

	humanize "github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
)

var (
	datasourceShowCmd = &cobra.Command{
		Use:     "show <organisation-ref> <name>",
		Aliases: []string{},
		Short:   "Show details of fixture",
		Run:     runDatasourceShow,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 2 {
				log.Fatal("Missing organisation or datasource")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceShowCmd)
}

func runDatasourceShow(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[1]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	ShowDetails(fileFixture)
}

// ShowDetails print out details of a file fixture, including its current version
func ShowDetails(fileFixture *filefixture.FileFixture) {
	fmt.Printf("Name:            %s\n", fileFixture.Name)
	fmt.Printf("UID:             %s\n", fileFixture.ID)
	fmt.Printf("Created:         %s (%s)\n", convertToLocalTZ(fileFixture.CreatedAt), humanize.Time(fileFixture.CreatedAt))
	fmt.Printf("Updated:         %s (%s)\n", convertToLocalTZ(fileFixture.UpdatedAt), humanize.Time(fileFixture.UpdatedAt))
	fmt.Printf("Current Version: %s\n", fileFixture.CurrentVersion.ID)
	fmt.Printf("  SHA256 Hash:   %s\n", fileFixture.CurrentVersion.Hash)
	fmt.Printf("  Size:          %s\n", humanize.Bytes(uint64(fileFixture.CurrentVersion.FileSize)))
	fmt.Printf("  Line Count:    %v\n", fileFixture.CurrentVersion.ItemCount)
	fmt.Printf("  Created:       %s (%s)\n", convertToLocalTZ(fileFixture.CurrentVersion.CreatedAt), humanize.Time(fileFixture.CurrentVersion.CreatedAt))
	fmt.Printf("Version(s):      %v\n", len(fileFixture.Versions))
	for _, version := range fileFixture.Versions {
		fmt.Printf("  - %s (created %s)\n", version.ID, convertToLocalTZ(version.CreatedAt))
	}
}
