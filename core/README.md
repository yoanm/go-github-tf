# core

## Constants

```golang
const DefaultBranchIdentifier = "default"
```

```golang
const TemplateMaxCount = 10
```

```golang
const TemplateMaxDepth = 10
```

## Variables

```golang
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
```

```golang
var (
    YamlAnchorDirectory *string
)
```

## Functions

### func [ApplyBranchProtectionsTemplate](./repo_config_computer.go#L85)

`func ApplyBranchProtectionsTemplate(c *GhRepoConfig, templates *TemplatesConfig) (err error)`

### func [ApplyBranchesTemplate](./repo_config_computer.go#L57)

`func ApplyBranchesTemplate(c *GhRepoConfig, templates *TemplatesConfig) (err error)`

### func [MapBranchToBranchProtectionRes](./gh2tf_repo_mapper.go#L295)

`func MapBranchToBranchProtectionRes(name string, c *GhBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config`

### func [MapDefaultBranchToBranchProtectionRes](./gh2tf_repo_mapper.go#L261)

`func MapDefaultBranchToBranchProtectionRes(c *GhDefaultBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config`

### func [MapToBranchProtectionRes](./gh2tf_repo_mapper.go#L317)

`func MapToBranchProtectionRes(c *GhBranchProtectionConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config`

### func [MapToBranchRes](./gh2tf_repo_mapper.go#L167)

`func MapToBranchRes(name string, c *GhBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranch.Config`

### func [MapToDefaultBranchRes](./gh2tf_repo_mapper.go#L222)

`func MapToDefaultBranchRes(c *GhDefaultBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchdefault.Config`

### func [MapToRepositoryRes](./gh2tf_repo_mapper.go#L24)

`func MapToRepositoryRes(c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) *ghrepository.Config`

### func [NewHclRepository](./repo_writer.go#L17)

`func NewHclRepository(repoTfId string, c *GhRepoConfig, valGen tfsig.ValueGenerator) *hclwrite.File`

### func [ValidateBranchProtectionTemplateConfig](./yaml_validator.go#L111)

`func ValidateBranchProtectionTemplateConfig(filePath string) (err error)`

### func [ValidateBranchTemplateConfig](./yaml_validator.go#L102)

`func ValidateBranchTemplateConfig(filePath string) (err error)`

### func [ValidateRepositoryConfig](./yaml_validator.go#L75)

`func ValidateRepositoryConfig(filePath string) (err error)`

### func [ValidateRepositoryConfigs](./yaml_validator.go#L84)

`func ValidateRepositoryConfigs(filePath string) (err error)`

### func [ValidateRepositoryTemplateConfig](./yaml_validator.go#L93)

`func ValidateRepositoryTemplateConfig(filePath string) (err error)`

## Types

### type [Config](./config_schema.go#L14)

`type Config struct { ... }`

#### func [NewConfig](./config_schema.go#L3)

`func NewConfig() *Config`

#### func (*Config) [AppendRepo](./config_schema.go#L22)

`func (c *Config) AppendRepo(repo *GhRepoConfig)`

#### func (*Config) [GetRepo](./config_schema.go#L26)

`func (c *Config) GetRepo(name string) *GhRepoConfig`

### type [MapperLink](./gh2tf_repo_mapper.go#L15)

`type MapperLink int`

#### Constants

```golang
const (
    LinkToRepository MapperLink = iota
    LinkToBranch
)
```

### type [Schema](./yaml_validator_schemas.go#L9)

`type Schema struct { ... }`

### type [SchemaList](./yaml_validator_schemas.go#L14)

`type SchemaList map[string]*Schema`

#### func (*SchemaList) [Compile](./yaml_validator_schemas.go#L57)

`func (list *SchemaList) Compile(url string) (*jsonschema.Schema, error)`

#### func (*SchemaList) [Find](./yaml_validator_schemas.go#L45)

`func (s *SchemaList) Find(url string) (*Schema, error)`

#### func (*SchemaList) [FindCompiled](./yaml_validator_schemas.go#L16)

`func (list *SchemaList) FindCompiled(url string) *jsonschema.Schema`

#### func (*SchemaList) [FindContent](./yaml_validator_schemas.go#L33)

`func (list *SchemaList) FindContent(url string) (*string, error)`

### type [TemplatesConfig](./config_schema.go#L36)

`type TemplatesConfig struct { ... }`

#### func (*TemplatesConfig) [GetBranch](./config_schema.go#L53)

`func (c *TemplatesConfig) GetBranch(name string) *GhBranchConfig`

#### func (*TemplatesConfig) [GetBranchProtection](./config_schema.go#L64)

`func (c *TemplatesConfig) GetBranchProtection(name string) *GhBranchProtectionConfig`

#### func (*TemplatesConfig) [GetRepo](./config_schema.go#L42)

`func (c *TemplatesConfig) GetRepo(name string) *GhRepoConfig`

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
