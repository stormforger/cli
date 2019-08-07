package filefixture

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/google/jsonapi"
)

// List is a list of FileFixtures, used for index action
type List struct {
	Fixtures []*FileFixture
}

// FileFixture represents a single file fixture
type FileFixture struct {
	ID             string     `jsonapi:"primary,file_fixture_structureds"`
	Name           string     `jsonapi:"attr,name"`
	CurrentVersion *Version   `jsonapi:"relation,current_version,omitempty"`
	Versions       []*Version `jsonapi:"relation,versions,omitempty"`
	CreatedAt      time.Time  `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt      time.Time  `jsonapi:"attr,updated_at,iso8601"`
}

// Version respresents a concrete version of a file fixture
type Version struct {
	ID         string    `jsonapi:"primary,file_fixture_version_structureds"`
	Hash       string    `jsonapi:"attr,hash"`
	FileSize   int       `jsonapi:"attr,file_size"`
	ItemCount  int       `jsonapi:"attr,item_count"`
	FieldNames []string  `jsonapi:"attr,field_names"`
	CreatedAt  time.Time `jsonapi:"attr,created_at,iso8601"`
	UpdatedAt  time.Time `jsonapi:"attr,updated_at,iso8601"`
}

// Unmarshal unmarshals a list of FileFixture records
func Unmarshal(input io.Reader) (List, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(FileFixture)))
	if err != nil {
		return List{}, err
	}

	result := List{}

	for _, item := range items {
		fixture, ok := item.(*FileFixture)
		if !ok {
			return List{}, fmt.Errorf("Type assertion failed")
		}

		result.Fixtures = append(result.Fixtures, fixture)
	}

	return result, nil
}

// UnmarshalFileFixture unmarshals a single FileFixture record
func UnmarshalFileFixture(input io.Reader) (*FileFixture, error) {
	fixture := new(FileFixture)
	err := jsonapi.UnmarshalPayload(input, fixture)
	if err != nil {
		return nil, err
	}

	return fixture, nil
}

// FindByName look up a FileFixture by name in List
func (fileFixtures List) FindByName(name string) FileFixture {
	for _, fileFixture := range fileFixtures.Fixtures {
		if fileFixture.Name == name {
			return *fileFixture
		}
	}

	return FileFixture{}
}

// ShowName displays the name, id and hash of a filefixture
func ShowName(input io.Reader) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(FileFixture)))

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fixture, _ := item.(*FileFixture)

		fmt.Printf("* %s (ID: %s, Content-Hash: %s, Fields: %s)\n", fixture.Name, fixture.ID, fixture.CurrentVersion.Hash, fixture.CurrentVersion.FieldNames)
	}
}

// ShowDetails print out details of a file fixture, including its current version
func ShowDetails(fileFixture *FileFixture) {
	fmt.Printf("Name:            %s\n", fileFixture.Name)
	fmt.Printf("UID:             %s\n", fileFixture.ID)
	fmt.Printf("Created:         %s (%s)\n", convertToLocalTZ(fileFixture.CreatedAt), humanize.Time(fileFixture.CreatedAt))
	fmt.Printf("Updated:         %s (%s)\n", convertToLocalTZ(fileFixture.UpdatedAt), humanize.Time(fileFixture.UpdatedAt))
	fmt.Printf("Current Version: %s\n", fileFixture.CurrentVersion.ID)
	fmt.Printf("  SHA256 Hash:   %s\n", fileFixture.CurrentVersion.Hash)
	fmt.Printf("  Size:          %s\n", humanize.Bytes(uint64(fileFixture.CurrentVersion.FileSize)))
	fmt.Printf("  Line Count:    %v\n", fileFixture.CurrentVersion.ItemCount)
	fmt.Printf("  Created:       %s (%s)\n", convertToLocalTZ(fileFixture.CurrentVersion.CreatedAt), humanize.Time(fileFixture.CurrentVersion.CreatedAt))
	fmt.Printf("Version(s):      %v\n", len(fileFixture.Versions))
	for _, version := range fileFixture.Versions {
		fmt.Printf("  - %s (created %s)\n", version.ID, convertToLocalTZ(version.CreatedAt))
	}
}

func convertToLocalTZ(timeToConvert time.Time) time.Time {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.Fatal(err)
	}

	return timeToConvert.In(loc)
}
