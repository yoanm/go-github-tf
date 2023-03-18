package core

import (
	"bytes"
	"os"

	"github.com/goccy/go-yaml"
)

func LoadRepositoriesFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error) {
	if err := ValidateRepositoryConfigs(filePath); err != nil {
		return nil, err
	}

	return LoadGhRepoConfigListFromFile(filePath, decoderOpts...)
}

func LoadRepositoryFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error) {
	if err := ValidateRepositoryConfig(filePath); err != nil {
		return nil, err
	}

	return LoadGhRepoConfigFromFile(filePath, decoderOpts...)
}

func LoadRepositoryTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error) {
	if err := ValidateRepositoryTemplateConfig(filePath); err != nil {
		return nil, err
	}

	return LoadGhRepoConfigFromFile(filePath, decoderOpts...)
}

func LoadBranchTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhBranchConfig, error) {
	if err := ValidateBranchTemplateConfig(filePath); err != nil {
		return nil, err
	}

	return LoadGhRepoBranchConfigFromFile(filePath, decoderOpts...)
}

func LoadBranchProtectionTemplateFromFile(
	filePath string,
	decoderOpts ...yaml.DecodeOption,
) (*GhBranchProtectionConfig, error) {
	if err := ValidateBranchProtectionTemplateConfig(filePath); err != nil {
		return nil, err
	}

	return LoadGhRepoBranchProtectionConfigFromFile(filePath, decoderOpts...)
}

// LoadGhRepoConfigFromFile loads the file content to GhRepoConfig struct
// No schema validation will be performed, use loadRepositoryFromFile or loadRepositoryTemplateFromFile instead !
func LoadGhRepoConfigFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error) {
	var (
		content []byte
		err     error
	)

	if content, err = os.ReadFile(filePath); err != nil {
		//nolint:wrapcheck // Expected to return unwrap error
		return nil, err
	}

	//nolint:exhaustruct // No need here, simple init
	config := &GhRepoConfig{}
	if err = newDecoder(content, decoderOpts...).Decode(config); err != nil {
		return nil, FileError(filePath, err)
	}

	return config, nil
}

// LoadGhRepoConfigListFromFile loads the file content to GhRepoConfig struct
// No schema validation will be performed, use loadRepositoriesFromFile instead !
func LoadGhRepoConfigListFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		//nolint:wrapcheck // Expected to return unwrap error
		return nil, err
	}

	var configs []*GhRepoConfig
	if err = newDecoder(content, decoderOpts...).Decode(&configs); err != nil {
		return nil, FileError(filePath, err)
	}

	return configs, nil
}

// LoadGhRepoBranchConfigFromFile loads the file content to GhBranchConfig struct
// No schema validation will be performed, use loadBranchTemplateFromFile instead !
func LoadGhRepoBranchConfigFromFile(
	filePath string,
	decoderOpts ...yaml.DecodeOption,
) (*GhBranchConfig, error) {
	var (
		content []byte
		err     error
	)

	if content, err = os.ReadFile(filePath); err != nil {
		//nolint:wrapcheck // Expected to return unwrap error
		return nil, err
	}

	//nolint:exhaustruct // No need here, simple init
	config := &GhBranchConfig{}
	if err = newDecoder(content, decoderOpts...).Decode(config); err != nil {
		return nil, FileError(filePath, err)
	}

	return config, nil
}

// LoadGhRepoBranchProtectionConfigFromFile loads the file content to GhBranchProtectionConfig struct
// No schema validation will be performed, use loadBranchProtectionTemplateFromFile instead !
func LoadGhRepoBranchProtectionConfigFromFile(
	filePath string,
	decoderOpts ...yaml.DecodeOption,
) (*GhBranchProtectionConfig, error) {
	var (
		content []byte
		err     error
	)

	if content, err = os.ReadFile(filePath); err != nil {
		//nolint:wrapcheck // Expected to return unwrap error
		return nil, err
	}

	//nolint:exhaustruct // No need here, simple init
	config := &GhBranchProtectionConfig{}
	if err = newDecoder(content, decoderOpts...).Decode(config); err != nil {
		return nil, FileError(filePath, err)
	}

	return config, nil
}

/** Private **/

func newDecoder(content []byte, decoderOpts ...yaml.DecodeOption) *yaml.Decoder {
	return yaml.NewDecoder(
		bytes.NewBuffer(content),
		append(getYamlValidatorDecoderOptions(), decoderOpts...)...,
	)
}

func getYamlValidatorDecoderOptions() []yaml.DecodeOption {
	list := []yaml.DecodeOption{yaml.DisallowDuplicateKey()}

	if YamlAnchorDirectory != nil {
		list = append(list, yaml.ReferenceDirs(*YamlAnchorDirectory))
	}

	return list
}
