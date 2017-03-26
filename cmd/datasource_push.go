package cmd

import (
	"fmt"
	"log"

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
		Raw        bool
		Delimiter  string
		FieldNames string
		Name       string
	}
)

func init() {
	datasourceCmd.AddCommand(datasourcePushCmd)

	datasourcePushCmd.Flags().BoolVarP(&pushOpts.Raw, "raw", "r", false, "Upload file as raw fixture")

	datasourcePushCmd.Flags().StringVarP(&pushOpts.Delimiter, "delimiter", "d", "", "Column Delimiter")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Name, "name", "n", "", "Name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.FieldNames, "fields", "f", "", "Name for the fields/columns, comma separated")
}

func runDataSourcePush(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Expecting exactly one argument: File Name to upload")
	}

	fileName := args[0]
	client := NewClient()

	var fixtureType string
	if pushOpts.Raw {
		fixtureType = "raw"
	} else {
		fixtureType = "structured"
	}

	params := &api.FileFixtureParams{
		Name:       pushOpts.Name,
		Type:       fixtureType,
		FieldNames: pushOpts.FieldNames,
		Delimiter:  pushOpts.Delimiter,
	}

	result, err := client.PushFileFixture(fileName, datasourceOpts.Organisation, params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
