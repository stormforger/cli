package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
	"github.com/stormforger/cli/internal/stringutil"
)

var (
	// testRunWatchCmd represents the test run watch command
	testRunWatchCmd = &cobra.Command{
		Use:   "watch <test-run-id>",
		Short: "Wait and watch for a active test run",
		Long: `Wait and watch for a active test run

watch will continue to look for the active test run until it reaches
a final state (like "done" or "aborted").

It will exit with 0 on success; 1 on test run errors (like "aborted")
and 2 if the given timeout was exceeded.`,
		Run: testRunWatch,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				log.Fatal("Expect exactly one argument: test run reference!")
			}
		},
	}

	testRunWatchOpts struct {
		MaxWatchTime time.Duration
	}
)

var successStates = []string{
	"analysing",
	"deploying",
	"done",
	"fetching_logs",
	"finished",
	"launching",
	"running",
	"starting",
}

func init() {
	TestRunCmd.AddCommand(testRunWatchCmd)

	testRunWatchCmd.Flags().DurationVar(&testRunWatchOpts.MaxWatchTime, "timeout", 0, "Maximum duration in seconds to watch")
}

func testRunWatch(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	watchTestRun(testRunUID, testRunWatchOpts.MaxWatchTime.Round(time.Second).Seconds(), rootOpts.OutputFormat)

	result := fetchTestRun(*client, testRunUID)

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
	}
}

func testRunOkay(testRun *testrun.TestRun) bool {
	return stringutil.InSlice(testRun.State, successStates)
}

func testRunSuccess(testRun *testrun.TestRun) bool {
	successStates := []string{
		"done",
	}

	return stringutil.InSlice(testRun.State, successStates)
}
