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
