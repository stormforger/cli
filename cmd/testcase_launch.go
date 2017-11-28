package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/jsonapi"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
)

var (
	// testRunLaunchCmd represents the test run launch command
	testRunLaunchCmd = &cobra.Command{
		Use:   "launch <test-case-ref>",
		Short: "Create and launch a new test run",
		Long:  `Create and launch a new test run based on given test case`,
		Run:   testRunLaunch,
	}

	testRunLaunchOpts struct {
		Title        string
		Notes        string
		Watch        bool
		MaxWatchTime int
	}
)

func init() {
	TestCaseCmd.AddCommand(testRunLaunchCmd)

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Title, "title", "t", "", "Descriptive title of test run")
	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Notes, "notes", "n", "", "Longer description (Markdown supported)")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.Watch, "watch", "w", false, "Automatically watch newly launched test run")
	testRunLaunchCmd.Flags().IntVar(&testRunLaunchOpts.MaxWatchTime, "watch-timeout", 0, "Maximum duration in seconds to watch")
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, response, err := client.TestRunCreate(args[0], testRunLaunchOpts.Title, testRunLaunchOpts.Notes)
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

			watchTestRun(testRun.ID, testRunLaunchOpts.MaxWatchTime)
		}

		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "Could not launch test run!")
		fmt.Println(response)

		os.Exit(1)
	}
}
