package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	trShareURLCmd = &cobra.Command{
		Use:     "share-url <test-run-ref>",
		Aliases: []string{"shareurl"},
		Args:    cobra.ExactArgs(1),
		Short:   "Generate a share url for the testrun",
		Example: "$ forge test-run share-url example-org/my-test-case/1\n" +
			"$ forge test-run share-url a17fac2 --expire-duration 24h",
		Run: shareUrl,
	}

	shareURLOpts struct {
		ExpireDuration time.Duration
	}
)

func init() {
	TestRunCmd.AddCommand(trShareURLCmd)

	trShareURLCmd.Flags().DurationVar(&shareURLOpts.ExpireDuration, "expire-duration", 0, "Expire duration for this token - zero leaves this up to the server")
}

func shareUrl(cmd *cobra.Command, args []string) {
	client := NewClient()

	testRunUID := getTestRunUID(*client, args[0])

	shareURL, err := client.TestRunShareURL(cmd.Context(), testRunUID, shareURLOpts.ExpireDuration)
	if err != nil {
		log.Fatal(err)
	}

	switch rootOpts.OutputFormat {
	case "json":
		enc := json.NewEncoder(cmd.OutOrStdout())
		enc.SetIndent("", "  ")
		enc.Encode(shareURL)
	default:
		fmt.Println("URL:\t\t", shareURL.URL)
		if shareURL.ExpiresAt != nil && !shareURL.ExpiresAt.IsZero() {
			fmt.Println("Expires at:\t", shareURL.ExpiresAt.Format(time.RFC3339))
		}
	}

}
