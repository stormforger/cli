package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
		var harFile string

		// FIXME this is the same as in testcase_validate.go. Can we extract and generalize this?
		if args[0] == "-" {
			harInput := readFromStdin()
			tmpFile, err := ioutil.TempFile(os.TempDir(), "forge-har")
			if err != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			harFile = tmpFile.Name()

			// TODO what/why this compact syntax with ;?
			if _, err := tmpFile.Write([]byte(harInput)); err != nil {
				log.Fatal(err)
			}

			if err := tmpFile.Close(); err != nil {
				log.Fatal(err)
			}

		} else {
			// FIXME check if file exists here?
			harFile = args[0]
		}

		client := NewClient()

		result, err := client.Har(harFile)
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
