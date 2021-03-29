package cmd

import (
	"io/ioutil"
	"path"
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

func GivenFolder(t *testing.T, files map[string][]byte) string {
	tempPath := t.TempDir()
	for name, content := range files {
		p := path.Join(tempPath, name)
		err := ioutil.WriteFile(p, content, 0444)
		require.NoError(t, err)
	}
	return tempPath
}
