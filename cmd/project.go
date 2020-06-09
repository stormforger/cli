package cmd

import "github.com/spf13/cobra"

var (
	projectCmd = &cobra.Command{
		Use:     "project",
		Aliases: []string{"ds"},
		Short:   "Work with multiple test cases and data sources",
	}

	projectOpts struct {
		manifestFile string
	}
)

func init() {
	RootCmd.AddCommand(projectCmd)

	projectCmd.PersistentFlags().StringVarP(&projectOpts.manifestFile, "config", "p", "", "Path to project configuration")
}
