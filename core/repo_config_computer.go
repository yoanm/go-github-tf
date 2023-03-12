package core

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

/** Public **/

func ComputeRepoConfig(base *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error) {
	var err error

	//nolint:exhaustruct // No need here, it's just the base structure
	config := &GhRepoConfig{}

	if base == nil {
		return config, nil
	}

	if base.Name == nil {
		return nil, fmt.Errorf("repository name is mandatory")
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

func ApplyRepositoryTemplate(c *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error) {
	if c == nil {
		return c, nil
	}

	tplList, err := loadRepoTemplatesFor(c, templates)
	if err != nil {
		return nil, err
	}

	return applyRepositoryTemplate(c, tplList), nil
}

func ApplyBranchesTemplate(c *GhRepoConfig, templates *TemplatesConfig) error {
	if c == nil {
		return nil
	}

	var err error

	if c.Branches != nil {
		for k, b := range *c.Branches {
			if b, err = ApplyBranchTemplate(b, templates); err != nil {
				return fmt.Errorf("branch %s: %w", k, err)
			}

			(*c.Branches)[k] = b
		}
	}

	if c.DefaultBranch != nil {
		b := &GhBranchConfig{
			SourceBranch: nil,
			SourceSha:    nil,
			BaseGhBranchConfig: BaseGhBranchConfig{
				ConfigTemplates: c.DefaultBranch.ConfigTemplates,
				Protection:      nil,
			},
		}

		if b, err = ApplyBranchTemplate(b, templates); err != nil {
			return fmt.Errorf("default branch: %w", err)
		}

		b.Merge(
			&GhBranchConfig{
				SourceBranch: nil,
				SourceSha:    nil,
				BaseGhBranchConfig: BaseGhBranchConfig{
					ConfigTemplates: nil,
					Protection:      c.DefaultBranch.Protection,
				},
			},
		)

		c.DefaultBranch.ConfigTemplates = nil
		c.DefaultBranch.Protection = b.Protection
	}

	return nil
}

func ApplyBranchProtectionsTemplate(c *GhRepoConfig, templates *TemplatesConfig) error {
	if c == nil {
		return nil
	}

	var err error

	if c.BranchProtections != nil {
		for k, b := range *c.BranchProtections {
			if b, err = ApplyBranchProtectionTemplate(b, templates); err != nil {
				return fmt.Errorf("branch protection #%d: %w", k, err)
			}

			(*c.BranchProtections)[k] = b
		}
	}

	mapDuplicatedBranchProtection(c)

	if c.DefaultBranch != nil && c.DefaultBranch.Protection != nil {
		emptyVal := ""

		wrapper := &GhBranchProtectionConfig{
			Pattern:                      &emptyVal,
			Forbid:                       &falseString,
			BaseGhBranchProtectionConfig: *c.DefaultBranch.Protection,
		}
		if wrapper, err = ApplyBranchProtectionTemplate(wrapper, templates); err != nil {
			return fmt.Errorf("default branch: %w", err)
		}

		c.DefaultBranch.Protection = &wrapper.BaseGhBranchProtectionConfig
	}

	if c.Branches != nil {
		for k, b := range *c.Branches {
			if b.Protection == nil {
				continue
			}

			wrapper := &GhBranchProtectionConfig{
				Pattern:                      &k,
				Forbid:                       &falseString,
				BaseGhBranchProtectionConfig: *b.Protection,
			}
			if wrapper, err = ApplyBranchProtectionTemplate(wrapper, templates); err != nil {
				return fmt.Errorf("branch %s: %w", k, err)
			}

			b.Protection = &wrapper.BaseGhBranchProtectionConfig
		}
	}

	return nil
}

func ApplyBranchProtectionTemplate(
	c *GhBranchProtectionConfig,
	templates *TemplatesConfig,
) (*GhBranchProtectionConfig, error) {
	if c == nil {
		return c, nil
	}

	tplList, err := loadBranchProtectionTemplatesFor(c.ConfigTemplates, templates)
	if err != nil {
		return nil, err
	}

	return applyBranchProtectionTemplate(c, tplList), nil
}

func ApplyBranchTemplate(c *GhBranchConfig, templates *TemplatesConfig) (*GhBranchConfig, error) {
	if c == nil {
		return c, nil
	}

	tplList, err := loadBranchTemplatesFor(c.ConfigTemplates, templates)
	if err != nil {
		return nil, err
	}

	return applyBranchTemplate(c, tplList), nil
}

/** Private **/

// Not easily doable with json-schema and applying templates might create duplicates.
func mapDuplicatedBranchProtection(conf *GhRepoConfig) {
	if conf.BranchProtections != nil {
		knowPattern := map[string]int{}

		configs := conf.BranchProtections
		for k, v := range *configs {
			if v.Pattern != nil {
				if knownKey, ok := knowPattern[*v.Pattern]; ok {
					log.Warn().Msgf(
						"Repository %s: A branch protection with '%s' pattern already exists (#%d) => applying #%d as template for #%d !", //nolint:lll
						*conf.Name,
						*v.Pattern,
						knownKey,
						knownKey,
						k,
					)

					(*configs)[knownKey] = applyBranchProtectionTemplate(
						v,
						[]*GhBranchProtectionConfig{(*configs)[knownKey]},
					)
					*configs = append((*configs)[:k], (*configs)[k+1:]...) // Remove the existing config from the list
				} else {
					knowPattern[*v.Pattern] = k
				}
			}
		}
	}
}

func applyRepositoryTemplate(to *GhRepoConfig, tplList []*GhRepoConfig) *GhRepoConfig {
	if len(tplList) == 0 {
		return to
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhRepoConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(to)
	// Remove the template as it has been applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func applyBranchTemplate(to *GhBranchConfig, tplList []*GhBranchConfig) *GhBranchConfig {
	if len(tplList) == 0 {
		return to
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhBranchConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(to)
	// Remove templates as they are applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func applyBranchProtectionTemplate(
	to *GhBranchProtectionConfig,
	tplList []*GhBranchProtectionConfig,
) *GhBranchProtectionConfig {
	if len(tplList) == 0 {
		return to
	}

	//nolint:exhaustruct // No need here, it's base structure
	newConfig := &GhBranchProtectionConfig{}

	for _, tpl := range tplList {
		newConfig.Merge(tpl)
	}

	newConfig.Merge(to)
	// Remove templates as they are applied
	newConfig.ConfigTemplates = nil

	return newConfig
}

func loadRepoTemplatesFor(to *GhRepoConfig, templates *TemplatesConfig) ([]*GhRepoConfig, error) {
	if to.ConfigTemplates == nil {
		return nil, nil
	}

	const tplType = "repository"

	if templates == nil {
		return nil, fmt.Errorf("unable to load %s template, no template available", tplType)
	}

	tplList, err := loadTemplateList(
		to.ConfigTemplates,
		func(s string) *GhRepoConfig {
			return templates.GetRepo(s)
		},
		func(c *GhRepoConfig) *[]string {
			return c.ConfigTemplates
		},
		tplType,
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

	const tplType = "branch protection"

	if templates == nil {
		return nil, fmt.Errorf("unable to load %s template, no template available", tplType)
	}

	tplList, err := loadTemplateList(
		tplNameToLoad,
		func(s string) *GhBranchProtectionConfig {
			return templates.GetBranchProtection(s)
		},
		func(c *GhBranchProtectionConfig) *[]string {
			return c.ConfigTemplates
		},
		tplType,
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

	const tplType = "branch"

	if templates == nil {
		return nil, fmt.Errorf("unable to load %s template, no template available", tplType)
	}

	tplList, err := loadTemplateList(
		tplNameToLoad,
		func(s string) *GhBranchConfig {
			return templates.GetBranch(s)
		},
		func(c *GhBranchConfig) *[]string {
			return c.ConfigTemplates
		},
		tplType,
	)
	if err != nil {
		return nil, err
	}

	return tplList, nil
}
