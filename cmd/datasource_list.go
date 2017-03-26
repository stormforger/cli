package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceListCmd = &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List fixtures",
		Run:     runDataSourceList,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceListCmd)
}

func runDataSourceList(cmd *cobra.Command, args []string) {
	client := NewClient()

	result, err := client.ListFileFixture(datasourceOpts.Organisation)
	if err != nil {
		log.Fatal(err)
	}

	type FileFixtureVersion struct {
		OriginalMd5Hash  string `json:"md5_hash"`
		ProcessedMd5Hash string `json:"processed_md5_hash"`
		FieldNames       string `json:"field_names"`
	}

	type FileFixtureAttributes struct {
		Name           string             `json:"name"`
		CurrentVersion FileFixtureVersion `json:"current_version"`
	}

	type FileFixture struct {
		ID         string                `json:"id"`
		Attributes FileFixtureAttributes `json:"attributes"`
	}

	type ResponseContainer struct {
		Entries []FileFixture `json:"data"`
	}

	var parsedResponse = new(ResponseContainer)
	err = json.Unmarshal(result, &parsedResponse)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range parsedResponse.Entries {
		fmt.Printf("* %s (ID: %s, Content-MD5: %s)\n", item.Attributes.Name, item.ID, item.Attributes.CurrentVersion.OriginalMd5Hash)
	}
}
