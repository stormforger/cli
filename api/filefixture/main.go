package filefixture

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/google/jsonapi"
)

// FileFixture TODO
type FileFixture struct {
	ID             string   `jsonapi:"primary,file_fixture_structureds"`
	Name           string   `jsonapi:"attr,name"`
	CurrentVersion *Version `jsonapi:"relation,current_version"`
}

// Version TODO
type Version struct {
	ID               int    `jsonapi:"primary,file_fixture_version_structureds"`
	OriginalMd5Hash  string `jsonapi:"attr,md5_hash"`
	ProcessedMd5Hash string `jsonapi:"attr,processed_md5_hash"`
	FieldNames       string `jsonapi:"attr,field_names"`
}

// ShowName TODO
func ShowName(input io.Reader) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(FileFixture)))

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		fixture, _ := item.(*FileFixture)

		fmt.Printf("* %s (ID: %s, Content-MD5: %s)\n", fixture.Name, fixture.ID, fixture.CurrentVersion.OriginalMd5Hash)
	}
}
