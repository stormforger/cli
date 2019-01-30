package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	testCasePushCmd = &cobra.Command{
		Use:   "push <test-case-file>|all",
		Short: "Creates or updates a test case",
		Long:  `Creates or updates a test case.`,
		Run:   runTestCasePush,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing test case to push (or 'all')")
			}

			fi, err := os.Stat(testCasePushOpts.manifestFile)
			if err != nil || !fi.Mode().IsRegular() {
				log.Fatalf("Manifest file not found: %s", testCasePushOpts.manifestFile)
			}

			manifestData, err := os.OpenFile(testCasePushOpts.manifestFile, os.O_RDONLY, 0755)
			if err != nil {
				log.Fatalf("Could not load manifest: %v", err)
			}

			manifestDefinition, err := loadManifest(manifestData)
			if err != nil {
				log.Fatal(err)
			}
			_, err = manifestDefinition.validate()
			if err != nil {
				log.Fatal(err)
			}

			testCasePushOpts.manifest = manifestDefinition
		},
	}

	testCasePushOpts = struct {
		manifestFile string
		manifest     manifest
	}{}
)

func init() {
	TestCaseCmd.AddCommand(testCasePushCmd)

	// option for manifest support
	testCasePushCmd.Flags().StringVarP(&testCasePushOpts.manifestFile, "manifest", "", "", "Read manifest file (other push flags will be ignored)")
}

func runTestCasePush(cmd *cobra.Command, args []string) {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellowBold := color.New(color.FgYellow, color.Bold).SprintFunc()
	okStatus := green("\u2713")
	failStatus := red("\u2717")
	warnStatus := yellowBold("!")

	// test cases to upsert
	var localFilesToUpload []string
	if len(args) == 1 && args[0] == "all" {
		localFilesToUpload = testCasePushOpts.manifest.allTCPaths()
	} else {
		localFilesToUpload = args
	}

	client := NewClient()
	orgUIDCache := make(map[string]string)
	someErrors := false

	for _, file := range localFilesToUpload {
		found, opts := testCasePushOpts.manifest.LookupTestCase(file)
		if !found {
			log.Fatalf("%s not found", file)
		}

		organizationUID := orgUIDCache[opts.Organisation]
		if organizationUID == "" {
			organizationUID = lookupOrganisationUID(*client, opts.Organisation)
			orgUIDCache[opts.Organisation] = organizationUID
		}

		reader, err := os.OpenFile(opts.Path, os.O_RDONLY, 0755)
		if err != nil {
			log.Fatal(err)
		}

		success, message, errValidation := client.TestCaseUpsert(organizationUID, opts.Name, opts.Comments, file, reader)
		if errValidation != nil {
			log.Fatal(errValidation)
		}

		errorMeta, err := api.UnmarshalErrorMeta(strings.NewReader(message))
		if err != nil {
			log.Fatal(err)
		}

		if rootOpts.OutputFormat == "json" {
			fmt.Println(message)

			if !success || len(errorMeta.Errors) > 0 {
				someErrors = true
			}

			continue
		}

		if success && len(errorMeta.Errors) == 0 {
			fmt.Printf("%s successfully pushed %s to %s/%s\n", okStatus, file, opts.Organisation, opts.Name)
			continue
		}

		if len(errorMeta.Errors) > 0 {
			fail := failStatus
			message := "could not push"
			if success {
				fail = warnStatus
				message = "pushed with warnings"
			}

			someErrors = true

			fmt.Printf("%s %s %s to %s/%s\n", fail, message, file, opts.Organisation, opts.Name)
			fmt.Fprintf(os.Stderr, " %s\n\n", errorMeta.Message)
			for i, e := range errorMeta.Errors {
				fmt.Fprintf(os.Stderr, "%d) %s: %s\n", i+1, e.Code, e.Title)
				fmt.Fprintf(os.Stderr, "%s\n", e.FormattedError)
			}
		} else {
			fmt.Printf("%s successfully pushed %s to %s/%s\n", okStatus, file, opts.Organisation, opts.Name)
		}
	}

	if someErrors {
		if rootOpts.OutputFormat != "json" {
			fmt.Fprintln(os.Stderr, "\nThere have been issues pushing test cases")
		}
		os.Exit(1)
	}
}
