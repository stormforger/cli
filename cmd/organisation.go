package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	organisationCmd = &cobra.Command{
		Use:     "organisation",
		Aliases: []string{"o", "orga", "organization"},
		Short:   "Manage organisations",
		Run: func(cmd *cobra.Command, args []string) {
			log.Fatal("Cannot be run without subcommand, like list!")
		},
	}
)

func init() {
	RootCmd.AddCommand(organisationCmd)
}
