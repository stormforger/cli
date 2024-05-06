package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestcaseBuild(t *testing.T) {
	var expectedOutput = GivenTestdataContents(t, t.Name()+"_output.js")

	var output strings.Builder
	err := MainTestcaseBuild(&output, filepath.Join("testdata", t.Name()+"_main.mjs"), nil)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, output.String())
}

func TestTestcaseBuild__DefineParam(t *testing.T) {
	var expectedOutput = GivenTestdataContents(t, t.Name()+"_output.js")
	defines := map[string]string{
		"defines.target": "\"http://example.com\"", // DOUBLE QUOTING!
	}

	var output strings.Builder
	err := MainTestcaseBuild(&output, filepath.Join("testdata", t.Name()+"_main.mjs"), defines)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, output.String())
}

// Same input as TestTestcaseBuild__DefineParam but we do NOT specify the define parameter.
// This must still build without an error
func TestTestcaseBuild__WithUndefinedVariable(t *testing.T) {
	var expectedOutput = GivenTestdataContents(t, t.Name()+"_output.js")

	var output strings.Builder
	err := MainTestcaseBuild(&output, filepath.Join("testdata", t.Name()+"_main.mjs"), nil)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, output.String())
}

func GivenTestdataContents(t *testing.T, file string) string {
	p := filepath.Join("testdata", file)

	data, err := os.ReadFile(p)
	require.NoError(t, err)

	return string(data)
}
