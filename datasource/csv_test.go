package datasource

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	invalidCsv string
	validCsv   string
)

func defaultValidatorOptions() Validator {
	return Validator{
		ColSeparator: ';',
		MaxErrors:    10,
	}
}

func TestMain(m *testing.M) {
	invalidCsvBytes, err := ioutil.ReadFile("../testdata/csv/invalid.csv")
	if err != nil {
		log.Fatal(err)
	}
	invalidCsv = string(invalidCsvBytes)

	validCsvBytes, err := ioutil.ReadFile("../testdata/csv/valid.csv")
	if err != nil {
		log.Fatal(err)
	}
	validCsv = string(validCsvBytes)

	os.Exit(m.Run())
}

func TestValidateCSVPositive(t *testing.T) {
	r := strings.NewReader(validCsv)

	results, err := ValidateCSV(r, defaultValidatorOptions())

	if err != nil {
		t.Error(err)
	}

	if len(results) > 0 {
		t.Error("Expected no errors, got ", results)
	}
}

func TestValidateCSVNegative(t *testing.T) {
	r := strings.NewReader(invalidCsv)

	results, err := ValidateCSV(r, defaultValidatorOptions())

	if err != nil {
		t.Error(err)
	}

	expected := []CsvError{
		{Row: 2, Column: 2, Message: "Fields may not contain newlines!"},
		{Row: 3, Column: 2, Message: "Fields may not contain the quoting character (\")!"},
		{Row: 17, Column: 1, Message: "Fields may not contain the column separator (;)!"},
	}

	if len(results) != len(expected) {
		t.Error("Expected ", len(expected), " errors, got ", len(results), results)
		return
	}

	for index := range results {
		if actual, expected := results[index], expected[index]; actual != expected {
			t.Error("Expected ", expected, " got ", actual)
		}
	}
}
