package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var (
	datasourceMoveCmd = &cobra.Command{
		Use:     "mv",
		Aliases: []string{"move", "rename"},
		Short:   "Rename a fixture",
		Run:     runDataSourceMove,
	}
)

func init() {
	datasourceCmd.AddCommand(datasourceMoveCmd)
}

func runDataSourceMove(cmd *cobra.Command, args []string) {
	client := NewClient()

	fileUID := args[0]
	newName := args[1]

	result, err := client.MoveFileFixture(datasourceOpts.Organisation, fileUID, newName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
