package cmd

import (
	"strings"
	"testing"

	"github.com/stormforger/cli/api"
	"github.com/stretchr/testify/assert"
)

func TestPrintValidationResultHuman_NoErrors(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, "test.js", true, api.ErrorPayload{
		Message: "TestCase updated.",
	})
	assert.Equal(t, buf.String(), "# FILE: test.js\nINFO: TestCase updated.\n")
}

func TestPrintValidationResultHuman_ValidationErrors(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, "test.js", true, api.ErrorPayload{
		Message: "TestCase updated, but validation errors occured.",
		Errors: []api.ErrorDetail{
			{Code: "E0", Title: "Validation Error"},
		},
	})
	assert.Equal(t, buf.String(), "# FILE: test.js\nWARN: TestCase updated, but validation errors occured.\n\n1) E0: Validation Error\n\n")
}

func TestPrintValidationResultHuman_Error(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, "test.js", false, api.ErrorPayload{
		Message: "TestCase update failed.",
		Errors: []api.ErrorDetail{
			{Code: "E0", Title: "Error"},
		},
	})
	assert.Equal(t, buf.String(), "# FILE: test.js\nERROR: TestCase update failed.\n\n1) E0: Error\n\n")
}
