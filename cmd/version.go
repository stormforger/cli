package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show forge version",
		Run: func(cmd *cobra.Command, args []string) {
			if versionOpts.Verbose == true {
				fmt.Println(buildInfo)
			} else {
				fmt.Println(buildInfo.ShortString())
			}
		},
	}

	versionOpts struct {
		Verbose bool
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolVar(&versionOpts.Verbose, "verbose", false, "Print build information")
}
