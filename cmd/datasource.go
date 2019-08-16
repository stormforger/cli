package cmd

import "github.com/spf13/cobra"

var (
	datasourceCmd = &cobra.Command{
		Use:     "datasource",
		Aliases: []string{"ds"},
		Short:   "Work with and manage data sources",
	}

	datasourceOpts struct {
		Organisation string
	}
)

func init() {
	RootCmd.AddCommand(datasourceCmd)

	datasourceCmd.PersistentFlags().StringVarP(&datasourceOpts.Organisation, "organisation", "o", "", "Name of the organisation")
}
