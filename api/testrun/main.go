package testrun

import (
	"fmt"
	"io"
	"reflect"

	"github.com/google/jsonapi"
)

// List is a list of TestRuns, used for index action
type List struct {
	TestRuns []*TestRun
}

// TestRun represents a single TestRun
type TestRun struct {
	ID        string `jsonapi:"primary,test_runs"`
	Scope     string `jsonapi:"attr,scope"`
	Title     string `jsonapi:"attr,title,omitempty"`
	Notes     string `jsonapi:"attr,notes,omitempty"`
	State     string `jsonapi:"attr,state,omitempty"`
	StartedBy string `jsonapi:"attr,started_by,omitempty"`
	StartedAt string `jsonapi:"attr,started_at,omitempty"`
	EndedAt   string `jsonapi:"attr,ended_at,omitempty"`

	// attributes for in progress
	EstimatedEnd string `jsonapi:"attr,estimated_end,omitempty"`
	Progress     int    `jsonapi:"attr,progress,omitempty"`
}

// NfrResultList is a list of NFR results
type NfrResultList struct {
	NfrResults []*NfrResult
}

// NfrResult describes a NFR check result
type NfrResult struct {
	ID               string `jsonapi:"primary,nfr_results"`
	Success          bool   `jsonapi:"attr,success"`
	Subject          string `jsonapi:"attr,subject"`
	SubjectAvailable bool   `jsonapi:"attr,subject_available"`
	SubjectUnit      string `jsonapi:"attr,subject_unit"`
	Expectation      string `jsonapi:"attr,expectation"`
	Type             string `jsonapi:"attr,nfr_type"`
	Disabled         bool   `jsonapi:"attr,disabled"`
	Filter           string `jsonapi:"attr,filter"`
	Metric           string `jsonapi:"attr,metric"`
}

// Unmarshal unmarshals a list of TestRun records
func Unmarshal(input io.Reader) (List, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(TestRun)))
	if err != nil {
		return List{}, err
	}

	result := List{}

	for _, item := range items {
		testRun, ok := item.(*TestRun)
		if !ok {
			return List{}, fmt.Errorf("Type assertion failed")
		}

		result.TestRuns = append(result.TestRuns, testRun)
	}

	return result, nil
}

// UnmarshalSingle unmarshals a single TestRun records
func UnmarshalSingle(input io.Reader) (TestRun, error) {
	item := new(TestRun)
	err := jsonapi.UnmarshalPayload(input, item)
	if err != nil {
		return TestRun{}, err
	}

	return *item, nil
}

// UnmarshalNfrResults unmarshals a list of NFR result records
func UnmarshalNfrResults(input io.Reader) (NfrResultList, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(NfrResult)))
	if err != nil {
		return NfrResultList{}, err
	}

	result := NfrResultList{}

	for _, item := range items {
		typedItem, ok := item.(*NfrResult)
		if !ok {
			return NfrResultList{}, fmt.Errorf("Type assertion failed")
		}

		result.NfrResults = append(result.NfrResults, typedItem)
	}

	return result, nil
}

// SubjectWithUnit formats the expectation inclusing the subject's unit
func (nfr *NfrResult) SubjectWithUnit() string {
	if nfr.SubjectUnit != "" {
		return nfr.Subject + " " + nfr.SubjectUnit
	}

	return nfr.Subject
}

// ExpectationWithUnit formats the expectation inclusing the subject's unit
func (nfr *NfrResult) ExpectationWithUnit() string {
	if nfr.SubjectUnit != "" {
		return nfr.Expectation + " " + nfr.SubjectUnit
	}

	return nfr.Expectation
}
