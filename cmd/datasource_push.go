package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	datasourcePushCmd = &cobra.Command{
		Use:   "push <organization-ref> <file>",
		Short: "Upload a file",
		Run:   runDataSourcePush,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Expecting one or more arguments: File(s) to upload")
			}

			if len(args) > 2 && (pushOpts.Name != "" || pushOpts.FieldNames != "") {
				log.Fatal("--name and --fields is not supported for multiple uploads")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}
		},
	}

	pushOpts struct {
		Raw            bool
		Delimiter      string
		FieldNames     string
		Name           string
		NamePrefixPath string
	}
)

func init() {
	datasourceCmd.AddCommand(datasourcePushCmd)

	datasourcePushCmd.Flags().BoolVarP(&pushOpts.Raw, "raw", "r", false, "Upload file as raw fixture")

	datasourcePushCmd.Flags().StringVarP(&pushOpts.Delimiter, "delimiter", "d", "", "Column Delimiter")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Name, "name", "n", "", "Name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.NamePrefixPath, "name-prefix-path", "p", "", "Prefix for name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.FieldNames, "fields", "f", "", "Name for the fields/columns, comma separated")
}

func runDataSourcePush(cmd *cobra.Command, args []string) {
	var fixtureNameFor func(string) string
	if len(args) == 2 {
		fixtureNameFor = func(fileName string) string {
			if pushOpts.Name != "" {
				return pushOpts.Name
			}

			return path.Base(args[1])
		}
	} else {
		fixtureNameFor = func(fileName string) string {
			return path.Base(fileName)
		}
	}

	var fixtureType string
	if pushOpts.Raw {
		fixtureType = "raw"
	} else {
		fixtureType = "structured"
	}

	client := NewClient()

	var params *api.FileFixtureParams
	for _, fileName := range args[1:] {
		fieldNames := pushOpts.FieldNames

		params = &api.FileFixtureParams{
			Name:       pushOpts.NamePrefixPath + fixtureNameFor(fileName),
			Type:       fixtureType,
			FieldNames: fieldNames,
			Delimiter:  pushOpts.Delimiter,
		}

		data, err := os.OpenFile(fileName, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}

		success, result, err := client.PushFileFixture(fileName, data, datasourceOpts.Organisation, params)
		if err != nil {
			log.Fatal(err)
		}

		if !success {
			fmt.Fprintln(os.Stderr, "Could not upload files as new data source!")
			fmt.Fprintln(os.Stderr, string(result))

			os.Exit(1)
		}

		fmt.Println(result)
	}
}
