package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTestRunResources(t *testing.T) {
	// UID
	assert.Equal(t, TestRunResources{UID: "1a2b3c"}, ExtractTestRunResources("1a2b3c"))

	// Org etc
	assert.Equal(t,
		TestRunResources{Organisation: "zeisss", TestCase: "simple", SequenceID: "10"},
		ExtractTestRunResources("zeisss/simple/test_runs/10"),
	)
	assert.Equal(t,
		TestRunResources{Organisation: "zeisss", TestCase: "simple", SequenceID: "10"},
		ExtractTestRunResources("zeisss/simple/10"),
	)
}
