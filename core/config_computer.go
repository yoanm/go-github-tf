package core

import (
	"fmt"
)

/** Public **/

func ComputeConfig(config *Config) (*Config, error) {
	if config == nil {
		return config, nil
	}

	computedConfig := NewConfig()

	errList := loadConfig(config, computedConfig)

	if len(errList) > 0 {
		return nil, ComputationError(MapToSortedList(errList))
	}

	return computedConfig, nil
}

/** Private **/

func loadConfig(config *Config, computedConfig *Config) map[string]error {
	errList := map[string]error{}

	// Compute repository config
	if config.Repos != nil {
		for k, base := range config.Repos {
			loadConfigRepo(config, computedConfig, base, errList, k)
		}
	}

	return errList
}

func loadConfigRepo(config *Config, computedConfig *Config, base *GhRepoConfig, errList map[string]error, index int) {
	if base.Name == nil {
		errList[fmt.Sprintf("Key %d", index)] = RepositoryNameIsMandatoryForRepoError(index)
	} else {
		computedRepo, computeError := ComputeRepoConfig(base, config.Templates)
		if computeError != nil {
			errList[*base.Name] = fmt.Errorf("repository %s: %w", *base.Name, computeError)
		} else {
			computedConfig.AppendRepo(computedRepo)
		}
	}
}
