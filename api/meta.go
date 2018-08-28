package api

import (
	"encoding/json"
	"io"
)

// Meta holds meta data of a JSONApi response. Currently
// only "links" are extracted.
type Meta struct {
	Links *Links `json:"links"`
}

// Links holds JSONAPI links
type Links struct {
	Self     string `json:"self"`
	SelfWeb  string `json:"self_web"`
	TestCase string `json:"test_case"`
}

// UnmarshalMeta will take a io.Reader and try to parse
// "meta" information from a JSONApi response.
func UnmarshalMeta(input io.Reader) (Meta, error) {
	var data struct {
		Meta *Meta `json:"data"`
	}
	if err := json.NewDecoder(input).Decode(&data); err != nil {
		return Meta{}, err
	}

	return *data.Meta, nil
}
