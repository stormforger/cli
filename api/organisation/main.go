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

// Unmarshal unmarshals a list of Organisations records
func Unmarshal(input io.Reader) (List, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(Organisation)))
	if err != nil {
		return List{}, err
	}

	result := List{}

	for _, item := range items {
		fixture, ok := item.(*Organisation)
		if !ok {
			return List{}, fmt.Errorf("Type assertion failed")
		}

		result.Organisations = append(result.Organisations, fixture)
	}

	return result, nil
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

// FindByNameOrUID look up a Organisation by name in List
func (organisations List) FindByNameOrUID(nameOrUID string) Organisation {
	// first, try to find test by UID
	for _, organisation := range organisations.Organisations {
		if organisation.ID == nameOrUID {
			return *organisation
		}
	}

	// then, try to find by name
	for _, organisation := range organisations.Organisations {
		if organisation.Name == nameOrUID {
			return *organisation
		}
	}

	return Organisation{}
}
