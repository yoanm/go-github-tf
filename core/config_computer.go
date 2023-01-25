package core

import (
	"fmt"
	"sort"
	"strings"
)

/** Public **/

func ComputeConfig(config *Config) (*Config, error) {
	if config == nil {
		return nil, nil
	}
	computedConfig := NewConfig()

	errList := map[string]error{}

	// Compute repository config
	if config.Repos != nil {
		for k, base := range config.Repos {
			if base.Name == nil {
				errList[fmt.Sprintf("Key %d", k)] = fmt.Errorf("repository name is missing for repo #%d", k)
			} else {
				computedRepo, computeError := ComputeRepoConfig(base, config.Templates)
				if computeError != nil {
					errList[*base.Name] = fmt.Errorf("repository %s: %s", *base.Name, computeError)
				} else {
					computedConfig.AppendRepo(computedRepo)
				}
			}
		}
	}

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
