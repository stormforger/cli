package cmd

import (
	"fmt"
	"log"
	"path"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	datasourcePushCmd = &cobra.Command{
		Use:   "push <file>",
		Short: "Upload a file",
		Run:   runDataSourcePush,
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
	datasourcePushCmd.Flags().StringVarP(&pushOpts.NamePrefixPath, "name-prefix-prefix", "p", "", "Prefix for name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.FieldNames, "fields", "f", "", "Name for the fields/columns, comma separated")
}

func runDataSourcePush(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("Expecting one or more arguments: File(s) to upload")
	}

	if len(args) > 1 && (pushOpts.Name != "" || pushOpts.FieldNames != "") {
		log.Fatal("--name and --fields is not supported for multiple uploads")
	}

	var fixtureNameFor func(string) string
	if len(args) == 1 {
		fixtureNameFor = func(fileName string) string {
			if pushOpts.Name != "" {
				return pushOpts.Name
			}

			return path.Base(args[0])
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
	for _, fileName := range args {
		fieldNames := pushOpts.FieldNames

		params = &api.FileFixtureParams{
			Name:       pushOpts.NamePrefixPath + fixtureNameFor(fileName),
			Type:       fixtureType,
			FieldNames: fieldNames,
			Delimiter:  pushOpts.Delimiter,
		}

		result, err := client.PushFileFixture(fileName, datasourceOpts.Organisation, params)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	}
}
