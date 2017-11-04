package testcase

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
