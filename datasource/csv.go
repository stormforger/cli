package datasource

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

const quotingRune = '"'

type Validator struct {
	ColSeparator rune
	MaxErrors    int
}

type CsvError struct {
	Row     int
	Column  int
	Message string
}

func DefaultValidatorOptions() Validator {
	return Validator{
		ColSeparator: ';',
		MaxErrors:    10,
	}
}

func ValidateCSV(input io.Reader, options Validator) ([]CsvError, error) {
	reader := csv.NewReader(bufio.NewReader(input))
	reader.Comma = options.ColSeparator
	reader.FieldsPerRecord = 0

	currentRow := 1
	errors := make([]CsvError, 0)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors, err
		}

		for index, value := range record {
			column := index + 1

			if csvError := checkField(value, options, currentRow, column); (csvError != CsvError{}) {
				errors = append(errors, csvError)
			}
		}

		if len(errors) >= options.MaxErrors {
			break
		}

		currentRow += 1
	}

	return errors, nil
}

func checkField(value string, options Validator, row int, column int) CsvError {
	var errorMessage string

	if strings.Contains(value, "\n") {
		errorMessage = "Fields may not contain newlines!"
	}

	if strings.Contains(value, "\r") {
		errorMessage = "Fields may not contain carriage return!"
	}

	if strings.ContainsRune(value, options.ColSeparator) {
		errorMessage = fmt.Sprintf("Fields may not contain the column separator (%v)!", string(options.ColSeparator))
	}

	if strings.ContainsRune(value, quotingRune) {
		errorMessage = fmt.Sprintf("Fields may not contain the quoting character (%v)!", string(quotingRune))
	}

	if errorMessage != "" {
		return CsvError{Row: row, Column: column, Message: errorMessage}
	} else {
		return CsvError{}
	}
}
