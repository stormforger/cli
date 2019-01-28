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

func (m *manifest) lookupDataSource(path string) dataSourceEntry {
	for _, dataSource := range *m.DataSources {
		if dataSource.Path == path {
			return dataSource
		}
	}

	return dataSourceEntry{}
}

func (m *manifest) allDSPaths() (out []string) {
	for _, source := range *m.DataSources {
		out = append(out, source.Path)
	}

	return
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
		return false, fmt.Errorf("Manifest: file not found: %s", ds.Path)
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

	return true, nil
}

type testCaseEntry struct {
	Path         string `yaml:"path"`
	Organisation string `yaml:"organisation"`
	Name         string `yaml:"name"`
}

func (m *manifest) LookupTestCase(path string) testCaseEntry {
	for _, testCase := range *m.TestCases {
		if testCase.Path == path {
			return testCase
		}
	}

	return testCaseEntry{}
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
