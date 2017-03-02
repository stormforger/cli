package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/buildinfo"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show forge version",
		Run: func(cmd *cobra.Command, args []string) {
			if versionOpts.Verbose == true {
				fmt.Println(buildinfo.BuildInfo)
			} else {
				fmt.Println(buildinfo.BuildInfo.ShortString())
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
