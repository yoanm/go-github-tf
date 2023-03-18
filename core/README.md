# core

Package core provides core functionality for GitHub config to terraform file conversion

## Constants

```golang
const (
    RepositoryTemplateType       = "repository"
    BranchTemplateType           = "branch"
    BranchProtectionTemplateType = "branch protection"

    TemplateMaxDepth = 10
    TemplateMaxCount = 10
)
```

```golang
const (
    DefaultBranchIdentifier = "default"
)
```

## Variables

```golang
var (
    ErrRepositoryNameIsMandatory = errors.New("repository name is mandatory")

    ErrWorkspacePathDoesntExist              = errors.New("workspace path doesn't exist")
    ErrWorkspacePathIsExpectedToBeADirectory = errors.New("workspace path is expected to be a directory")

    ErrNoTemplateAvailable                 = errors.New("not found as none available")
    ErrNoRepositoryTemplateAvailable       = fmt.Errorf("%s template %w", RepositoryTemplateType, ErrNoTemplateAvailable)
    ErrNoBranchTemplateAvailable           = fmt.Errorf("%s template %w", BranchTemplateType, ErrNoTemplateAvailable)
    ErrNoBranchProtectionTemplateAvailable = fmt.Errorf(
        "%s template %w",
        BranchProtectionTemplateType,
        ErrNoTemplateAvailable,
    )

    ErrTemplateNotFound                 = errors.New("not found")
    ErrRepositoryTemplateNotFound       = fmt.Errorf("%s template %w", RepositoryTemplateType, ErrTemplateNotFound)
    ErrBranchTemplateNotFound           = fmt.Errorf("%s template %w", BranchTemplateType, ErrTemplateNotFound)
    ErrBranchProtectionTemplateNotFound = fmt.Errorf("%s template %w", BranchProtectionTemplateType, ErrTemplateNotFound)

    ErrMaxTemplateCount = errors.New("maximum template count reached")
    ErrMaxTemplateDepth = errors.New("maximum template depth reached")

    ErrDuringWriteTerraformFiles = errors.New("error while writing terraform files")
    ErrDuringFileGeneration      = errors.New("error while generating files")
    ErrDuringComputation         = errors.New("error during computation")

    ErrSchemaValidation        = errors.New("schema validation error")
    ErrEmptySchema             = errors.New("empty schema")
    ErrSchemaNotFound          = errors.New("schema not found")
    ErrSchemaIsNil             = errors.New("schema is nil")
    ErrDuringSchemaCompilation = errors.New("error during schema compilation")

    ErrFileError             = errors.New("file")
    ErrBranchError           = errors.New("branch")
    ErrDefaultBranchError    = errors.New("default branch")
    ErrBranchProtectionError = errors.New("branch protection")
)
```

```golang
var (
    //nolint:gochecknoglobals //Easier to manage it as exported variable
    YamlAnchorDirectory *string
    //nolint:gochecknoglobals //Easier to manage it as exported variable
    Schemas = &SchemaList{
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
)
```

## Functions

### func [ApplyBranchProtectionsTemplate](./repo_config_computer.go#L112)

`func ApplyBranchProtectionsTemplate(config *GhRepoConfig, templates *TemplatesConfig) error`

### func [ApplyBranchesTemplate](./repo_config_computer.go#L63)

`func ApplyBranchesTemplate(repoConfig *GhRepoConfig, templates *TemplatesConfig) error`

### func [BranchError](./errors.go#L49)

`func BranchError(branch string, err error) error`

### func [BranchProtectionError](./errors.go#L57)

`func BranchProtectionError(index int, err error) error`

### func [ComputationError](./errors.go#L77)

`func ComputationError(errList []error) error`

### func [DefaultBranchError](./errors.go#L53)

`func DefaultBranchError(err error) error`

### func [EmptySchemaError](./errors.go#L161)

`func EmptySchemaError(url string) error`

### func [FileError](./errors.go#L61)

`func FileError(filepath string, err error) error`

### func [FileGenerationError](./errors.go#L65)

`func FileGenerationError(msgList []string) error`

### func [GenerateHclRepoFiles](./terraform_writer.go#L19)

`func GenerateHclRepoFiles(configList []*GhRepoConfig) (map[string]*hclwrite.File, error)`

### func [JoinErrors](./errors.go#L199)

`func JoinErrors(errList []error, separator string) error`

### func [LoadTemplate](./template_loader.go#L12)

`func LoadTemplate[T any](
    tplName string,
    loaderFn func(s string) *T,
    finderFn func(c *T) *[]string,
    tplType string,
    path ...string,
) ([]*T, error)`

### func [LoadTemplateList](./template_loader.go#L42)

`func LoadTemplateList[T any](
    tplNameList *[]string,
    loaderFn func(s string) *T,
    finderFn func(c *T) *[]string,
    tplType string,
    path ...string,
) ([]*T, error)`

### func [MapBranchToBranchProtectionRes](./gh2tf_repo_mapper.go#L400)

`func MapBranchToBranchProtectionRes(
    name string,
    branchConfig *GhBranchConfig,
    valGen tfsig.ValueGenerator,
    repo *GhRepoConfig,
    repoTfId string,
    links ...MapperLink,
) *ghbranchprotect.Config`

### func [MapDefaultBranchToBranchProtectionRes](./gh2tf_repo_mapper.go#L357)

`func MapDefaultBranchToBranchProtectionRes(
    branchConfig *GhDefaultBranchConfig,
    valGen tfsig.ValueGenerator,
    repo *GhRepoConfig,
    repoTfId string,
    links ...MapperLink,
) *ghbranchprotect.Config`

### func [MapToBranchProtectionRes](./gh2tf_repo_mapper.go#L430)

`func MapToBranchProtectionRes(
    branchProtectionConfig *GhBranchProtectionConfig,
    valGen tfsig.ValueGenerator,
    repo *GhRepoConfig,
    repoTfId string,
    links ...MapperLink,
) *ghbranchprotect.Config`

### func [MapToBranchRes](./gh2tf_repo_mapper.go#L225)

`func MapToBranchRes(
    name string,
    branchConfig *GhBranchConfig,
    valGen tfsig.ValueGenerator,
    repo *GhRepoConfig,
    repoTfId string,
    links ...MapperLink,
) *ghbranch.Config`

### func [MapToDefaultBranchRes](./gh2tf_repo_mapper.go#L308)

`func MapToDefaultBranchRes(
    branchConfig *GhDefaultBranchConfig,
    valGen tfsig.ValueGenerator,
    repo *GhRepoConfig,
    repoTfId string,
    links ...MapperLink,
) *ghbranchdefault.Config`

### func [MapToRepositoryRes](./gh2tf_repo_mapper.go#L28)

`func MapToRepositoryRes(repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) *ghrepository.Config`

### func [MaxTemplateCountReachedError](./errors.go#L119)

`func MaxTemplateCountReachedError(tplType string, path []string) error`

### func [MaxTemplateDepthReachedError](./errors.go#L135)

`func MaxTemplateDepthReachedError(tplType string, path []string) error`

### func [NewHclRepository](./repo_writer.go#L17)

`func NewHclRepository(repoTfId string, repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator) *hclwrite.File`

### func [NoTemplateAvailableError](./errors.go#L106)

`func NoTemplateAvailableError(tplType string) error`

### func [RepositoryNameIsMandatoryForConfigIndexError](./errors.go#L81)

`func RepositoryNameIsMandatoryForConfigIndexError(index int) error`

### func [RepositoryNameIsMandatoryForRepoError](./errors.go#L85)

`func RepositoryNameIsMandatoryForRepoError(index int) error`

### func [SchemaCompilationError](./errors.go#L173)

`func SchemaCompilationError(url string, msg string) error`

### func [SchemaIsNilError](./errors.go#L169)

`func SchemaIsNilError(url string) error`

### func [SchemaNotFoundError](./errors.go#L165)

`func SchemaNotFoundError(url string) error`

### func [SchemaValidationError](./errors.go#L151)

`func SchemaValidationError(path string, location string, msg string) error`

### func [SortErrorsByKey](./errors.go#L181)

`func SortErrorsByKey(errList map[string]error) []error`

### func [TerraformFileWritingErrors](./errors.go#L177)

`func TerraformFileWritingErrors(errList []error) error`

### func [UnknownTemplateError](./errors.go#L89)

`func UnknownTemplateError(tplType string, tplName string) error`

### func [ValidateBranchProtectionTemplateConfig](./yaml_validator.go#L113)

`func ValidateBranchProtectionTemplateConfig(filePath string) error`

### func [ValidateBranchTemplateConfig](./yaml_validator.go#L104)

`func ValidateBranchTemplateConfig(filePath string) error`

### func [ValidateRepositoryConfig](./yaml_validator.go#L77)

`func ValidateRepositoryConfig(filePath string) error`

### func [ValidateRepositoryConfigs](./yaml_validator.go#L86)

`func ValidateRepositoryConfigs(filePath string) error`

### func [ValidateRepositoryTemplateConfig](./yaml_validator.go#L95)

`func ValidateRepositoryTemplateConfig(filePath string) error`

### func [WorkspacePathDoesntExistError](./errors.go#L69)

`func WorkspacePathDoesntExistError(path string) error`

### func [WorkspacePathIsExpectedToBeADirectoryError](./errors.go#L73)

`func WorkspacePathIsExpectedToBeADirectoryError(path string) error`

### func [WriteTerraformFiles](./terraform_writer.go#L55)

`func WriteTerraformFiles(rootPath string, files map[string]*hclwrite.File) error`

## Types

### type [BaseGhBranchConfig](./repo_schema.go#L339)

`type BaseGhBranchConfig struct { ... }`

#### func (*BaseGhBranchConfig) [Merge](./repo_schema.go#L346)

`func (to *BaseGhBranchConfig) Merge(from *BaseGhBranchConfig)`

### type [BaseGhBranchProtectionConfig](./repo_schema.go#L395)

`type BaseGhBranchProtectionConfig struct { ... }`

#### func (*BaseGhBranchProtectionConfig) [Merge](./repo_schema.go#L412)

`func (to *BaseGhBranchProtectionConfig) Merge(from *BaseGhBranchProtectionConfig)`

### type [Config](./config_schema.go#L14)

`type Config struct { ... }`

#### func [ComputeConfig](./config_computer.go#L9)

`func ComputeConfig(config *Config) (*Config, error)`

#### func [NewConfig](./config_schema.go#L3)

`func NewConfig() *Config`

#### func (*Config) [AppendRepo](./config_schema.go#L22)

`func (c *Config) AppendRepo(repo *GhRepoConfig)`

#### func (*Config) [GetRepo](./config_schema.go#L26)

`func (c *Config) GetRepo(name string) *GhRepoConfig`

### type [GhBranchConfig](./repo_schema.go#L378)

`type GhBranchConfig struct { ... }`

#### func [ApplyBranchTemplate](./repo_config_computer.go#L154)

`func ApplyBranchTemplate(branchConfig *GhBranchConfig, templates *TemplatesConfig) (*GhBranchConfig, error)`

#### func [LoadBranchTemplateFromFile](./yaml_loader.go#L34)

`func LoadBranchTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhBranchConfig, error)`

#### func [LoadGhRepoBranchConfigFromFile](./yaml_loader.go#L94)

`func LoadGhRepoBranchConfigFromFile(
    filePath string,
    decoderOpts ...yaml.DecodeOption,
) (*GhBranchConfig, error)`

LoadGhRepoBranchConfigFromFile loads the file content to GhBranchConfig struct
No schema validation will be performed, use loadBranchTemplateFromFile instead !

#### func (*GhBranchConfig) [Merge](./repo_schema.go#L384)

`func (to *GhBranchConfig) Merge(from *GhBranchConfig)`

### type [GhBranchProtectPRReviewConfig](./repo_schema.go#L496)

`type GhBranchProtectPRReviewConfig struct { ... }`

#### func (*GhBranchProtectPRReviewConfig) [Merge](./repo_schema.go#L505)

`func (to *GhBranchProtectPRReviewConfig) Merge(from *GhBranchProtectPRReviewConfig)`

### type [GhBranchProtectPRReviewDismissalsConfig](./repo_schema.go#L525)

`type GhBranchProtectPRReviewDismissalsConfig struct { ... }`

#### func (*GhBranchProtectPRReviewDismissalsConfig) [Merge](./repo_schema.go#L531)

`func (to *GhBranchProtectPRReviewDismissalsConfig) Merge(from *GhBranchProtectPRReviewDismissalsConfig)`

### type [GhBranchProtectPushesConfig](./repo_schema.go#L451)

`type GhBranchProtectPushesConfig struct { ... }`

#### func (*GhBranchProtectPushesConfig) [Merge](./repo_schema.go#L457)

`func (to *GhBranchProtectPushesConfig) Merge(from *GhBranchProtectPushesConfig)`

### type [GhBranchProtectStatusChecksConfig](./repo_schema.go#L482)

`type GhBranchProtectStatusChecksConfig struct { ... }`

#### func (*GhBranchProtectStatusChecksConfig) [Merge](./repo_schema.go#L487)

`func (to *GhBranchProtectStatusChecksConfig) Merge(from *GhBranchProtectStatusChecksConfig)`

### type [GhBranchProtectionConfig](./repo_schema.go#L466)

`type GhBranchProtectionConfig struct { ... }`

#### func [ApplyBranchProtectionTemplate](./repo_config_computer.go#L138)

`func ApplyBranchProtectionTemplate(
    branchProtectionConfig *GhBranchProtectionConfig,
    templates *TemplatesConfig,
) (*GhBranchProtectionConfig, error)`

#### func [LoadBranchProtectionTemplateFromFile](./yaml_loader.go#L42)

`func LoadBranchProtectionTemplateFromFile(
    filePath string,
    decoderOpts ...yaml.DecodeOption,
) (*GhBranchProtectionConfig, error)`

#### func [LoadGhRepoBranchProtectionConfigFromFile](./yaml_loader.go#L119)

`func LoadGhRepoBranchProtectionConfigFromFile(
    filePath string,
    decoderOpts ...yaml.DecodeOption,
) (*GhBranchProtectionConfig, error)`

LoadGhRepoBranchProtectionConfigFromFile loads the file content to GhBranchProtectionConfig struct
No schema validation will be performed, use loadBranchProtectionTemplateFromFile instead !

#### func (*GhBranchProtectionConfig) [Merge](./repo_schema.go#L472)

`func (to *GhBranchProtectionConfig) Merge(from *GhBranchProtectionConfig)`

### type [GhBranchProtectionsConfig](./repo_schema.go#L115)

`type GhBranchProtectionsConfig []*GhBranchProtectionConfig`

#### func (*GhBranchProtectionsConfig) [Merge](./repo_schema.go#L117)

`func (to *GhBranchProtectionsConfig) Merge(from *GhBranchProtectionsConfig)`

### type [GhBranchesConfig](./repo_schema.go#L95)

`type GhBranchesConfig map[string]*GhBranchConfig`

#### func (*GhBranchesConfig) [Merge](./repo_schema.go#L97)

`func (to *GhBranchesConfig) Merge(from *GhBranchesConfig)`

### type [GhDefaultBranchConfig](./repo_schema.go#L363)

`type GhDefaultBranchConfig struct { ... }`

#### func (*GhDefaultBranchConfig) [Merge](./repo_schema.go#L368)

`func (to *GhDefaultBranchConfig) Merge(from *GhDefaultBranchConfig)`

### type [GhRepoConfig](./repo_schema.go#L6)

`type GhRepoConfig struct { ... }`

#### func [ApplyRepositoryTemplate](./repo_config_computer.go#L50)

`func ApplyRepositoryTemplate(repoConfig *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error)`

#### func [ComputeRepoConfig](./repo_config_computer.go#L11)

`func ComputeRepoConfig(base *GhRepoConfig, templates *TemplatesConfig) (*GhRepoConfig, error)`

#### func [LoadGhRepoConfigFromFile](./yaml_loader.go#L55)

`func LoadGhRepoConfigFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error)`

LoadGhRepoConfigFromFile loads the file content to GhRepoConfig struct
No schema validation will be performed, use loadRepositoryFromFile or loadRepositoryTemplateFromFile instead !

#### func [LoadGhRepoConfigListFromFile](./yaml_loader.go#L77)

`func LoadGhRepoConfigListFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error)`

LoadGhRepoConfigListFromFile loads the file content to GhRepoConfig struct
No schema validation will be performed, use loadRepositoriesFromFile instead !

#### func [LoadRepositoriesFromFile](./yaml_loader.go#L10)

`func LoadRepositoriesFromFile(filePath string, decoderOpts ...yaml.DecodeOption) ([]*GhRepoConfig, error)`

#### func [LoadRepositoryFromFile](./yaml_loader.go#L18)

`func LoadRepositoryFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error)`

#### func [LoadRepositoryTemplateFromFile](./yaml_loader.go#L26)

`func LoadRepositoryTemplateFromFile(filePath string, decoderOpts ...yaml.DecodeOption) (*GhRepoConfig, error)`

#### func (*GhRepoConfig) [Merge](./repo_schema.go#L23)

`func (to *GhRepoConfig) Merge(from *GhRepoConfig)`

### type [GhRepoFileTemplatesConfig](./repo_schema.go#L226)

`type GhRepoFileTemplatesConfig struct { ... }`

#### func (*GhRepoFileTemplatesConfig) [Merge](./repo_schema.go#L231)

`func (to *GhRepoFileTemplatesConfig) Merge(from *GhRepoFileTemplatesConfig)`

### type [GhRepoMiscellaneousConfig](./repo_schema.go#L148)

`type GhRepoMiscellaneousConfig struct { ... }`

#### func (*GhRepoMiscellaneousConfig) [Merge](./repo_schema.go#L167)

`func (to *GhRepoMiscellaneousConfig) Merge(from *GhRepoMiscellaneousConfig)`

### type [GhRepoPRBranchConfig](./repo_schema.go#L325)

`type GhRepoPRBranchConfig struct { ... }`

#### func (*GhRepoPRBranchConfig) [Merge](./repo_schema.go#L330)

`func (to *GhRepoPRBranchConfig) Merge(from *GhRepoPRBranchConfig)`

### type [GhRepoPRCommitConfig](./repo_schema.go#L311)

`type GhRepoPRCommitConfig struct { ... }`

#### func (*GhRepoPRCommitConfig) [Merge](./repo_schema.go#L316)

`func (to *GhRepoPRCommitConfig) Merge(from *GhRepoPRCommitConfig)`

### type [GhRepoPRMergeStrategyConfig](./repo_schema.go#L289)

`type GhRepoPRMergeStrategyConfig struct { ... }`

#### func (*GhRepoPRMergeStrategyConfig) [Merge](./repo_schema.go#L300)

`func (to *GhRepoPRMergeStrategyConfig) Merge(from *GhRepoPRMergeStrategyConfig)`

### type [GhRepoPagesConfig](./repo_schema.go#L210)

`type GhRepoPagesConfig struct { ... }`

#### func (*GhRepoPagesConfig) [Merge](./repo_schema.go#L216)

`func (to *GhRepoPagesConfig) Merge(from *GhRepoPagesConfig)`

### type [GhRepoPullRequestConfig](./repo_schema.go#L240)

`type GhRepoPullRequestConfig struct { ... }`

#### func (*GhRepoPullRequestConfig) [Merge](./repo_schema.go#L247)

`func (to *GhRepoPullRequestConfig) Merge(from *GhRepoPullRequestConfig)`

### type [GhRepoSecurityConfig](./repo_schema.go#L541)

`type GhRepoSecurityConfig struct { ... }`

#### func (*GhRepoSecurityConfig) [Merge](./repo_schema.go#L545)

`func (to *GhRepoSecurityConfig) Merge(from *GhRepoSecurityConfig)`

### type [GhRepoTemplateConfig](./repo_schema.go#L134)

`type GhRepoTemplateConfig struct { ... }`

#### func (*GhRepoTemplateConfig) [Merge](./repo_schema.go#L139)

`func (to *GhRepoTemplateConfig) Merge(from *GhRepoTemplateConfig)`

### type [GhRepoTerraformConfig](./repo_schema.go#L553)

`type GhRepoTerraformConfig struct { ... }`

#### func (*GhRepoTerraformConfig) [Merge](./repo_schema.go#L558)

`func (to *GhRepoTerraformConfig) Merge(from *GhRepoTerraformConfig)`

### type [MapperLink](./gh2tf_repo_mapper.go#L14)

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

### type [SchemaList](./yaml_validator_schemas.go#L15)

`type SchemaList map[string]*Schema`

#### func (*SchemaList) [Compile](./yaml_validator_schemas.go#L66)

`func (s *SchemaList) Compile(url string) (*jsonschema.Schema, error)`

#### func (*SchemaList) [Find](./yaml_validator_schemas.go#L53)

`func (s *SchemaList) Find(url string) (*Schema, error)`

#### func (*SchemaList) [FindCompiled](./yaml_validator_schemas.go#L17)

`func (s *SchemaList) FindCompiled(url string) *jsonschema.Schema`

#### func (*SchemaList) [FindContent](./yaml_validator_schemas.go#L40)

`func (s *SchemaList) FindContent(url string) (*string, error)`

### type [TemplatesConfig](./config_schema.go#L36)

`type TemplatesConfig struct { ... }`

#### func (*TemplatesConfig) [GetBranch](./config_schema.go#L54)

`func (c *TemplatesConfig) GetBranch(name string) *GhBranchConfig`

#### func (*TemplatesConfig) [GetBranchProtection](./config_schema.go#L66)

`func (c *TemplatesConfig) GetBranchProtection(name string) *GhBranchProtectionConfig`

#### func (*TemplatesConfig) [GetRepo](./config_schema.go#L42)

`func (c *TemplatesConfig) GetRepo(name string) *GhRepoConfig`

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
