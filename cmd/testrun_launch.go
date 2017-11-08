package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	// testRunLaunchCmd represents the test run launch command
	testRunLaunchCmd = &cobra.Command{
		Use:   "launch <test-case-ref>",
		Short: "Create and launch a new test run",
		Long:  `Create and launch a new test run based on given test case`,
		Run:   testRunLaunch,
	}
)

func init() {
	TestRunCmd.AddCommand(testRunLaunchCmd)
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	title := "title"
	notes := "notes"

	status, response, err := client.TestRunCreate(args[0], title, notes)
	if err != nil {
		log.Fatal(err)
	}

	if status {
		fmt.Println(response)

		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "Could not launch test run!")
		fmt.Println(response)

		os.Exit(1)
	}
}
