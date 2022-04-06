package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api/testrun"
	"github.com/stormforger/cli/internal/stringutil"
)

var (
	// testRunListCmd represents the calllog command
	testRunListCmd = &cobra.Command{
		Use:     "list <test-case-ref>",
		Aliases: []string{"ls"},
		Short:   "List of completed test runs",
		Long:    `List of completed test runs.`,
		Run:     testRunList,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				log.Fatal("Too many arguments")
			}

			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			segments := strings.Split(args[0], "/")

			if len(segments) > 2 {
				log.Fatal("Invalid argument: <test-case-ref> has to be like organisation-name/test-case-name or test-case-uid")
			}
		},
	}

	testRunListOpts struct {
		Archived bool
	}
)

func init() {
	TestRunCmd.AddCommand(testRunListCmd)

	testRunListCmd.Flags().BoolVarP(&testRunListOpts.Archived, "archived", "", false, "List archived test runs")
}

func testRunList(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := mustLookupTestCase(client, args[0])

	filter := ""
	if testRunListOpts.Archived {
		filter = "archived"
	}

	status, result, err := client.TestRunList(testCaseUID, filter)
	if err != nil {
		log.Fatal(err)
	}

	if !status {
		fmt.Fprintln(os.Stderr, "Could not list test runs!")
		fmt.Fprintln(os.Stderr, string(result))

		os.Exit(1)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(result))
		return
	}

	items, err := testrun.Unmarshal(bytes.NewReader(result))
	if err != nil {
		log.Fatal(err)
	}

	if rootOpts.OutputFormat == "human" {
		fmt.Println("ID        Started At           Title")
	}
	for _, item := range items.TestRuns {
		if rootOpts.OutputFormat == "human" {
			fmt.Printf("%s  %s %s\n",
				item.ID,
				stringutil.Coalesce(item.StartedAt, "<no startet_at date>"),
				stringutil.Coalesce(item.Title, "<no title>"),
			)
		} else { // plain
			fmt.Printf("%s\n", item.ID)
		}
	}
}
