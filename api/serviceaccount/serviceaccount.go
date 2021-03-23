package serviceaccount

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/google/jsonapi"
)

// List is a list of ServiceAccounts, used for index action
type List struct {
	ServiceAccounts []*ServiceAccount
}

// ServiceAccount represents a single ServiceAccount
type ServiceAccount struct {
	UID                        string     `jsonapi:"primary,service_account"`
	TokenLabel                 string     `jsonapi:"attr,token_label"`
	GeneratedAt                *time.Time `jsonapi:"attr,generated_at,iso8601"`
	MostRecentAPIAccessAt      *time.Time `jsonapi:"attr,most_recent_api_access_at,iso8601"`
	MostRecentAPIClientVersion string     `jsonapi:"attr,most_recent_api_client_version"`

	// AccessToken is returned once when creating a new service account
	AccessToken string `jsonapi:"attr,access_token"`
}

// Unmarshal unmarshals a list of TestCase records
func UnmarshalList(input io.Reader) (List, error) {
	items, err := jsonapi.UnmarshalManyPayload(input, reflect.TypeOf(new(ServiceAccount)))
	if err != nil {
		return List{}, err
	}

	result := List{}

	for _, item := range items {
		serviceAccount, ok := item.(*ServiceAccount)
		if !ok {
			return List{}, fmt.Errorf("Type assertion failed")
		}

		result.ServiceAccounts = append(result.ServiceAccounts, serviceAccount)
	}

	return result, nil
}

// Unmarshal unmarshals a single ServiceAccount record
func Unmarshal(input io.Reader) (*ServiceAccount, error) {
	serviceAccount := new(ServiceAccount)
	err := jsonapi.UnmarshalPayload(input, serviceAccount)
	if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}
