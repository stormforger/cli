package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
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
		MaxWatchTime int
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

	testRunWatchCmd.Flags().IntVar(&testRunWatchOpts.MaxWatchTime, "timeout", 0, "Maximum duration in seconds to watch")
}

func testRunWatch(cmd *cobra.Command, args []string) {
	client := NewClient()

	started := time.Now()

	for true {
		runningSince := time.Now().Sub(started).Seconds()

		testRun, response, err := client.TestRunWatch(args[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response)

		if !testRunOkay(&testRun) {
			os.Exit(1)
		}

		if testRunSuccess(&testRun) {
			os.Exit(0)
		}

		if testRunWatchOpts.MaxWatchTime > 0 && int(runningSince) > testRunWatchOpts.MaxWatchTime {
			os.Exit(2)
		}

		time.Sleep(5 * time.Second)
	}
}

func testRunOkay(testRun *api.TestRun) bool {
	return stringInSlice(testRun.State, successStates)
}

func testRunSuccess(testRun *api.TestRun) bool {
	successStates := []string{
		"done",
	}

	return stringInSlice(testRun.State, successStates)
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
