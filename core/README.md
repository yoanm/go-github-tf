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

### func [GenerateHclRepoFiles](./terraform_writer.go#L22)

`func GenerateHclRepoFiles(configList []*GhRepoConfig) (map[string]*hclwrite.File, error)`

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

### func [WriteTerraformFiles](./terraform_writer.go#L78)

`func WriteTerraformFiles(rootPath string, files map[string]*hclwrite.File) (err error)`

## Types

### type [BaseGhBranchConfig](./repo_schema.go#L284)

`type BaseGhBranchConfig struct { ... }`

#### func (*BaseGhBranchConfig) [Merge](./repo_schema.go#L290)

`func (to *BaseGhBranchConfig) Merge(from *BaseGhBranchConfig)`

### type [BaseGhBranchProtectionConfig](./repo_schema.go#L336)

`type BaseGhBranchProtectionConfig struct { ... }`

#### func (*BaseGhBranchProtectionConfig) [Merge](./repo_schema.go#L349)

`func (to *BaseGhBranchProtectionConfig) Merge(from *BaseGhBranchProtectionConfig)`

### type [Config](./config_schema.go#L14)

`type Config struct { ... }`

#### func [ComputeConfig](./config_computer.go#L11)

`func ComputeConfig(config *Config) (*Config, error)`

#### func [NewConfig](./config_schema.go#L3)

`func NewConfig() *Config`

#### func (*Config) [AppendRepo](./config_schema.go#L22)

`func (c *Config) AppendRepo(repo *GhRepoConfig)`

#### func (*Config) [GetRepo](./config_schema.go#L26)

`func (c *Config) GetRepo(name string) *GhRepoConfig`

### type [GhBranchConfig](./repo_schema.go#L319)

`type GhBranchConfig struct { ... }`

#### func [ApplyBranchTemplate](./repo_config_computer.go#L143)

`func ApplyBranchTemplate(c *GhBranchConfig, templates *TemplatesConfig) (*GhBranchConfig, error)`

#### func [LoadBranchTemplateFromFile](./yaml_loader.go#L34)

`func LoadBranchTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhBranchConfig, error)`

#### func [LoadGhRepoBranchConfigFromFile](./yaml_loader.go#L84)

`func LoadGhRepoBranchConfigFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (config *GhBranchConfig, err error)`

LoadGhRepoBranchConfigFromFile loads the file content to GhBranchConfig struct
No schema validation will be performed, use loadBranchTemplateFromFile instead !

#### func (*GhBranchConfig) [Merge](./repo_schema.go#L325)

`func (to *GhBranchConfig) Merge(from *GhBranchConfig)`

### type [GhBranchProtectPRReviewConfig](./repo_schema.go#L424)

`type GhBranchProtectPRReviewConfig struct { ... }`

#### func (*GhBranchProtectPRReviewConfig) [Merge](./repo_schema.go#L433)

`func (to *GhBranchProtectPRReviewConfig) Merge(from *GhBranchProtectPRReviewConfig)`

### type [GhBranchProtectPRReviewDismissalsConfig](./repo_schema.go#L450)

`type GhBranchProtectPRReviewDismissalsConfig struct { ... }`

#### func (*GhBranchProtectPRReviewDismissalsConfig) [Merge](./repo_schema.go#L456)

`func (to *GhBranchProtectPRReviewDismissalsConfig) Merge(from *GhBranchProtectPRReviewDismissalsConfig)`

### type [GhBranchProtectPushesConfig](./repo_schema.go#L380)

`type GhBranchProtectPushesConfig struct { ... }`

#### func (*GhBranchProtectPushesConfig) [Merge](./repo_schema.go#L385)

`func (to *GhBranchProtectPushesConfig) Merge(from *GhBranchProtectPushesConfig)`

### type [GhBranchProtectStatusChecksConfig](./repo_schema.go#L410)

`type GhBranchProtectStatusChecksConfig struct { ... }`

#### func (*GhBranchProtectStatusChecksConfig) [Merge](./repo_schema.go#L415)

`func (to *GhBranchProtectStatusChecksConfig) Merge(from *GhBranchProtectStatusChecksConfig)`

### type [GhBranchProtectionConfig](./repo_schema.go#L394)

`type GhBranchProtectionConfig struct { ... }`

#### func [ApplyBranchProtectionTemplate](./repo_config_computer.go#L131)

`func ApplyBranchProtectionTemplate(c *GhBranchProtectionConfig, templates *TemplatesConfig) (*GhBranchProtectionConfig, error)`

#### func [LoadBranchProtectionTemplateFromFile](./yaml_loader.go#L42)

`func LoadBranchProtectionTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhBranchProtectionConfig, error)`

#### func [LoadGhRepoBranchProtectionConfigFromFile](./yaml_loader.go#L100)

`func LoadGhRepoBranchProtectionConfigFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (config *GhBranchProtectionConfig, err error)`

LoadGhRepoBranchProtectionConfigFromFile loads the file content to GhBranchProtectionConfig struct
No schema validation will be performed, use loadBranchProtectionTemplateFromFile instead !

#### func (*GhBranchProtectionConfig) [Merge](./repo_schema.go#L400)

`func (to *GhBranchProtectionConfig) Merge(from *GhBranchProtectionConfig)`

### type [GhBranchProtectionsConfig](./repo_schema.go#L90)

`type GhBranchProtectionsConfig []*GhBranchProtectionConfig`

#### func (*GhBranchProtectionsConfig) [Merge](./repo_schema.go#L92)

`func (to *GhBranchProtectionsConfig) Merge(from *GhBranchProtectionsConfig)`

### type [GhBranchesConfig](./repo_schema.go#L71)

`type GhBranchesConfig map[string]*GhBranchConfig`

#### func (*GhBranchesConfig) [Merge](./repo_schema.go#L73)

`func (to *GhBranchesConfig) Merge(from *GhBranchesConfig)`

### type [GhDefaultBranchConfig](./repo_schema.go#L304)

`type GhDefaultBranchConfig struct { ... }`

#### func (*GhDefaultBranchConfig) [Merge](./repo_schema.go#L309)

`func (to *GhDefaultBranchConfig) Merge(from *GhDefaultBranchConfig)`

### type [GhRepoConfig](./repo_schema.go#L3)

`type GhRepoConfig struct { ... }`

#### func [ApplyRepositoryTemplate](./repo_config_computer.go#L44)

`func ApplyRepositoryTemplate(c *GhRepoConfig, templates *TemplatesConfig) (newConfig *GhRepoConfig, err error)`

#### func [ComputeRepoConfig](./repo_config_computer.go#L11)

`func ComputeRepoConfig(base *GhRepoConfig, templates *TemplatesConfig) (c *GhRepoConfig, err error)`

#### func [LoadGhRepoConfigFromFile](./yaml_loader.go#L52)

`func LoadGhRepoConfigFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (config *GhRepoConfig, err error)`

LoadGhRepoConfigFromFile loads the file content to GhRepoConfig struct
No schema validation will be performed, use loadRepositoryFromFile or loadRepositoryTemplateFromFile instead !

#### func [LoadGhRepoConfigListFromFile](./yaml_loader.go#L68)

`func LoadGhRepoConfigListFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error)`

LoadGhRepoConfigListFromFile loads the file content to GhRepoConfig struct
No schema validation will be performed, use loadRepositoriesFromFile instead !

#### func [LoadRepositoriesFromFile](./yaml_loader.go#L11)

`func LoadRepositoriesFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error)`

#### func [LoadRepositoryFromFile](./yaml_loader.go#L18)

`func LoadRepositoryFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error)`

#### func [LoadRepositoryTemplateFromFile](./yaml_loader.go#L26)

`func LoadRepositoryTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error)`

#### func (*GhRepoConfig) [Merge](./repo_schema.go#L17)

`func (to *GhRepoConfig) Merge(from *GhRepoConfig)`

### type [GhRepoFileTemplatesConfig](./repo_schema.go#L185)

`type GhRepoFileTemplatesConfig struct { ... }`

#### func (*GhRepoFileTemplatesConfig) [Merge](./repo_schema.go#L190)

`func (to *GhRepoFileTemplatesConfig) Merge(from *GhRepoFileTemplatesConfig)`

### type [GhRepoMiscellaneousConfig](./repo_schema.go#L120)

`type GhRepoMiscellaneousConfig struct { ... }`

#### func (*GhRepoMiscellaneousConfig) [Merge](./repo_schema.go#L135)

`func (to *GhRepoMiscellaneousConfig) Merge(from *GhRepoMiscellaneousConfig)`

### type [GhRepoPRBranchConfig](./repo_schema.go#L269)

`type GhRepoPRBranchConfig struct { ... }`

#### func (*GhRepoPRBranchConfig) [Merge](./repo_schema.go#L274)

`func (to *GhRepoPRBranchConfig) Merge(from *GhRepoPRBranchConfig)`

### type [GhRepoPRCommitConfig](./repo_schema.go#L255)

`type GhRepoPRCommitConfig struct { ... }`

#### func (*GhRepoPRCommitConfig) [Merge](./repo_schema.go#L260)

`func (to *GhRepoPRCommitConfig) Merge(from *GhRepoPRCommitConfig)`

### type [GhRepoPRMergeStrategyConfig](./repo_schema.go#L237)

`type GhRepoPRMergeStrategyConfig struct { ... }`

#### func (*GhRepoPRMergeStrategyConfig) [Merge](./repo_schema.go#L244)

`func (to *GhRepoPRMergeStrategyConfig) Merge(from *GhRepoPRMergeStrategyConfig)`

### type [GhRepoPagesConfig](./repo_schema.go#L169)

`type GhRepoPagesConfig struct { ... }`

#### func (*GhRepoPagesConfig) [Merge](./repo_schema.go#L175)

`func (to *GhRepoPagesConfig) Merge(from *GhRepoPagesConfig)`

### type [GhRepoPullRequestConfig](./repo_schema.go#L199)

`type GhRepoPullRequestConfig struct { ... }`

#### func (*GhRepoPullRequestConfig) [Merge](./repo_schema.go#L206)

`func (to *GhRepoPullRequestConfig) Merge(from *GhRepoPullRequestConfig)`

### type [GhRepoSecurityConfig](./repo_schema.go#L466)

`type GhRepoSecurityConfig struct { ... }`

#### func (*GhRepoSecurityConfig) [Merge](./repo_schema.go#L470)

`func (to *GhRepoSecurityConfig) Merge(from *GhRepoSecurityConfig)`

### type [GhRepoTemplateConfig](./repo_schema.go#L106)

`type GhRepoTemplateConfig struct { ... }`

#### func (*GhRepoTemplateConfig) [Merge](./repo_schema.go#L111)

`func (to *GhRepoTemplateConfig) Merge(from *GhRepoTemplateConfig)`

### type [GhRepoTerraformConfig](./repo_schema.go#L478)

`type GhRepoTerraformConfig struct { ... }`

#### func (*GhRepoTerraformConfig) [Merge](./repo_schema.go#L483)

`func (to *GhRepoTerraformConfig) Merge(from *GhRepoTerraformConfig)`

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
