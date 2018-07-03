package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/jsonapi"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/testrun"
)

var (
	// testRunLaunchCmd represents the test run launch command
	testRunLaunchCmd = &cobra.Command{
		Use:   "launch <test-case-ref>",
		Short: "Create and launch a new test run",
		Long: `Create and launch a new test run based on given test case

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.

Examples
--------
* launch by organisation and test case name

  forge test-case launch acme-inc/checkout

* alternatively the test case UID can also be provided

  forge test-case launch xPSX5KXM

`,
		Run: testRunLaunch,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			if len(args) > 1 {
				log.Fatal("Too many arguments")
			}
		},
	}

	testRunLaunchOpts struct {
		Title        string
		Notes        string
		Watch        bool
		MaxWatchTime time.Duration
		CheckNFR     string
		DisableGzip  bool
		SkipWait     bool
		DumpTraffic  bool
	}
)

func init() {
	TestCaseCmd.AddCommand(testRunLaunchCmd)

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Title, "title", "t", "", "Descriptive title of test run")
	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Notes, "notes", "n", "", "Longer description (Markdown supported)")

	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.Watch, "watch", "w", false, "Automatically watch newly launched test run")
	testRunLaunchCmd.Flags().DurationVar(&testRunLaunchOpts.MaxWatchTime, "watch-timeout", 0, "Maximum duration in seconds to watch")

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.CheckNFR, "nfr-check-file", "", "", "Check test result against NFR definition (implies --watch)")

	// options for debugging
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.DisableGzip, "disable-gzip", "", false, "Globally disable gzip")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.SkipWait, "skip-wait", "", false, "Ignore defined waits")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.DumpTraffic, "dump-traffic", "", false, "Create traffic dump")
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := lookupTestCase(*client, args[0])

	launchOptions := api.TestRunLaunchOptions{
		Title:       testRunLaunchOpts.Title,
		Notes:       testRunLaunchOpts.Notes,
		DisableGzip: testRunLaunchOpts.DisableGzip,
		SkipWait:    testRunLaunchOpts.SkipWait,
		DumpTraffic: testRunLaunchOpts.DumpTraffic,
	}

	status, response, err := client.TestRunCreate(testCaseUID, launchOptions)
	if err != nil {
		log.Fatal(err)
	}

	if status {
		if testRunLaunchOpts.Watch || testRunLaunchOpts.CheckNFR != "" {
			testRun := new(testrun.TestRun)
			err = jsonapi.UnmarshalPayload(strings.NewReader(response), testRun)
			if err != nil {
				log.Fatal(err)
			}

			watchTestRun(testRun.ID, testRunLaunchOpts.MaxWatchTime.Round(time.Second).Seconds(), rootOpts.OutputFormat)

			if testRunLaunchOpts.CheckNFR != "" {
				fmt.Println("Test finished, running non-functional checks...")

				fileName := filepath.Base(testRunLaunchOpts.CheckNFR)
				nfrData, err := os.OpenFile(testRunLaunchOpts.CheckNFR, os.O_RDONLY, 0755)
				if err != nil {
					log.Fatal(err)
				}

				runNfrCheck(*client, testRun.ID, fileName, nfrData)
			} else {
				result := fetchTestRun(*client, testRun.ID)
				fmt.Println(string(result))
			}
		} else {
			fmt.Println(response)
		}

		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "Could not launch test run!")
		fmt.Println(response)

		os.Exit(1)
	}
}
