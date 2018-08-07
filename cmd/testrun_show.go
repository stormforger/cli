package cmd

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/testrun"
)

var (
	// testRunShowCmd represents the calllog command
	testRunShowCmd = &cobra.Command{
		Use:   "show <test-run-ref>",
		Short: "Show test run details",
		Long:  `Show test run details.`,
		Run:   testRunShow,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("Expect exactly one argument: test run reference!")
			}
		},
	}

	testRunShowOpts struct {
		Type       string
		Full       bool
		OutputFile string
	}
)

func init() {
	TestRunCmd.AddCommand(testRunShowCmd)
}

func testRunShow(cmd *cobra.Command, args []string) {
	client := NewClient()

	result := fetchTestRun(*client, args[0])

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	testRun, err := testrun.UnmarshalSingle(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	// FIXME can we integrate this into testrun.UnmarshalSingle somehow?
	meta, err := api.UnmarshalMeta(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s, %s)\n", testRun.Scope, testRun.State, testRun.ID)
	fmt.Printf("Report  %s\n", meta.Links.SelfWeb)
	if testRun.Title != "" {
		fmt.Printf("Title   %s\n", testRun.Title)
	}
	fmt.Printf("Started %s\n", testRun.StartedAt)
	fmt.Printf("Ended   %s\n", testRun.EndedAt)
	fmt.Print("Configuration:\n")
	fmt.Printf("  Setup: %s cluster in %s\n", testRun.TestConfiguration.ClusterSizing, testRun.TestConfiguration.ClusterRegion)

	if testRun.TestConfiguration.DisableGzip {
		fmt.Print("  [\u2713] Disabled GZIP\n")
	}
	if testRun.TestConfiguration.SkipWait {
		fmt.Print("  [\u2713] Skip Waits\n")
	}
	if testRun.TestConfiguration.DumpTrafficFull {
		fmt.Print("  [\u2713] Traffic Dump\n")
	}
	if testRun.TestConfiguration.SessionValidationMode {
		fmt.Print("  [\u2713] Session Validation Mode\n")
	}

	if testRun.Notes != "" {
		fmt.Printf("\nNotes:\n%s\n", testRun.Notes)
	}
}
