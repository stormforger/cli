package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	harCmd = &cobra.Command{
		Use:   "har <file>",
		Short: "Convert HAR to test case",
		Long:  `Will convert a given HAR archive into a StormForger test case definition. Pass - for the file to read from stdin.`,
		Args:  cobra.ExactArgs(1),
		Run:   runHar,
	}
)

func runHar(cmd *cobra.Command, args []string) {
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
