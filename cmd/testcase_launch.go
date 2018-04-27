package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/jsonapi"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
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
	}
)

func init() {
	TestCaseCmd.AddCommand(testRunLaunchCmd)

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Title, "title", "t", "", "Descriptive title of test run")
	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Notes, "notes", "n", "", "Longer description (Markdown supported)")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.Watch, "watch", "w", false, "Automatically watch newly launched test run")
	testRunLaunchCmd.Flags().DurationVar(&testRunLaunchOpts.MaxWatchTime, "watch-timeout", 0, "Maximum duration in seconds to watch")
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := lookupTestCase(*client, args[0])

	status, response, err := client.TestRunCreate(testCaseUID, testRunLaunchOpts.Title, testRunLaunchOpts.Notes)
	if err != nil {
		log.Fatal(err)
	}

	if status {
		fmt.Println(response)

		if testRunLaunchOpts.Watch {
			testRun := new(api.TestRun)
			err = jsonapi.UnmarshalPayload(strings.NewReader(response), testRun)
			if err != nil {
				log.Fatal(err)
			}

			watchTestRun(testRun.ID, testRunLaunchOpts.MaxWatchTime.Round(time.Second).Seconds())
		}

		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "Could not launch test run!")
		fmt.Println(response)

		os.Exit(1)
	}
}
