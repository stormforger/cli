package cmd

import (
	"io/ioutil"
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

func GivenTestdataContents(t *testing.T, file string) string {
	p := filepath.Join("testdata", file)

	data, err := ioutil.ReadFile(p)
	require.NoError(t, err)

	return string(data)
}
