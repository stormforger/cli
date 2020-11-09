package esbundle

import (
	"fmt"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/go-sourcemap/sourcemap"
)

// Result of the bundle process.
type Result struct {
	CompiledContent string
	SourceMapper    SourceMapper
}

// SourceMapper can be used to find the origin of any line/column of the compiled content.
type SourceMapper func(genLine, genColumn int) (source, name string, line, column int, ok bool)

// Bundle returns either a Result or a list of Errors and an error.
func Bundle(inputFile string, replacements map[string]string) (Result, error) {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{inputFile},
		Bundle:            true,
		Write:             false,
		MinifyWhitespace:  false,
		MinifyIdentifiers: false,
		MinifySyntax:      false,
		LogLevel:          esbuild.LogLevelInfo,
		Platform:          esbuild.PlatformNode,
		Defines:           replacements,
		Sourcemap:         esbuild.SourceMapExternal,
		Outdir:            ".",
		// Loaders:           map[string]esbuild.Loader{".jsm": esbuild.LoaderJS},
		// ResolveExtensions: []string{".jsm"},
		// Externals:         []string{"stormforger"},
	})

	if len(result.Errors) > 0 {
		// NOTE: esbuild prints the errors to the console by itself.
		return Result{}, fmt.Errorf("esbuild failed")
	}

	var bundle Result

	// result.OutputFiles should have generated two files: the compiles source (.js) and the source map (.map)
	// map to result.{CompiledContent,SourceMapper}
	for _, file := range result.OutputFiles {
		if strings.HasSuffix(file.Path, ".map") {
			smap, err := sourcemap.Parse("", file.Contents)
			if err != nil {
				panic(err)
			}
			bundle.SourceMapper = smap.Source

		} else if strings.HasSuffix(file.Path, ".js") {
			bundle.CompiledContent = string(file.Contents)
		} else {
			panic("unknown file: " + file.Path)
		}
	}
	return bundle, nil
}
