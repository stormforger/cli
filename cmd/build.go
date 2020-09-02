package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Ping the StormForger API",
		Long:  `Ping the StormForger API and try to authenticate.`,
		Run:   runBuildCmd,
	}

	buildOpts struct {
		Replacements []string
	}
)

func init() {
	RootCmd.AddCommand(buildCmd)

	buildCmd.PersistentFlags().StringSliceVar(&buildOpts.Replacements, "define", []string{}, "Substitute a list of K=V while parsing: env=production,debug=false")
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

	tmpfile, err := ioutil.TempFile("", "forge-js-build")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{args[0]},
		Outfile:     tmpfile.Name(),
		Bundle:      true,
		Write:       true,
		LogLevel:    esbuild.LogLevelInfo,
		Platform:    esbuild.PlatformNode,
		Sourcemap:   esbuild.SourceMapInline,
		Defines:     defines,
		Externals:   []string{"stormforger"},
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	fmt.Print(string(result.OutputFiles[0].Contents))
}
