package esbundle

import (
	"fmt"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func Bundle(inputFile string, replacements map[string]string) (string, error) {
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
		// Loaders:           map[string]esbuild.Loader{".jsm": esbuild.LoaderJS},
		// ResolveExtensions: []string{".jsm"},
		// Externals:         []string{"stormforger"},
	})

	if len(result.Errors) == 0 {
		return string(result.OutputFiles[0].Contents), nil
	}

	return "", fmt.Errorf("Error")
}
