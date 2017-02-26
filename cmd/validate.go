package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"github.com/stormforger/cli/datasource"
)

var (
	validateCmd = &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate source files for data source usage",
		Long: `validate will check a given file for its compliance.

		File fixtures may not contain the following in fields:
		* newlines (NL, \n), carriage return (CR, \r)
		* \0 or \1`,
		Example: "forge validate --input users.csv",
		Run:     RunValidate,
	}

	validateOpts struct {
		InputFile    string
		ColSeparator string
		MaxErrors    int
	}
)

type CsvError struct {
	Row     int
	Column  int
	Message string
}

func RunValidate(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal(errors.New("The file argument is required!"))
	}

	validateOpts.InputFile = args[0]

	// FIXME where should be put the validation/check
	//       of the input string and conversion to rune
	columnSeparator, _ := utf8.DecodeRuneInString(validateOpts.ColSeparator[0:])

	validatorConfig := datasource.Validator{
		ColSeparator: columnSeparator,
		MaxErrors:    validateOpts.MaxErrors,
	}

	fd, err := os.Open(validateOpts.InputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	csvErrors, error := datasource.ValidateCSV(fd, validatorConfig)
	if error != nil {
		log.Fatal(error)
	}

	if len(csvErrors) > 0 {
		messageAddition := ""
		if len(csvErrors) >= validateOpts.MaxErrors {
			messageAddition = fmt.Sprintf(" (only showing the first %v)", validateOpts.MaxErrors)
		}
		fmt.Println(fmt.Sprintf("Found %v errors%v:", len(csvErrors), messageAddition))
		for _, error := range csvErrors {
			fmt.Println(fmt.Sprintf("line %v, column %v: "+error.Message, error.Row, error.Column))
		}
	}
}

func init() {
	DatasourceCmd.AddCommand(validateCmd)

	validateCmd.PersistentFlags().StringVar(&validateOpts.ColSeparator, "separator", ",", "Column separator")
	validateCmd.PersistentFlags().IntVar(&validateOpts.MaxErrors, "max-errors", 10, "Stop when encountering more errors")
}
