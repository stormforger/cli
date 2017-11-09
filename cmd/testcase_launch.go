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

	testRunLaunchOpts struct {
		Title string
		Notes string
	}
)

func init() {
	TestCaseCmd.AddCommand(testRunLaunchCmd)

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Title, "title", "t", "", "Descriptive title of test run")
	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Notes, "notes", "n", "", "Longer description (Markdown supported)")
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	status, response, err := client.TestRunCreate(args[0], testRunLaunchOpts.Title, testRunLaunchOpts.Notes)
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
