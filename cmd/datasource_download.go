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
		Use:     "get <organisation-ref> <name>",
		Aliases: []string{"download"},
		Short:   "Download file fixture",
		Args:    cobra.ExactArgs(2),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			datasourceOpts.Organisation = lookupOrganisationUID(NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
		Run:               runDatasourceDownload,
		ValidArgsFunction: completeOrga,
	}

	downloadOpts struct {
		Version string
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceDownloadCmd)

	datasourceDownloadCmd.Flags().StringVarP(&downloadOpts.Version, "version", "v", "current", "Version to download")
}

func runDatasourceDownload(cmd *cobra.Command, args []string) {
	client := NewClient()
	fileName := args[1]

	fileFixture := findFixtureByName(*client, datasourceOpts.Organisation, fileName)

	success, result, err := client.DownloadFileFixture(datasourceOpts.Organisation, fileFixture.ID, downloadOpts.Version)
	if err != nil {
		log.Fatal(err)
	}

	if !success {
		log.Fatalf("Could not download %s: %s\n", fileName, result)
	}

	_, err = io.Copy(os.Stdout, bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}
}
