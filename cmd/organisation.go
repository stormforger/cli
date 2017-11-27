package cmd

import "github.com/spf13/cobra"

var (
	organisationCmd = &cobra.Command{
		Use:     "organisation",
		Aliases: []string{"o", "orga", "organization"},
		Short:   "Manage organisations",
	}
)

func init() {
	RootCmd.AddCommand(organisationCmd)
}
