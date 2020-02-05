package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, stringInSlice("hi", []string{"hello", "hi", "hallo"}))
	assert.False(t, stringInSlice("huhu", []string{"hello", "hi", "hallo"}))
}
