package cmd

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	datasourceDownloadCmd = &cobra.Command{
		Use:              "get <file-name>",
		Aliases:          []string{"download"},
		Short:            "Download file fixture",
		Run:              runDatasourceDownload,
		PersistentPreRun: ensureDatasourceDownloadOptions,
	}

	downloadOpts struct {
		Version string
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceDownloadCmd)

	datasourceDownloadCmd.Flags().StringVarP(&downloadOpts.Version, "version", "v", "current", "Version to download")
}

func ensureDatasourceDownloadOptions(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File name to download")
	}

	datasourceOpts.Organisation = findFirstNonEmpty([]string{datasourceOpts.Organisation, readOrganisationUIDFromFile(), rootOpts.DefaultOrganisation})

	if datasourceOpts.Organisation == "" {
		log.Fatal("Missing organization")
	}
}

func runDatasourceDownload(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[0]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	result, err := client.DownloadFileFixture(datasourceOpts.Organisation, fileFixture.ID, downloadOpts.Version)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}
}
