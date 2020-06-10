package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/stormforger/cli/api/filefixture"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	datasourcePushCmd = &cobra.Command{
		Use:   "push <organisation-ref> <file>",
		Short: "Upload a file",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				log.Fatal("Expecting one or more arguments: organisation and file(s) to upload")
			}

			if pushOpts.Raw && (pushOpts.FieldNames != "" || pushOpts.Delimiter != "" || pushOpts.FirstRowHeaders) {
				log.Fatal("Raw file fixtures do not support --fields, --delimiter and --auto-field-names")
			}

			if len(args) > 2 && (pushOpts.Name != "" || pushOpts.FieldNames != "") {
				log.Fatal("--name and --fields is not supported for multiple uploads")
			}

			if pushOpts.FieldNames != "" && pushOpts.FirstRowHeaders {
				log.Fatal("--fields and --auto-field-names are mutual exclusive")
			}

			if pushOpts.Delimiter != "" && utf8.RuneCountInString(pushOpts.Delimiter) > 1 {
				log.Fatal("Delimiter can only be one character!")
			}

			datasourceOpts.Organisation = lookupOrganisationUID(NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organisation")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			client := NewClient()
			MainDataSourcePush(client, pushOpts, args)
		},
	}

	pushOpts PushOpts
)

type PushOpts struct {
	Raw             bool
	Delimiter       string
	FieldNames      string
	Name            string
	NamePrefixPath  string
	FirstRowHeaders bool
}

func init() {
	datasourceCmd.AddCommand(datasourcePushCmd)

	// type of FF
	datasourcePushCmd.Flags().BoolVarP(&pushOpts.Raw, "raw", "r", false, "Enable to upload file as raw fixture")

	// general options
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Name, "name", "n", "", "Name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.NamePrefixPath, "name-prefix-path", "p", "", "Prefix for name for the file fixture")

	// options for structured FF
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Delimiter, "delimiter", "d", "", "Column Delimiter (defaults to ',' for structured file fixtures)")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.FieldNames, "fields", "f", "", "Name for the fields/columns, comma separated (,)")
	datasourcePushCmd.Flags().BoolVar(&pushOpts.FirstRowHeaders, "auto-field-names", false, "Enable to interpret first row as headers")
}

func MainDataSourcePush(client *api.Client, pushOpts PushOpts, args []string) {
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

	var params *api.FileFixtureParams
	someErrors := false
	for _, fileName := range args[1:] {
		fieldNames := pushOpts.FieldNames

		params = &api.FileFixtureParams{
			Name:            pushOpts.NamePrefixPath + fixtureNameFor(fileName),
			Type:            fixtureType,
			FieldNames:      fieldNames,
			Delimiter:       pushOpts.Delimiter,
			FirstRowHeaders: pushOpts.FirstRowHeaders,
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
			fmt.Fprintf(os.Stderr, "%v could not be create or updated! Error: %s\n", params.Name, string(result))

			someErrors = true
			continue
		}

		if rootOpts.OutputFormat == "json" {
			fmt.Println(string(result))
		} else {
			if fixtureType == "structured" {
				fixture, err := filefixture.UnmarshalFileFixture(bytes.NewReader(result))
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf(
					"%s successfully uploaded! Name: %s, Items: %d, Columns: %s\n",
					params.Name,
					fixture.Name,
					fixture.CurrentVersion.ItemCount,
					strings.Join(fixture.CurrentVersion.FieldNames, ", "),
				)
			} else {
				fmt.Printf("%s uploaded!\n", params.Name)
			}
		}
	}

	if someErrors {
		os.Exit(1)
	}
}
