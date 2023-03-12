package core

import (
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

type Schema struct {
	Content *string

	compiled *jsonschema.Schema
}
type SchemaList map[string]*Schema

func (s *SchemaList) FindCompiled(url string) *jsonschema.Schema {
	schema, err := s.Find(url)
	if err != nil {
		panic(err)
	}

	if schema.compiled == nil {
		compiled, err2 := s.Compile(url)
		if err2 != nil {
			panic(err2)
		}

		schema.compiled = compiled
	}

	return schema.compiled
}

func (s *SchemaList) FindContent(url string) (*string, error) {
	schema, err := s.Find(url)
	if err != nil {
		return nil, err
	}

	if schema.Content == nil {
		return nil, fmt.Errorf("%q has an empty schema", url)
	}

	return schema.Content, nil
}

func (s *SchemaList) Find(url string) (*Schema, error) {
	schema, ok := (*s)[url]
	if !ok {
		return nil, fmt.Errorf("%q not found", url)
	}

	if schema == nil {
		return nil, fmt.Errorf("%q is nil", url)
	}

	return schema, nil
}

func (s *SchemaList) Compile(url string) (*jsonschema.Schema, error) {
	val, err := jsonschema.Compile(url)
	if err != nil {
		return nil, fmt.Errorf("error during %q compilation: %w", url, err)
	}

	return val, nil
}
