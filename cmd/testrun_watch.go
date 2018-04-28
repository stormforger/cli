package cmd

import (
	"bytes"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
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
	}

	testRunWatchOpts struct {
		MaxWatchTime time.Duration
	}
)

var successStates = []string{
	"launching",
	"deploying",
	"starting",
	"running",
	"fetching_logs",
	"log_fetched",
	"analysing",
	"done",
}

func init() {
	TestRunCmd.AddCommand(testRunWatchCmd)

	testRunWatchCmd.Flags().DurationVar(&testRunWatchOpts.MaxWatchTime, "timeout", 0, "Maximum duration in seconds to watch")
}

func testRunWatch(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	result := fetchTestRun(*client, testRunUID)
	testRun, err := testrun.UnmarshalSingle(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	watchTestRun(testRun.ID, testRunWatchOpts.MaxWatchTime.Round(time.Second).Seconds())
}

func testRunOkay(testRun *testrun.TestRun) bool {
	return stringInSlice(testRun.State, successStates)
}

func testRunSuccess(testRun *testrun.TestRun) bool {
	successStates := []string{
		"done",
	}

	return stringInSlice(testRun.State, successStates)
}
