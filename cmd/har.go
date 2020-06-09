package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	harCmd = &cobra.Command{
		Use:   "har FILE",
		Short: "Convert HAR to test case",
		Long:  `Will convert a given HAR archive into a StormForger test case definition.`,
		Run:   runHar,
	}
)

func runHar(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Missing argument: HAR file or - to read from stdin")
	}

	fileName, harFile, err := readFromStdinOrReadFromArgument(args[0], "stdin")
	if err != nil {
		log.Fatal(err)
	}

	client := NewClient()

	result, err := client.Har(fileName, harFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func init() {
	RootCmd.AddCommand(harCmd)
}
