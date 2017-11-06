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
