package core

type GhRepoConfig struct {
	Name              *string                    `yaml:"name,omitempty"`
	ConfigTemplates   *[]string                  `yaml:"_templates,omitempty,flow"`
	Visibility        *string                    `yaml:"visibility,omitempty"`
	Description       *string                    `yaml:"description,omitempty"`
	DefaultBranch     *GhDefaultBranchConfig     `yaml:"default-branch,omitempty"`
	Branches          *GhBranchesConfig          `yaml:"branches,omitempty"`
	BranchProtections *GhBranchProtectionsConfig `yaml:"branch-protections,omitempty"`
	PullRequests      *GhRepoPullRequestConfig   `yaml:"pull-requests,omitempty"`
	Security          *GhRepoSecurityConfig      `yaml:"security,omitempty"`
	Miscellaneous     *GhRepoMiscellaneousConfig `yaml:"misc,omitempty"`
	Terraform         *GhRepoTerraformConfig     `yaml:"terraform,omitempty"`
}

func (to *GhRepoConfig) Merge(from *GhRepoConfig) { //nolint:gocognit,cyclop // Hard to factorize, more understandable as is
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Name, from.Name)
	mergeSliceIfNotNil(&to.ConfigTemplates, from.ConfigTemplates)
	mergeStringIfNotNil(&to.Visibility, from.Visibility)
	mergeStringIfNotNil(&to.Description, from.Description)

	if from.Miscellaneous != nil {
		if to.Miscellaneous == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Miscellaneous = &GhRepoMiscellaneousConfig{}
		}

		to.Miscellaneous.Merge(from.Miscellaneous)
	}

	if from.PullRequests != nil {
		if to.PullRequests == nil {
			//nolint:exhaustruct // No need here, simple init
			to.PullRequests = &GhRepoPullRequestConfig{}
		}

		to.PullRequests.Merge(from.PullRequests)
	}

	if from.DefaultBranch != nil {
		if to.DefaultBranch == nil {
			//nolint:exhaustruct // No need here, simple init
			to.DefaultBranch = &GhDefaultBranchConfig{}
		}

		to.DefaultBranch.Merge(from.DefaultBranch)
	}

	if from.Branches != nil {
		if to.Branches == nil {
			to.Branches = &GhBranchesConfig{}
		}

		to.Branches.Merge(from.Branches)
	}

	if from.Security != nil {
		if to.Security == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Security = &GhRepoSecurityConfig{}
		}

		to.Security.Merge(from.Security)
	}

	if from.Terraform != nil {
		if to.Terraform == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Terraform = &GhRepoTerraformConfig{}
		}

		to.Terraform.Merge(from.Terraform)
	}

	if from.BranchProtections != nil {
		if to.BranchProtections == nil {
			to.BranchProtections = &GhBranchProtectionsConfig{}
		}

		to.BranchProtections.Merge(from.BranchProtections)
	}
}

type GhBranchesConfig map[string]*GhBranchConfig

func (to *GhBranchesConfig) Merge(from *GhBranchesConfig) {
	if from == nil {
		return
	}

	for k, v := range *from {
		existingVal, exists := (*to)[k]
		if exists {
			existingVal.Merge(v)
		} else {
			//nolint:exhaustruct // No need here, it's base structure
			newVal := &GhBranchConfig{}
			newVal.Merge(v)
			(*to)[k] = newVal
		}
	}
}

type GhBranchProtectionsConfig []*GhBranchProtectionConfig

func (to *GhBranchProtectionsConfig) Merge(from *GhBranchProtectionsConfig) {
	if from == nil {
		return
	}
	// Duplicate every 'from' items to avoid overflow later
	newItems := make(GhBranchProtectionsConfig, len(*from))

	for k, v := range *from {
		//nolint:exhaustruct // No need here, it's base structure
		newItem := &GhBranchProtectionConfig{}
		newItem.Merge(v)
		newItems[k] = newItem
	}

	*to = append(*to, newItems...)
}

type GhRepoTemplateConfig struct {
	Source    *string `yaml:"source,omitempty"`
	FullClone *string `yaml:"full-clone,omitempty"`
}

func (to *GhRepoTemplateConfig) Merge(from *GhRepoTemplateConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Source, from.Source)
	mergeStringIfNotNil(&to.FullClone, from.FullClone)
}

type GhRepoMiscellaneousConfig struct {
	Topics        *[]string                  `yaml:"topics,omitempty,flow"`
	AutoInit      *string                    `yaml:"auto-init,omitempty"`
	Archived      *string                    `yaml:"archived,omitempty"`
	IsTemplate    *string                    `yaml:"is-template,omitempty"`
	HomepageUrl   *string                    `yaml:"homepage-url,omitempty"`
	HasIssues     *string                    `yaml:"issues,omitempty"`
	HasWiki       *string                    `yaml:"wiki,omitempty"`
	HasProjects   *string                    `yaml:"projects,omitempty"`
	HasDownloads  *string                    `yaml:"downloads,omitempty"`
	Template      *GhRepoTemplateConfig      `yaml:"template,omitempty"`
	Pages         *GhRepoPagesConfig         `yaml:"pages,omitempty"`
	FileTemplates *GhRepoFileTemplatesConfig `yaml:"file-templates,omitempty"`
}

func (to *GhRepoMiscellaneousConfig) Merge(from *GhRepoMiscellaneousConfig) {
	if from == nil {
		return
	}

	mergeSliceIfNotNil(&to.Topics, from.Topics)
	mergeStringIfNotNil(&to.AutoInit, from.AutoInit)
	mergeStringIfNotNil(&to.Archived, from.Archived)
	mergeStringIfNotNil(&to.IsTemplate, from.IsTemplate)
	mergeStringIfNotNil(&to.HomepageUrl, from.HomepageUrl)
	mergeStringIfNotNil(&to.HasIssues, from.HasIssues)
	mergeStringIfNotNil(&to.HasWiki, from.HasWiki)
	mergeStringIfNotNil(&to.HasProjects, from.HasProjects)
	mergeStringIfNotNil(&to.HasDownloads, from.HasDownloads)

	if from.Template != nil {
		if to.Template == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Template = &GhRepoTemplateConfig{}
		}

		to.Template.Merge(from.Template)
	}

	if from.Pages != nil {
		if to.Pages == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Pages = &GhRepoPagesConfig{}
		}

		to.Pages.Merge(from.Pages)
	}

	if from.FileTemplates != nil {
		if to.FileTemplates == nil {
			//nolint:exhaustruct // No need here, simple init
			to.FileTemplates = &GhRepoFileTemplatesConfig{}
		}

		to.FileTemplates.Merge(from.FileTemplates)
	}
}

type GhRepoPagesConfig struct {
	Domain       *string `yaml:"domain,omitempty"`
	SourceBranch *string `yaml:"source-branch,omitempty"`
	SourcePath   *string `yaml:"source-path,omitempty"`
}

func (to *GhRepoPagesConfig) Merge(from *GhRepoPagesConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Domain, from.Domain)
	mergeStringIfNotNil(&to.SourceBranch, from.SourceBranch)
	mergeStringIfNotNil(&to.SourcePath, from.SourcePath)
}

type GhRepoFileTemplatesConfig struct {
	Gitignore *string `yaml:"gitignore,omitempty"`
	License   *string `yaml:"license,omitempty"`
}

func (to *GhRepoFileTemplatesConfig) Merge(from *GhRepoFileTemplatesConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Gitignore, from.Gitignore)
	mergeStringIfNotNil(&to.License, from.License)
}

type GhRepoPullRequestConfig struct {
	MergeStrategy *GhRepoPRMergeStrategyConfig `yaml:"merge-strategy,omitempty"`
	MergeCommit   *GhRepoPRCommitConfig        `yaml:"merge-commit,omitempty"`
	SquashCommit  *GhRepoPRCommitConfig        `yaml:"squash-commit,omitempty"`
	Branch        *GhRepoPRBranchConfig        `yaml:"branch,omitempty"`
}

func (to *GhRepoPullRequestConfig) Merge(from *GhRepoPullRequestConfig) {
	if from == nil {
		return
	}

	if from.MergeStrategy != nil {
		if to.MergeStrategy == nil {
			//nolint:exhaustruct // No need here, simple init
			to.MergeStrategy = &GhRepoPRMergeStrategyConfig{}
		}

		to.MergeStrategy.Merge(from.MergeStrategy)
	}

	if from.MergeCommit != nil {
		if to.MergeCommit == nil {
			//nolint:exhaustruct // No need here, simple init
			to.MergeCommit = &GhRepoPRCommitConfig{}
		}

		to.MergeCommit.Merge(from.MergeCommit)
	}

	if from.SquashCommit != nil {
		if to.SquashCommit == nil {
			//nolint:exhaustruct // No need here, simple init
			to.SquashCommit = &GhRepoPRCommitConfig{}
		}

		to.SquashCommit.Merge(from.SquashCommit)
	}

	if from.Branch != nil {
		if to.Branch == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Branch = &GhRepoPRBranchConfig{}
		}

		to.Branch.Merge(from.Branch)
	}
}

type GhRepoPRMergeStrategyConfig struct {
	AllowMerge     *string `yaml:"merge,omitempty"`
	AllowRebase    *string `yaml:"rebase,omitempty"`
	AllowSquash    *string `yaml:"squash,omitempty"`
	AllowAutoMerge *string `yaml:"auto-merge,omitempty"`
}

func (to *GhRepoPRMergeStrategyConfig) Merge(from *GhRepoPRMergeStrategyConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.AllowMerge, from.AllowMerge)
	mergeStringIfNotNil(&to.AllowRebase, from.AllowRebase)
	mergeStringIfNotNil(&to.AllowSquash, from.AllowSquash)
	mergeStringIfNotNil(&to.AllowAutoMerge, from.AllowAutoMerge)
}

type GhRepoPRCommitConfig struct {
	Title   *string `yaml:"title,omitempty"`
	Message *string `yaml:"message,omitempty"`
}

func (to *GhRepoPRCommitConfig) Merge(from *GhRepoPRCommitConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Title, from.Title)
	mergeStringIfNotNil(&to.Message, from.Message)
}

type GhRepoPRBranchConfig struct {
	SuggestUpdate       *string `yaml:"suggest-update,omitempty"`
	DeleteBranchOnMerge *string `yaml:"delete-on-merge,omitempty"`
}

func (to *GhRepoPRBranchConfig) Merge(from *GhRepoPRBranchConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.SuggestUpdate, from.SuggestUpdate)
	mergeStringIfNotNil(&to.DeleteBranchOnMerge, from.DeleteBranchOnMerge)
}

type BaseGhBranchConfig struct {
	ConfigTemplates *[]string `yaml:"_templates,omitempty,flow"`

	Protection *BaseGhBranchProtectionConfig `yaml:"protection,omitempty"`
}

func (to *BaseGhBranchConfig) Merge(from *BaseGhBranchConfig) {
	if from == nil {
		return
	}

	mergeSliceIfNotNil(&to.ConfigTemplates, from.ConfigTemplates)

	if from.Protection != nil {
		if to.Protection == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Protection = &BaseGhBranchProtectionConfig{}
		}

		to.Protection.Merge(from.Protection)
	}
}

type GhDefaultBranchConfig struct {
	Name               *string `yaml:"name,omitempty"`
	BaseGhBranchConfig `yaml:",inline"`
}

func (to *GhDefaultBranchConfig) Merge(from *GhDefaultBranchConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Name, from.Name)

	(&to.BaseGhBranchConfig).Merge(&from.BaseGhBranchConfig)
}

type GhBranchConfig struct {
	SourceBranch       *string `yaml:"source-branch,omitempty"`
	SourceSha          *string `yaml:"source-sha,omitempty"`
	BaseGhBranchConfig `yaml:",inline"`
}

func (to *GhBranchConfig) Merge(from *GhBranchConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.SourceBranch, from.SourceBranch)
	mergeStringIfNotNil(&to.SourceSha, from.SourceSha)

	(&to.BaseGhBranchConfig).Merge(&from.BaseGhBranchConfig)
}

type BaseGhBranchProtectionConfig struct {
	ConfigTemplates *[]string `yaml:"_templates,omitempty,flow"`

	EnforceAdmins        *string `yaml:"enforce-admins,omitempty"`
	AllowsDeletion       *string `yaml:"deletion,omitempty"`
	RequireLinearHistory *string `yaml:"linear-history,omitempty"`
	RequireSignedCommits *string `yaml:"signed-commits,omitempty"`

	Pushes             *GhBranchProtectPushesConfig       `yaml:"pushes,omitempty"`
	StatusChecks       *GhBranchProtectStatusChecksConfig `yaml:"status-checks,omitempty"`
	PullRequestReviews *GhBranchProtectPRReviewConfig     `yaml:"pull-request-reviews,omitempty"`
}

func (to *BaseGhBranchProtectionConfig) Merge(from *BaseGhBranchProtectionConfig) {
	if from == nil {
		return
	}

	mergeSliceIfNotNil(&to.ConfigTemplates, from.ConfigTemplates)
	mergeStringIfNotNil(&to.EnforceAdmins, from.EnforceAdmins)
	mergeStringIfNotNil(&to.AllowsDeletion, from.AllowsDeletion)
	mergeStringIfNotNil(&to.RequireLinearHistory, from.RequireLinearHistory)
	mergeStringIfNotNil(&to.RequireSignedCommits, from.RequireSignedCommits)

	if from.Pushes != nil {
		if to.Pushes == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Pushes = &GhBranchProtectPushesConfig{}
		}

		to.Pushes.Merge(from.Pushes)
	}

	if from.StatusChecks != nil {
		if to.StatusChecks == nil {
			//nolint:exhaustruct // No need here, simple init
			to.StatusChecks = &GhBranchProtectStatusChecksConfig{}
		}

		to.StatusChecks.Merge(from.StatusChecks)
	}

	if from.PullRequestReviews != nil {
		if to.PullRequestReviews == nil {
			//nolint:exhaustruct // No need here, simple init
			to.PullRequestReviews = &GhBranchProtectPRReviewConfig{}
		}

		to.PullRequestReviews.Merge(from.PullRequestReviews)
	}
}

type GhBranchProtectPushesConfig struct {
	AllowsForcePushes *string   `yaml:"force-push,omitempty"`
	PushRestrictions  *[]string `yaml:"restrict-to,omitempty,flow"`
}

func (to *GhBranchProtectPushesConfig) Merge(from *GhBranchProtectPushesConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.AllowsForcePushes, from.AllowsForcePushes)
	mergeSliceIfNotNil(&to.PushRestrictions, from.PushRestrictions)
}

type GhBranchProtectionConfig struct {
	Pattern                      *string `yaml:"pattern,omitempty"`
	Forbid                       *string `yaml:"forbid,omitempty"`
	BaseGhBranchProtectionConfig `yaml:",inline"`
}

func (to *GhBranchProtectionConfig) Merge(from *GhBranchProtectionConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Pattern, from.Pattern)
	mergeStringIfNotNil(&to.Forbid, from.Forbid)
	(&to.BaseGhBranchProtectionConfig).Merge(&from.BaseGhBranchProtectionConfig)
}

type GhBranchProtectStatusChecksConfig struct {
	Strict   *string   `yaml:"strict,omitempty"`
	Required *[]string `yaml:"required,omitempty,flow"`
}

func (to *GhBranchProtectStatusChecksConfig) Merge(from *GhBranchProtectStatusChecksConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Strict, from.Strict)
	mergeSliceIfNotNil(&to.Required, from.Required)
}

type GhBranchProtectPRReviewConfig struct {
	BypasserList          *[]string `yaml:"bypassers,omitempty,flow"`
	CodeownerApprovals    *string   `yaml:"codeowner-approvals,omitempty"`
	ResolvedConversations *string   `yaml:"resolved-conversations,omitempty"`
	ApprovalCount         *string   `yaml:"approval-count,omitempty"`

	Dismissals *GhBranchProtectPRReviewDismissalsConfig `yaml:"dismissals,omitempty"`
}

func (to *GhBranchProtectPRReviewConfig) Merge(from *GhBranchProtectPRReviewConfig) {
	if from == nil {
		return
	}

	mergeSliceIfNotNil(&to.BypasserList, from.BypasserList)
	mergeStringIfNotNil(&to.CodeownerApprovals, from.CodeownerApprovals)
	mergeStringIfNotNil(&to.ResolvedConversations, from.ResolvedConversations)
	mergeStringIfNotNil(&to.ApprovalCount, from.ApprovalCount)

	if from.Dismissals != nil {
		if to.Dismissals == nil {
			//nolint:exhaustruct // No need here, simple init
			to.Dismissals = &GhBranchProtectPRReviewDismissalsConfig{}
		}

		to.Dismissals.Merge(from.Dismissals)
	}
}

type GhBranchProtectPRReviewDismissalsConfig struct {
	Staled     *string   `yaml:"staled,omitempty"`
	Restrict   *string   `yaml:"restrict,omitempty"`
	RestrictTo *[]string `yaml:"restrict-to,omitempty,flow"`
}

func (to *GhBranchProtectPRReviewDismissalsConfig) Merge(from *GhBranchProtectPRReviewDismissalsConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.Staled, from.Staled)
	mergeStringIfNotNil(&to.Restrict, from.Restrict)
	mergeSliceIfNotNil(&to.RestrictTo, from.RestrictTo)
}

type GhRepoSecurityConfig struct {
	VulnerabilityAlerts *string `yaml:"vulnerability-alerts,omitempty"`
}

func (to *GhRepoSecurityConfig) Merge(from *GhRepoSecurityConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.VulnerabilityAlerts, from.VulnerabilityAlerts)
}

type GhRepoTerraformConfig struct {
	ArchiveOnDestroy                    *string `yaml:"archive-on-destroy,omitempty"`
	IgnoreVulnerabilityAlertsDuringRead *string `yaml:"ignore-vulnerability-alerts-during-read,omitempty"`
}

func (to *GhRepoTerraformConfig) Merge(from *GhRepoTerraformConfig) {
	if from == nil {
		return
	}

	mergeStringIfNotNil(&to.ArchiveOnDestroy, from.ArchiveOnDestroy)
	mergeStringIfNotNil(&to.IgnoreVulnerabilityAlertsDuringRead, from.IgnoreVulnerabilityAlertsDuringRead)
}

// mergeStringIfNotNil ensures that updating 'from' afterward doesn't affect 'to' and vice versa
// It override the string behind 'to' pointer by the string behind 'from' pointer.
func mergeStringIfNotNil(to **string, from *string) {
	if from != nil {
		// !! to != nil is assume here !!
		if *to == nil {
			empty := ""
			*to = &empty
		}
		// Set the underlying value instead of the pointer to avoid overflow later
		**to = *from
	}
}

// mergeSliceIfNotNil ensures that updating 'from' afterward doesn't affect 'to' and vice versa
// It use append function to create a new slice combining 'from' and 'to' items.
func mergeSliceIfNotNil(to **[]string, from *[]string) {
	if from != nil {
		// !! to != nil is assume here !!
		if *to == nil {
			*to = &[]string{}
		}

		**to = append(**to, *from...)
	}
}
