package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Meta holds meta data of a JSONApi response. Currently
// only "links" are extracted.
type Meta struct {
	Links *Links `json:"links"`
}

// Links holds JSONAPI links
type Links struct {
	Self     string `json:"self"`
	SelfWeb  string `json:"self_web"`
	TestCase string `json:"test_case"`
}

// ErrorPayload holds the list of returned JSONAPI errors
type ErrorPayload struct {
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors"`
}

// ErrorDetail holds data on a specific JSONAPI error
type ErrorDetail struct {
	Code                string          `json:"code"`
	Title               string          `json:"title"`
	Detail              string          `json:"detail"`
	MetaRaw             json.RawMessage `json:"meta"`
	FormattedError      string
	EvaluationErrorMeta *EvaluationErrorMeta
}

// EvaluationErrorMeta holds meta data on JS Evaluation errors
type EvaluationErrorMeta struct {
	Message  string                 `json:"message"`
	RawStack string                 `json:"raw_stack"`
	Name     string                 `json:"name"`
	Stack    []EvaluationStackFrame `json:"stack"`
}

// EvaluationStackFrame represents a stack frame returned by
// evaluation errors
type EvaluationStackFrame struct {
	Context   string `json:"context"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Eval      bool   `json:"eval"`
	Anonymous bool   `json:"anonymous"`
	Internal  bool   `json:"internal"`
}

// UnmarshalMeta will take a io.Reader and try to parse
// "meta" information from a JSONApi response.
func UnmarshalMeta(input io.Reader) (Meta, error) {
	var data struct {
		Meta *Meta `json:"data"`
	}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return Meta{}, err
	}

	return *data.Meta, nil
}

// UnmarshalErrorMeta will take the response (io.Reader) and
// will extract JSONAPI errors.
func UnmarshalErrorMeta(input io.Reader) (ErrorPayload, error) {
	var data ErrorPayload
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return ErrorPayload{}, err
	}

	for i, e := range data.Errors {
		switch e.Code {
		case "E0":
			data.Errors[i].FormattedError = e.Detail
		case "E23":
			errorMeta := new(EvaluationErrorMeta)
			err := json.Unmarshal(e.MetaRaw, errorMeta)
			if err != nil {
				log.Fatal(err)
			}

			data.Errors[i].FormattedError = errorMeta.String()
			data.Errors[i].EvaluationErrorMeta = errorMeta
		}
	}

	return data, nil
}

func (e EvaluationErrorMeta) String() string {
	backtrace := fmt.Sprintf("%s: %s\n", e.Name, e.Message)
	for _, frame := range e.Stack {
		location := ""
		if frame.Anonymous == false {
			location = fmt.Sprintf("%s:%d:%d", frame.File, frame.Line, frame.Column)
		} else {
			location = "<anonymous>"
		}

		if frame.Context != "" {
			backtrace += fmt.Sprintf("    at %s (%s)\n", frame.Context, location)
		} else {
			backtrace += fmt.Sprintf("    at %s\n", location)
		}
	}

	return backtrace
}
