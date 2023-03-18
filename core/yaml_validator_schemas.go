package core

import (
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
		return nil, EmptySchemaError(url)
	}

	return schema.Content, nil
}

func (s *SchemaList) Find(url string) (*Schema, error) {
	schema, ok := (*s)[url]
	if !ok {
		return nil, SchemaNotFoundError(url)
	}

	if schema == nil {
		return nil, SchemaIsNilError(url)
	}

	return schema, nil
}

func (s *SchemaList) Compile(url string) (*jsonschema.Schema, error) {
	val, err := jsonschema.Compile(url)
	if err != nil {
		return nil, SchemaCompilationError(url, err.Error())
	}

	return val, nil
}
