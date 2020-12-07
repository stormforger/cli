package cmd

import (
	"strings"
	"testing"

	"github.com/stormforger/cli/api"
	"github.com/stretchr/testify/assert"
)

func TestPrintValidationResultHuman_NoErrors(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, true, api.ErrorPayload{
		Message: "TestCase updated.",
	})
	assert.Equal(t, buf.String(), "INFO: TestCase updated.\n")
}

func TestPrintValidationResultHuman_ValidationErrors(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, true, api.ErrorPayload{
		Message: "TestCase updated, but validation errors occured.",
		Errors: []api.ErrorDetail{
			{Code: "E0", Title: "Validation Error"},
		},
	})
	assert.Equal(t, buf.String(), "WARN: TestCase updated, but validation errors occured.\n\n1) E0: Validation Error\n")
}

func TestPrintValidationResultHuman_Error(t *testing.T) {
	var buf strings.Builder
	printValidationResultHuman(&buf, false, api.ErrorPayload{
		Message: "TestCase update failed.",
		Errors: []api.ErrorDetail{
			{Code: "E0", Title: "Error"},
		},
	})
	assert.Equal(t, buf.String(), "ERROR: TestCase update failed.\n\n1) E0: Error\n")
}
