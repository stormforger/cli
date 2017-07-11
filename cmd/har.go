package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	harCmd = &cobra.Command{
		Use:   "har",
		Short: "Convert HAR to test case",
		Long:  `Will convert a given HAR archive into a StormForger test case definition.`,
		Run:   runHar,
	}

	harOpts struct {
		SkipAssets bool
	}
)

func runHar(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		fileName, harFile, err := readFromStdinOrReadFirstArgument(args, "har.json")
		if err != nil {
			log.Fatal(err)
		}

		client := NewClient()

		result, err := client.Har(fileName, harFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	} else {
		log.Fatal("Missing argument; HAR file or - to read from stdin")
	}
}

func init() {
	RootCmd.AddCommand(harCmd)

	harCmd.Flags().BoolVarP(&harOpts.SkipAssets, "skip-assets", "s", false, "Ignore assets?")
}
