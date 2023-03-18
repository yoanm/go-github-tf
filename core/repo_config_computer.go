package core

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

/** Public **/

func ComputeRepoConfig(base *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error) {
	var err error

	if base == nil {
		return base, nil
	}

	//nolint:exhaustruct // No need here, it's just the base structure
	config := &GhRepoConfig{}

	if base.Name == nil {
		return nil, ErrRepositoryNameIsMandatory
	}

	config.Merge(base)

	configTrace(fmt.Sprintf("Config after merge: %s", *base.Name), config)

	if config, err = ApplyRepositoryTemplate(config, templates); err != nil {
		return nil, err
	} else {
		configTrace(fmt.Sprintf("Config after repo template: %s", *base.Name), config)
	}

	if err = ApplyBranchesTemplate(config, templates); err != nil {
		return nil, err
	} else {
		configTrace(fmt.Sprintf("Config after branch template: %s", *base.Name), config)
	}

	if err = ApplyBranchProtectionsTemplate(config, templates); err != nil {
		return nil, err
	}

	configTrace(fmt.Sprintf("Final config: %s", *base.Name), config)

	return config, nil
}

func ApplyRepositoryTemplate(repoConfig *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error) {
	if repoConfig == nil {
		return repoConfig, nil
	}

	tplList, err := loadRepoTemplatesFor(repoConfig, templates)
	if err != nil {
		return nil, err
	}

	return applyRepositoryTemplate(repoConfig, tplList), nil
}

func ApplyBranchesTemplate(repoConfig *GhRepoConfig, templates *TemplatesConfig) error {
	if repoConfig == nil {
		return nil
	}

	var err error

	if repoConfig.Branches != nil {
		for k, b := range *repoConfig.Branches {
			if b, err = ApplyBranchTemplate(b, templates); err != nil {
				return BranchError(k, err)
			}

			(*repoConfig.Branches)[k] = b
		}
	}

	if repoConfig.DefaultBranch != nil {
		branchConfig := &GhBranchConfig{
			SourceBranch: nil,
			SourceSha:    nil,
			BaseGhBranchConfig: BaseGhBranchConfig{
				ConfigTemplates: repoConfig.DefaultBranch.ConfigTemplates,
				Protection:      nil,
			},
		}

		if branchConfig, err = ApplyBranchTemplate(branchConfig, templates); err != nil {
			return DefaultBranchError(err)
		}

		branchConfig.Merge(
			&GhBranchConfig{
				SourceBranch: nil,
				SourceSha:    nil,
				BaseGhBranchConfig: BaseGhBranchConfig{
					ConfigTemplates: nil,
					Protection:      repoConfig.DefaultBranch.Protection,
				},
			},
		)

		repoConfig.DefaultBranch.ConfigTemplates = nil
		repoConfig.DefaultBranch.Protection = branchConfig.Protection
	}

	return nil
}

func ApplyBranchProtectionsTemplate(config *GhRepoConfig, templates *TemplatesConfig) error {
	if config == nil {
		return nil
	}

	var err error

	if config.BranchProtections != nil {
		for k, b := range *config.BranchProtections {
			if b, err = ApplyBranchProtectionTemplate(b, templates); err != nil {
				return BranchProtectionError(k, err)
			}

			(*config.BranchProtections)[k] = b
		}
	}

	mapDuplicatedBranchProtection(config)

	if err2 := applyDefaultBranchProtectionTemplate(config, templates); err2 != nil {
		return err2
	}

	return applyBranchesBranchProtectionTemplate(config, templates)
}

func ApplyBranchProtectionTemplate(
	branchProtectionConfig *GhBranchProtectionConfig,
	templates *TemplatesConfig,
) (*GhBranchProtectionConfig, error) {
	if branchProtectionConfig == nil {
		return branchProtectionConfig, nil
	}

	tplList, err := loadBranchProtectionTemplatesFor(branchProtectionConfig.ConfigTemplates, templates)
	if err != nil {
		return nil, err
	}

	return applyBranchProtectionTemplate(branchProtectionConfig, tplList), nil
}

func ApplyBranchTemplate(branchConfig *GhBranchConfig, templates *TemplatesConfig) (*GhBranchConfig, error) {
	if branchConfig == nil {
		return branchConfig, nil
	}

	tplList, err := loadBranchTemplatesFor(branchConfig.ConfigTemplates, templates)
	if err != nil {
		return nil, err
	}

	return applyBranchTemplate(branchConfig, tplList), nil
}

/** Private **/

// Not easily doable with json-schema and applying templates might create duplicates.
func mapDuplicatedBranchProtection(conf *GhRepoConfig) {
	if conf.BranchProtections != nil {
		knowPattern := map[string]int{}

		configs := conf.BranchProtections
		for pattern, branchProtectionConfig := range *configs {
			if branchProtectionConfig.Pattern != nil {
				if knownKey, ok := knowPattern[*branchProtectionConfig.Pattern]; ok {
					log.Warn().Msgf(
						"Repository %s: A branch protection with '%s' pattern already exists (#%d) => applying #%d as template for #%d !", //nolint:lll
						*conf.Name,
						*branchProtectionConfig.Pattern,
						knownKey,
						knownKey,
						pattern,
					)

					(*configs)[knownKey] = applyBranchProtectionTemplate(
						branchProtectionConfig,
						[]*GhBranchProtectionConfig{(*configs)[knownKey]},
					)
					*configs = append((*configs)[:pattern], (*configs)[pattern+1:]...) // Remove the existing config from the list
				} else {
					knowPattern[*branchProtectionConfig.Pattern] = pattern
				}
			}
		}
	}
}

func applyRepositoryTemplate(toConfig *GhRepoConfig, tplList []*GhRepoConfig) *GhRepoConfig {
	if len(tplList) == 0 {
		return toConfig
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhRepoConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(toConfig)
	// Remove the template as it has been applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func applyBranchTemplate(toConfig *GhBranchConfig, tplList []*GhBranchConfig) *GhBranchConfig {
	if len(tplList) == 0 {
		return toConfig
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhBranchConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(toConfig)
	// Remove templates as they are applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func applyBranchProtectionTemplate(
	configReceiver *GhBranchProtectionConfig,
	tplList []*GhBranchProtectionConfig,
) *GhBranchProtectionConfig {
	if len(tplList) == 0 {
		return configReceiver
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhBranchProtectionConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(configReceiver)
	// Remove templates as they are applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func loadRepoTemplatesFor(toConfig *GhRepoConfig, templates *TemplatesConfig) ([]*GhRepoConfig, error) {
	if toConfig.ConfigTemplates == nil {
		return nil, nil
	}

	if templates == nil {
		return nil, NoTemplateAvailableError(RepositoryTemplateType)
	}

	tplList, err := LoadTemplateList(
		toConfig.ConfigTemplates,
		func(s string) *GhRepoConfig {
			return templates.GetRepo(s)
		},
		func(c *GhRepoConfig) *[]string {
			return c.ConfigTemplates
		},
		RepositoryTemplateType,
	)
	if err != nil {
		return nil, err
	}

	return tplList, nil
}

func loadBranchProtectionTemplatesFor(
	tplNameToLoad *[]string,
	templates *TemplatesConfig,
) ([]*GhBranchProtectionConfig, error) {
	if tplNameToLoad == nil {
		return nil, nil
	}

	if templates == nil {
		return nil, NoTemplateAvailableError(BranchProtectionTemplateType)
	}

	tplList, err := LoadTemplateList(
		tplNameToLoad,
		func(s string) *GhBranchProtectionConfig {
			return templates.GetBranchProtection(s)
		},
		func(c *GhBranchProtectionConfig) *[]string {
			return c.ConfigTemplates
		},
		BranchProtectionTemplateType,
	)
	if err != nil {
		return nil, err
	}

	return tplList, nil
}

func loadBranchTemplatesFor(tplNameToLoad *[]string, templates *TemplatesConfig) ([]*GhBranchConfig, error) {
	if tplNameToLoad == nil {
		return nil, nil
	}

	if templates == nil {
		return nil, NoTemplateAvailableError(BranchTemplateType)
	}

	tplList, err := LoadTemplateList(
		tplNameToLoad,
		func(s string) *GhBranchConfig {
			return templates.GetBranch(s)
		},
		func(c *GhBranchConfig) *[]string {
			return c.ConfigTemplates
		},
		BranchTemplateType,
	)
	if err != nil {
		return nil, err
	}

	return tplList, nil
}

func applyBranchesBranchProtectionTemplate(config *GhRepoConfig, templates *TemplatesConfig) error {
	if config.Branches != nil {
		for branchName, branchConfig := range *config.Branches {
			if branchConfig.Protection == nil {
				continue
			}

			var err error

			wrapper := &GhBranchProtectionConfig{
				Pattern:                      &branchName,
				Forbid:                       &falseString,
				BaseGhBranchProtectionConfig: *branchConfig.Protection,
			}
			if wrapper, err = ApplyBranchProtectionTemplate(wrapper, templates); err != nil {
				return BranchError(branchName, err)
			}

			branchConfig.Protection = &wrapper.BaseGhBranchProtectionConfig
		}
	}

	return nil
}

func applyDefaultBranchProtectionTemplate(config *GhRepoConfig, templates *TemplatesConfig) error {
	if config.DefaultBranch != nil && config.DefaultBranch.Protection != nil {
		emptyVal := ""

		var err error

		wrapper := &GhBranchProtectionConfig{
			Pattern:                      &emptyVal,
			Forbid:                       &falseString,
			BaseGhBranchProtectionConfig: *config.DefaultBranch.Protection,
		}
		if wrapper, err = ApplyBranchProtectionTemplate(wrapper, templates); err != nil {
			return DefaultBranchError(err)
		}

		config.DefaultBranch.Protection = &wrapper.BaseGhBranchProtectionConfig
	}

	return nil
}
