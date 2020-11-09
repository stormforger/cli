package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/internal/esbundle"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build a test case",
		Long:  `Build a test case`,
		Run:   runBuildCmd,
	}

	buildOpts struct {
		Replacements []string
	}
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.PersistentFlags().StringArrayVar(&buildOpts.Replacements, "define", []string{}, "Substitute a list of K=V while parsing: debug=false")
}

func runBuildCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Missing argument: Entry file")
	}

	defines := make(map[string]string)
	for _, kv := range buildOpts.Replacements {
		equals := strings.IndexByte(kv, '=')
		if equals == -1 {
			log.Fatalf("Missing \"=\": %q", kv)
		}

		defines[kv[:equals]] = kv[equals+1:]
	}

	res, err := esbundle.Bundle(args[0], defines)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(res.CompiledContent)
}
