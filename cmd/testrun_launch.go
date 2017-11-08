package cmd

import (
	"fmt"
	"log"

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

	_, response, err := client.TestRunCreate(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}
