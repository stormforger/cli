package cmd

import (
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
	log.Fatal("NOT IMPLEMENTED")
}
