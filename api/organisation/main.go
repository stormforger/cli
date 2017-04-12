package organisation

import (
	"fmt"
	"io"
	"log"
	"reflect"

	"github.com/google/jsonapi"
)

// List is a list of Organisations, used for index action
type List struct {
	Organisations []*Organisation
}

// Organisation represents a single organisation
type Organisation struct {
	ID   string `jsonapi:"primary,organisations"`
	Name string `jsonapi:"attr,name"`
}

// ShowNames displays the name and uid of organisations
func ShowNames(input io.Reader) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(Organisation)))

	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		organisation, _ := item.(*Organisation)

		fmt.Printf("* %s (ID: %s)\n", organisation.Name, organisation.ID)
	}
}
