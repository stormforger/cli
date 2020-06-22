package stringutil_test

import (
	"testing"

	"github.com/stormforger/cli/internal/stringutil"
	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, stringutil.InSlice("hi", []string{"hello", "hi", "hallo"}))
	assert.False(t, stringutil.InSlice("huhu", []string{"hello", "hi", "hallo"}))
}
