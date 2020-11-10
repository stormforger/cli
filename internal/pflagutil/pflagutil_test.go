package pflagutil_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"github.com/stormforger/cli/internal/pflagutil"
)

func TestKeyValueFlag(t *testing.T) {
	var m map[string]string
	val := &pflagutil.KeyValueFlag{Map: &m}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	fs.Var(val, "define", "define docstring")

	err := fs.Parse([]string{"--define", "a=1", "--define", "b=2", "--define", "name=\"Harry\""})
	assert.NoError(t, err)

	assert.Len(t, m, 3)
	assert.Equal(t, "1", m["a"])
	assert.Equal(t, "2", m["b"])

	// NOTE If you cannot reproduce this behavior on your system your shell eats your quotes >:)
	assert.Equal(t, "\"Harry\"", m["name"])
}
