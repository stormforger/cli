package testcase

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/google/jsonapi"
)

// List is a list of TestCases, used for index action
type List struct {
	TestCases []*TestCase
}

// TestCase represents a single TestCase
type TestCase struct {
	ID          string `jsonapi:"primary,test_cases"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
}

// ShowNames displays the name and uid of organisations
func ShowNames(input io.Reader) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(TestCase)))

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		testCase, _ := item.(*TestCase)

		fmt.Printf("* %s (ID: %s)\n", testCase.Name, testCase.ID)
	}
}

// Unmarshal unmarshals a list of TestCase records
func Unmarshal(input io.Reader) (List, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(TestCase)))
	if err != nil {
		return List{}, err
	}

	result := List{}

	for _, item := range items {
		testcas, ok := item.(*TestCase)
		if !ok {
			return List{}, fmt.Errorf("Type assertion failed")
		}

		result.TestCases = append(result.TestCases, testcas)
	}

	return result, nil
}

// FindByNameOrUID look up a TestCase by name in List
func (testcases List) FindByNameOrUID(nameOrUID string) TestCase {
	// first, try to find test case by UID
	for _, testCase := range testcases.TestCases {
		if testCase.ID == nameOrUID {
			return *testCase
		}
	}

	// then, try to find case by name
	for _, testCase := range testcases.TestCases {
		if testCase.Name == nameOrUID {
			return *testCase
		}
	}

	return TestCase{}
}
