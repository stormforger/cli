package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/filefixture"

	"github.com/spf13/cobra"
)

var (
	datasourcePushCmd = &cobra.Command{
		Use:   "push <organization-ref> <file>",
		Short: "Upload a file",
		Run:   runDataSourcePush,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// manifest file was given?
			if manifestFile != "" {
				if len(args) < 1 {
					log.Fatal("Missing files to push (or 'all')")
				}

				fi, err := os.Stat(manifestFile)
				if err != nil || !fi.Mode().IsRegular() {
					log.Fatalf("Manifest file not found: %s", manifestFile)
				}

				manifestData, err := os.OpenFile(manifestFile, os.O_RDONLY, 0755)
				if err != nil {
					log.Fatalf("Could not load manifest: %v", err)
				}

				manifestDefinition, err = loadManifest(manifestData)
				if err != nil {
					log.Fatal(err)
				}
				_, err = manifestDefinition.validate()
				if err != nil {
					log.Fatal(err)
				}

				useManifest = true

				return
			}

			// no manifest given, use arguments
			if len(args) < 2 {
				log.Fatal("Expecting one or more arguments: File(s) to upload (if --manifest not provided)")
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

			datasourceOpts.Organisation = lookupOrganisationUID(*NewClient(), args[0])
			if datasourceOpts.Organisation == "" {
				log.Fatal("Missing organization")
			}
		},
	}

	pushOpts pushCmdOpts

	useManifest        bool
	manifestFile       string
	manifestDefinition manifest
)

type pushCmdOpts struct {
	Raw             bool
	Delimiter       string
	FieldNames      string
	Name            string
	NamePrefixPath  string
	FirstRowHeaders bool
	Organisation    string
	OrganisationUID string
}

func init() {
	datasourceCmd.AddCommand(datasourcePushCmd)

	// type of FF
	datasourcePushCmd.Flags().BoolVarP(&pushOpts.Raw, "raw", "r", false, "Upload file as raw fixture")

	// general options
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Name, "name", "n", "", "Name for the file fixture")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.NamePrefixPath, "name-prefix-path", "p", "", "Prefix for name for the file fixture")

	// options for structured FF
	datasourcePushCmd.Flags().StringVarP(&pushOpts.Delimiter, "delimiter", "d", "", "Column Delimiter")
	datasourcePushCmd.Flags().StringVarP(&pushOpts.FieldNames, "fields", "f", "", "Name for the fields/columns, comma separated (,)")
	datasourcePushCmd.Flags().BoolVarP(&pushOpts.FirstRowHeaders, "auto-field-names", "", false, "Interpret first row as headers")

	// option for manifest support
	datasourcePushCmd.Flags().StringVarP(&manifestFile, "manifest", "", "", "Read manifest file (other push flags will be ignored)")
}

func lookupPushOpts(localFileName string, args []string) pushCmdOpts {
	if !useManifest {
		if pushOpts.Name == "" {
			pushOpts.Name = pushOpts.NamePrefixPath + path.Base(localFileName)
		}

		if pushOpts.Organisation == "" {
			pushOpts.Organisation = datasourceOpts.Organisation
		}

		return pushOpts
	}

	found, ds := manifestDefinition.lookupDataSource(localFileName)
	if !found || ds.Path != "" {

		uploadName := ds.Name
		if uploadName == "" {
			uploadName = path.Base(localFileName)
		}

		po := pushCmdOpts{
			Name:            uploadName,
			Raw:             ds.Raw,
			Delimiter:       ds.Delimiter,
			FieldNames:      strings.Join(ds.Fields, ","),
			FirstRowHeaders: ds.AutoColumnNames,
			OrganisationUID: lookupOrganisationUID(*NewClient(), ds.Organisation),
		}

		if po.OrganisationUID == "" {
			log.Fatalf("Organization %s not found\n", ds.Organisation)
		}

		return po
	}

	log.Fatalf("Manifest: Definition for data source %v not found\n", localFileName)

	return pushCmdOpts{}
}

func runDataSourcePush(cmd *cobra.Command, args []string) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	okStatus := green("\u2713")
	failStatus := red("\u2717")

	// files to upload
	var localFilesToUpload []string
	if !useManifest {
		localFilesToUpload = args[1:]
	} else {
		if len(args) == 1 && args[0] == "all" {
			localFilesToUpload = manifestDefinition.allDSPaths()
		} else {
			localFilesToUpload = args
		}
	}

	var params *api.FileFixtureParams
	someErrors := false
	client := NewClient()
	for _, localFileName := range localFilesToUpload {
		options := lookupPushOpts(localFileName, args)

		params = &api.FileFixtureParams{
			Name:            options.Name,
			Raw:             options.Raw,
			FieldNames:      options.FieldNames,
			Delimiter:       options.Delimiter,
			FirstRowHeaders: options.FirstRowHeaders,
		}

		data, err := os.OpenFile(localFileName, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}

		success, result, err := client.PushFileFixture(localFileName, data, options.OrganisationUID, params)
		if err != nil {
			log.Fatal(err)
		}

		if !success {
			fmt.Fprintf(os.Stderr, "%s %v could not be create or updated! Error: %s\n", failStatus, params.Name, string(result))

			someErrors = true
			continue
		}

		if rootOpts.OutputFormat == "json" {
			fmt.Println(string(result))
		} else {
			if !options.Raw {
				fixture, err := filefixture.UnmarshalFileFixture(bytes.NewReader(result))
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf(
					"%s %s successfully uploaded! Name: %s, Items: %d, Columns: %s\n",
					okStatus,
					params.Name,
					fixture.Name,
					fixture.CurrentVersion.ItemCount,
					strings.Join(fixture.CurrentVersion.FieldNames, ", "),
				)
			} else {
				fmt.Printf("%s %s uploaded!\n", okStatus, params.Name)
			}
		}
	}

	if someErrors {
		os.Exit(1)
	}
}
