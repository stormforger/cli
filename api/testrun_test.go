package api

import (
	"net/http"
	"net/http/httptest"
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

// Test that an api response of 404 returns an error with "test-case or logs not found"
func TestTestRunCallLog_404NotFound(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(func() { srv.Close() })

	const testCaseUID = "01234567abc"
	c := NewClient(srv.URL, "")
	r, err := c.TestRunCallLog(testCaseUID, false)
	assert.Nil(t, r)
	assert.NotNil(t, err)

	assert.Equal(t, "test-run or logs not found", err.Error())
}
