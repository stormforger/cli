package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"

	yaml "gopkg.in/yaml.v2"
)

const maxDataSourceSize = 100 * 1024 * 1024 // 100MB

type manifest struct {
	DataSources *[]dataSourceEntry `yaml:"data_sources"`
	TestCases   *[]testCaseEntry   `yaml:"test_cases"`
}

func (m *manifest) lookupDataSource(path string) (bool, dataSourceEntry) {
	for _, dataSource := range *m.DataSources {
		if dataSource.Path == path {
			return true, dataSource
		}
	}

	return false, dataSourceEntry{}
}

func (m *manifest) allDSPaths() (out []string) {
	for _, source := range *m.DataSources {
		out = append(out, source.Path)
	}

	return
}

func (m *manifest) allTCPaths() (out []string) {
	for _, source := range *m.TestCases {
		out = append(out, source.Path)
	}

	return
}

func (m *manifest) validate() (bool, error) {
	names := make(map[string]bool)
	paths := make(map[string]bool)
	for _, ds := range *m.DataSources {
		_, err := ds.validate()
		if err != nil {
			return false, err
		}

		if names[ds.Name] || paths[ds.Path] {
			return false, fmt.Errorf("Manifest: Data sources not unique, %v | %v is defined multiple times", ds.Name, ds.Path)
		}

		names[ds.Name] = true
		paths[ds.Path] = true
	}

	names = make(map[string]bool)
	paths = make(map[string]bool)
	for _, tc := range *m.TestCases {
		_, err := tc.validate()
		if err != nil {
			return false, err
		}

		if names[tc.Name] || paths[tc.Path] {
			return false, fmt.Errorf("Manifest: Test Cases not unique, %v | %v is defined multiple times", tc.Name, tc.Path)
		}

		names[tc.Name] = true
		paths[tc.Path] = true

	}

	return true, nil
}

func (m *manifest) LookupTestCase(path string) (bool, testCaseEntry) {
	for _, testCase := range *m.TestCases {
		if testCase.Path == path {
			return true, testCase
		}
	}

	return false, testCaseEntry{}
}

type dataSourceEntry struct {
	Path            string   `yaml:"path"`
	Organisation    string   `yaml:"organisation"`
	Name            string   `yaml:"name"`
	Delimiter       string   `yaml:"delimiter"`
	AutoColumnNames bool     `yaml:"auto_column_names"`
	Fields          []string `yaml:"fields"`
	Raw             bool     `yaml:"raw"`
}

func (ds *dataSourceEntry) validate() (bool, error) {
	fi, err := os.Stat(ds.Path)
	if err != nil || !fi.Mode().IsRegular() {
		return false, fmt.Errorf("Manifest: data source file not found: %s", ds.Path)
	}

	if fi.Size() > maxDataSourceSize {
		return false, fmt.Errorf("Manifest: %s is larger then 100MB", ds.Path)
	}

	if ds.Raw && (ds.Delimiter != "" || ds.AutoColumnNames || len(ds.Fields) > 0) {
		return false, fmt.Errorf("Manifest: Raw file fixtures do not support fields, delimiter and auto-field-names")
	}

	if ds.AutoColumnNames && len(ds.Fields) > 0 {
		return false, fmt.Errorf("Manifest: fields and auto-field-names are mutually exclusive")
	}

	if ds.Delimiter != "" && utf8.RuneCountInString(ds.Delimiter) > 1 {
		return false, fmt.Errorf("Manifest: Delimiter can only be one character (%s)", ds.Path)
	}

	return false, nil
}

type testCaseEntry struct {
	Path         string `yaml:"path"`
	Organisation string `yaml:"organisation"`
	Name         string `yaml:"name"`
	Comments     string `yaml:"notes"`
}

func (tc *testCaseEntry) validate() (bool, error) {
	fi, err := os.Stat(tc.Path)
	if err != nil || !fi.Mode().IsRegular() {
		return false, fmt.Errorf("Manifest: test case file not found: %s", tc.Path)
	}

	if tc.Organisation == "" {
		return false, fmt.Errorf("Manifest: Organisation missing: %s", tc.Path)
	}

	if tc.Name == "" {
		return false, fmt.Errorf("Manifest: Test Case name missing: %s", tc.Path)
	}

	return false, nil
}

func loadManifest(data io.Reader) (manifest, error) {
	manifestDefinition, err := ioutil.ReadAll(data)
	if err != nil {
		return manifest{}, err
	}

	m := manifest{}
	err = yaml.Unmarshal([]byte(manifestDefinition), &m)
	if err != nil {
		return manifest{}, err
	}

	return m, nil
}
