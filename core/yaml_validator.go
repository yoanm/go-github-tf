package core

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	YamlAnchorDirectory *string

	//go:embed schemas/repo.json
	repositoryConfigSchema string

	//go:embed schemas/repos.json
	repositoriesConfigSchema string

	//go:embed schemas/branch-protection.json
	branchProtectionSchema string

	//go:embed schemas/branch-branch-protection-template.json
	branchBranchProtectionTemplateSchema string

	//go:embed schemas/branch-branch-protection.json
	branchBranchProtectionSchema string

	//go:embed schemas/branch-protection-template.json
	branchProtectionTemplateSchema string

	//go:embed schemas/branch.json
	branchSchema string

	//go:embed schemas/branch-template.json
	branchTemplateSchema string

	//go:embed schemas/default-branch.json
	defaultBranchSchema string

	//go:embed schemas/repo-template.json
	repositoryTemplateSchema string
)

var Schemas = &SchemaList{
	"map:///repo.json":                              {Content: &repositoryConfigSchema},
	"map:///repos.json":                             {Content: &repositoriesConfigSchema},
	"map:///branch-protection.json":                 {Content: &branchProtectionSchema},
	"map:///branch-protection-template.json":        {Content: &branchProtectionTemplateSchema},
	"map:///branch-branch-protection.json":          {Content: &branchBranchProtectionSchema},
	"map:///branch-branch-protection-template.json": {Content: &branchBranchProtectionTemplateSchema},
	"map:///branch-template.json":                   {Content: &branchTemplateSchema},
	"map:///branch.json":                            {Content: &branchSchema},
	"map:///default-branch.json":                    {Content: &defaultBranchSchema},
	"map:///repo-template.json":                     {Content: &repositoryTemplateSchema},
}

func init() {
	jsonschema.Loaders["map"] = func(url string) (io.ReadCloser, error) {
		schema, err := Schemas.FindContent(url)
		if err != nil {
			return nil, err
		}

		return io.NopCloser(strings.NewReader(*schema)), nil
	}
}

/** Public **/

func ValidateRepositoryConfig(filePath string) (err error) {
	var i interface{}
	if err = loadAsInterface(filePath, &i); err != nil {
		return err
	}

	return _normalizeValidationError(filePath, Schemas.FindCompiled("map:///repo.json").Validate(i))
}

func ValidateRepositoryConfigs(filePath string) (err error) {
	var i interface{}
	if err = loadAsInterface(filePath, &i); err != nil {
		return err
	}

	return _normalizeValidationError(filePath, Schemas.FindCompiled("map:///repos.json").Validate(i))
}

func ValidateRepositoryTemplateConfig(filePath string) (err error) {
	var i interface{}
	if err = loadAsInterface(filePath, &i); err != nil {
		return err
	}

	return _normalizeValidationError(filePath, Schemas.FindCompiled("map:///repo-template.json").Validate(i))
}

func ValidateBranchTemplateConfig(filePath string) (err error) {
	var i interface{}
	if err = loadAsInterface(filePath, &i); err != nil {
		return err
	}

	return _normalizeValidationError(filePath, Schemas.FindCompiled("map:///branch-template.json").Validate(i))
}

func ValidateBranchProtectionTemplateConfig(filePath string) (err error) {
	var i interface{}
	if err = loadAsInterface(filePath, &i); err != nil {
		return err
	}

	return _normalizeValidationError(filePath, Schemas.FindCompiled("map:///branch-protection-template.json").Validate(i))
}

/** Private **/

func loadAsInterface(filePath string, v *interface{}, decoderOpts ...yaml.DecodeOption) error {
	var content []byte
	var err error
	if content, err = os.ReadFile(filePath); err != nil {
		return err
	}

	if err = newDecoder(content, decoderOpts...).Decode(v); err != nil {
		return fmt.Errorf("file %s: %s", filePath, err)
	}

	return nil
}

// USE ONLY ON THAT FILE - START
// e is expected to be an 'error' from jsonschema lib, which is supposed to be a '*jsonschema.ValidationError'
// (but error type can't be cast to jsonschema.ValidationError type as is)
func _normalizeValidationError(filePath string, e interface{}) error {
	if e == nil {
		return nil
	}
	vErr := e.(*jsonschema.ValidationError)

	log.Trace().Msgf("File %s: original validation error => %s", filePath, e)
	err := _validationErrorLeaf(vErr)

	return fmt.Errorf("file %s: %s %s", filePath, err.InstanceLocation, err.Message)
}

func _validationErrorLeaf(ve *jsonschema.ValidationError) *jsonschema.ValidationError {
	if len(ve.Causes) > 0 {
		return _validationErrorLeaf(ve.Causes[0])
	}

	return ve
}

// USE ONLY ON THAT FILE - END
