package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/filefixture"
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
}

func runDatasourceDownload(cmd *cobra.Command, args []string) {
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
	// TODO how to make this better?
	if fileFixture.ID == "" {
		log.Fatal(fmt.Printf("Filefixture %s not found!", fileName))
	}

	result, err := client.DownloadFileFixture(datasourceOpts.Organisation, fileFixture.ID, downloadOpts.Version)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(os.Stdout, bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}
}
