package core

import (
	"fmt"
	"sort"
	"strings"
)

/** Public **/

func ComputeConfig(config *Config) (*Config, error) {
	computedConfig := NewConfig()
	if config == nil {
		return computedConfig, nil
	}

	errList := loadConfig(config, computedConfig)

	if len(errList) > 0 {
		msgList := []string{"error during computation:"}
		// sort file to always get a predictable output (for tests mostly)
		keys := make([]string, 0, len(errList))
		for k := range errList {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, file := range keys {
			err := errList[file]
			msgList = append(msgList, fmt.Sprintf("\t - %s", err))
		}

		return nil, fmt.Errorf(strings.Join(msgList, "\n"))
	}

	return computedConfig, nil
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
		errList[fmt.Sprintf("Key %d", index)] = fmt.Errorf("repository name is missing for repo #%d", index)
	} else {
		computedRepo, computeError := ComputeRepoConfig(base, config.Templates)
		if computeError != nil {
			errList[*base.Name] = fmt.Errorf("repository %s: %w", *base.Name, computeError)
		} else {
			computedConfig.AppendRepo(computedRepo)
		}
	}
}
