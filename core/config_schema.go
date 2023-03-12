package core

func NewConfig() *Config {
	return &Config{
		Templates: &TemplatesConfig{
			Repos:             map[string]*GhRepoConfig{},
			Branches:          map[string]*GhBranchConfig{},
			BranchProtections: map[string]*GhBranchProtectionConfig{},
		},
		Repos: []*GhRepoConfig{},
	}
}

type Config struct {
	Templates *TemplatesConfig `yaml:"templates,omitempty"`
	Repos     []*GhRepoConfig  `yaml:"repos,omitempty"`

	// Org ...
	// Teams ...
}

func (c *Config) AppendRepo(repo *GhRepoConfig) {
	c.Repos = append(c.Repos, repo)
}

func (c *Config) GetRepo(name string) *GhRepoConfig {
	for _, r := range c.Repos {
		if r.Name != nil && *r.Name == name {
			return r
		}
	}

	return nil
}

type TemplatesConfig struct {
	Repos             map[string]*GhRepoConfig             `yaml:"repos,omitempty"`
	Branches          map[string]*GhBranchConfig           `yaml:"branches,omitempty"`
	BranchProtections map[string]*GhBranchProtectionConfig `yaml:"branchProtection,omitempty"`
}

func (c *TemplatesConfig) GetRepo(name string) *GhRepoConfig {
	if c.Repos == nil {
		return nil
	}

	if tpl, ok := c.Repos[name]; ok {
		return tpl
	}

	return nil
}

func (c *TemplatesConfig) GetBranch(name string) *GhBranchConfig {
	if c.Branches == nil {
		return nil
	}

	if tpl, ok := c.Branches[name]; ok {
		return tpl
	}

	return nil
}

func (c *TemplatesConfig) GetBranchProtection(name string) *GhBranchProtectionConfig {
	if c.BranchProtections == nil {
		return nil
	}

	if tpl, ok := c.BranchProtections[name]; ok {
		return tpl
	}

	return nil
}
