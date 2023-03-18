package core

import (
	"fmt"
	"sort"
)

/** Public **/

func ComputeConfig(config *Config) (*Config, error) {
	computedConfig := NewConfig()
	if config == nil {
		return computedConfig, nil
	}

	errList := loadConfig(config, computedConfig)

	if len(errList) > 0 {
		return nil, ComputationError(generateComputationErrorMessages(errList))
	}

	return computedConfig, nil
}

func generateComputationErrorMessages(errList map[string]error) []string {
	msgList := []string{}
	// sort file to always get a predictable output (for tests mostly)
	keys := make([]string, 0, len(errList))
	for k := range errList {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, file := range keys {
		msgList = append(msgList, errList[file].Error())
	}

	return msgList
}

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
		errList[fmt.Sprintf("Key %d", index)] = RepositoryNameIsMissingForRepoError(index)
	} else {
		computedRepo, computeError := ComputeRepoConfig(base, config.Templates)
		if computeError != nil {
			errList[*base.Name] = fmt.Errorf("%w: repository %s", computeError, *base.Name)
		} else {
			computedConfig.AppendRepo(computedRepo)
		}
	}
}
