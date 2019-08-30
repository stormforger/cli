package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/stormforger/cli/api"
	"github.com/stormforger/cli/api/testrun"
)

var (
	// testRunLaunchCmd represents the test run launch command
	testRunLaunchCmd = &cobra.Command{
		Use:   "launch <test-case-ref>",
		Short: "Create and launch a new test run",
		Long: fmt.Sprintf(`Create and launch a new test run based on given test case

<test-case-ref> can be 'organisation-name/test-case-name' or 'test-case-uid'.

Examples
--------
* Launch by organisation and test case name

  forge test-case launch acme-inc/checkout

* Alternatively the test case UID can also be provided

  forge test-case launch xPSX5KXM


Configuration
-------------
You can specify configuration for a test run that will overwrite what is defined
in your JavaScript definition.

* Available cluster sizings:
  * %s

Available cluster regions are available at https://docs.stormforger.com/reference/test-cluster/#cluster-region
`,
			strings.Join(validSizings, "\n  * ")),
		Run: testRunLaunch,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Missing argument: test case reference")
			}

			if len(args) > 1 {
				log.Fatal("Too many arguments")
			}

			if testRunLaunchOpts.DumpTraffic && testRunLaunchOpts.CheckNFR != "" {
				log.Fatal("--dump-traffic and --nfr-check-file are mutual exclusive")
			}

			if testRunLaunchOpts.ClusterRegion != "" && !stringInSlice(testRunLaunchOpts.ClusterRegion, validRegions) {
				log.Fatalf("%s is not a valid region", testRunLaunchOpts.ClusterRegion)
			}

			if testRunLaunchOpts.ClusterSizing != "" && !stringInSlice(testRunLaunchOpts.ClusterSizing, validSizings) {
				log.Fatalf("%s is not a valid sizing", testRunLaunchOpts.ClusterSizing)
			}
		},
	}

	testRunLaunchOpts struct {
		OpenInBrowser         bool
		Title                 string
		Notes                 string
		ClusterRegion         string
		ClusterSizing         string
		Watch                 bool
		MaxWatchTime          time.Duration
		CheckNFR              string
		DisableGzip           bool
		SkipWait              bool
		DumpTraffic           bool
		SessionValidationMode bool
		Validate              bool
	}

	validRegions = []string{
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-south-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ca-central-1",
		"eu-central-1",
		"eu-north-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}

	validSizings = []string{
		"preflight",
		"tiny",
		"small",
		"medium",
		"large",
		"xlarge",
		"2xlarge",
	}
)

func init() {
	TestCaseCmd.AddCommand(testRunLaunchCmd)

	testRunLaunchCmd.Flags().BoolVar(&testRunLaunchOpts.OpenInBrowser, "open", false, "Open test run in browser")

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Title, "title", "t", "", "Descriptive title of test run")
	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.Notes, "notes", "n", "", "Longer description (Markdown supported)")

	testRunLaunchCmd.Flags().StringVar(&testRunLaunchOpts.ClusterRegion, "region", "", "Region to start test in")
	testRunLaunchCmd.Flags().StringVar(&testRunLaunchOpts.ClusterSizing, "sizing", "", "Cluster sizing to use")

	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.Watch, "watch", "w", false, "Automatically watch newly launched test run")
	testRunLaunchCmd.Flags().DurationVar(&testRunLaunchOpts.MaxWatchTime, "watch-timeout", 0, "Maximum duration in seconds to watch")

	testRunLaunchCmd.Flags().StringVarP(&testRunLaunchOpts.CheckNFR, "nfr-check-file", "", "", "Check test result against NFR definition (implies --watch)")

	// options for debugging
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.DisableGzip, "disable-gzip", "", false, "Globally disable gzip")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.SkipWait, "skip-wait", "", false, "Ignore defined waits")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.DumpTraffic, "dump-traffic", "", false, "Create traffic dump")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.SessionValidationMode, "session-validation-mode", "", false, "Enable session validation mode")
	testRunLaunchCmd.Flags().BoolVarP(&testRunLaunchOpts.Validate, "validate", "", false, "Perform validation run")
}

func testRunLaunch(cmd *cobra.Command, args []string) {
	client := NewClient()

	testCaseUID := lookupTestCase(*client, args[0])

	launchOptions := api.TestRunLaunchOptions{
		Title:                 testRunLaunchOpts.Title,
		Notes:                 testRunLaunchOpts.Notes,
		ClusterRegion:         testRunLaunchOpts.ClusterRegion,
		ClusterSizing:         testRunLaunchOpts.ClusterSizing,
		DisableGzip:           testRunLaunchOpts.DisableGzip,
		SkipWait:              testRunLaunchOpts.SkipWait,
		DumpTraffic:           testRunLaunchOpts.DumpTraffic,
		SessionValidationMode: testRunLaunchOpts.SessionValidationMode,
	}

	if testRunLaunchOpts.Validate {
		launchOptions.SessionValidationMode = true
		launchOptions.ClusterSizing = "preflight"
	}

	status, response, err := client.TestRunCreate(testCaseUID, launchOptions)
	if err != nil {
		log.Fatal(err)
	}

	if !status {
		fmt.Fprintln(os.Stderr, "Could not launch test run!")
		fmt.Println(response)

		os.Exit(1)
	}

	testRun, err := testrun.UnmarshalSingle(strings.NewReader(response))
	if err != nil {
		log.Fatal(err)
	}

	if rootOpts.OutputFormat == "json" {
		fmt.Println(string(response))
	} else {
		// FIXME can we integrate this into testrun.UnmarshalSingle somehow?
		meta, err := api.UnmarshalMeta(strings.NewReader(response))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(`Launching test %s
UID: %s
Web URL: %s
`,
			testRun.Scope,
			testRun.ID,
			meta.Links.SelfWeb,
		)

		fmt.Printf("Configuration: %s cluster in %s\n", testRun.TestConfiguration.ClusterSizing, testRun.TestConfiguration.ClusterRegion)

		if testRun.TestConfiguration.DisableGzip {
			fmt.Print("  [\u2713] Disabled GZIP\n")
		}
		if testRun.TestConfiguration.SkipWait {
			fmt.Print("  [\u2713] Skip Waits\n")
		}
		if testRun.TestConfiguration.DumpTrafficFull {
			fmt.Print("  [\u2713] Traffic Dump\n")
		}
		if testRun.TestConfiguration.SessionValidationMode {
			fmt.Print("  [\u2713] Session Validation Mode\n")
		}

		if testRunLaunchOpts.OpenInBrowser {
			fmt.Printf("Opening %s in browser...\n", meta.Links.SelfWeb)
			err = browser.OpenURL(meta.Links.SelfWeb)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if testRunLaunchOpts.Watch || testRunLaunchOpts.CheckNFR != "" || testRunLaunchOpts.Validate {
		if rootOpts.OutputFormat != "json" {
			fmt.Println("\nWatching...")
		}

		watchTestRun(testRun.ID, testRunLaunchOpts.MaxWatchTime.Round(time.Second).Seconds(), rootOpts.OutputFormat)

		if testRunLaunchOpts.CheckNFR != "" || testRunLaunchOpts.Validate {
			fmt.Println("Test finished, running non-functional checks...")

			fileName := ""
			var nfrData io.Reader
			if testRunLaunchOpts.CheckNFR != "" {
				fileName = filepath.Base(testRunLaunchOpts.CheckNFR)
				nfrData, err = os.OpenFile(testRunLaunchOpts.CheckNFR, os.O_RDONLY, 0755)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fileName = "validation.yml"
				nfrData = bytes.NewBufferString(`version: "0.1"
requirements:
- test.completed: true
- checks:
    select: success_rate
    test: ["=", 1]
- http.error_ratio:
    test: ["=", 0]`)
			}

			runNfrCheck(*client, testRun.ID, fileName, nfrData)
		} else {
			result := fetchTestRun(*client, testRun.ID)
			fmt.Println(string(result))
		}
	}
}
