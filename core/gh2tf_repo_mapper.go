package core

import (
	"fmt"
	"strings"

	"github.com/yoanm/go-gh2tf/ghbranch"
	"github.com/yoanm/go-gh2tf/ghbranchdefault"
	"github.com/yoanm/go-gh2tf/ghbranchprotect"
	"github.com/yoanm/go-gh2tf/ghrepository"
	"github.com/yoanm/go-tfsig"
)

type MapperLink int

const (
	LinkToRepository MapperLink = iota
	LinkToBranch
)

const (
	DefaultBranchIdentifier = "default"
)

//nolint:gochecknoglobals // Easier than duplicate it everywhere needed
var falseString = "false"

// Replace branch protection pattern's special chars by dedicated string in order to avoid ID collision (pattern is used for TF ressource generation)
// in case every special chars are replaced by the same string => "?.*", "[.]" or "?/?" would all lead to the same "---" ID for instance
const patternToIdReplacer  = strings.NewReplacer(
	".", "_DOT_",
	"/", "_SLASH_",
	"\\", "_ESC_"
	// fnmatch special chars
	"*", "_STAR_",
	"[", "_SEQ-O_",
	"]", "_SEQ-C_",
	"?", "_Q-MARK_",
	"!", "_EX-MARK_",
)

func MapToRepositoryRes(repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) *ghrepository.Config {
	if repoConfig == nil {
		return nil
	}

	topics,
		autoInit,
		archived,
		homepageUrl,
		hasIssues,
		hasProjects,
		hasWiki,
		hasDownloads,
		page,
		template := mapMiscellaneous(repoConfig, valGen)

	allowMergeCommit,
		allowRebaseMerge,
		allowSquashMerge,
		allowAutoMerge,
		mergeCommitTitle,
		mergeCommitMessage,
		squashMergeCommitTitle,
		squashMergeCommitMessage,
		deleteBranchOnMerge := mapPullRequest(repoConfig)
	// Security
	var vulnerabilityAlerts *string
	if repoConfig.Security != nil {
		vulnerabilityAlerts = repoConfig.Security.VulnerabilityAlerts
	}

	// Terraform
	var archiveOnDestroy *string
	if repoConfig.Terraform != nil {
		archiveOnDestroy = repoConfig.Terraform.ArchiveOnDestroy
	}

	return &ghrepository.Config{
		ValueGenerator: valGen,
		Identifier:     repoTfId,

		Name:                     repoConfig.Name,
		Visibility:               repoConfig.Visibility,
		Archived:                 archived,
		Description:              repoConfig.Description,
		AutoInit:                 autoInit,
		HasIssues:                hasIssues,
		HasProjects:              hasProjects,
		HasWiki:                  hasWiki,
		HasDownloads:             hasDownloads,
		HomepageUrl:              homepageUrl,
		Topics:                   topics,
		VulnerabilityAlerts:      vulnerabilityAlerts,
		AllowMergeCommit:         allowMergeCommit,
		AllowRebaseMerge:         allowRebaseMerge,
		AllowSquashMerge:         allowSquashMerge,
		AllowAutoMerge:           allowAutoMerge,
		MergeCommitTitle:         mergeCommitTitle,
		MergeCommitMessage:       mergeCommitMessage,
		SquashMergeCommitTitle:   squashMergeCommitTitle,
		SquashMergeCommitMessage: squashMergeCommitMessage,
		DeleteBranchOnMerge:      deleteBranchOnMerge,
		ArchiveOnDestroy:         archiveOnDestroy,
		// SuggestUpdate:            suggestUpdate,
		// IsTemplate:               isTemplate,
		// LicenseTemplate:          licenseTemplate
		// GitignoreTemplate:        gitignoreTemplate

		Page:     page,
		Template: template,
	}
}

func MapToBranchRes(
	name string,
	branchConfig *GhBranchConfig,
	valGen tfsig.ValueGenerator,
	repo *GhRepoConfig,
	repoTfId string,
	links ...MapperLink,
) *ghbranch.Config {
	if branchConfig == nil {
		return nil
	}

	var (
		sourceBranch *string
		repoName     *string
	)

	identifier := fmt.Sprintf("%s-%s", repoTfId, tfsig.ToTerraformIdentifier(name))

	if branchConfig.SourceBranch != nil && *branchConfig.SourceBranch != name {
		sourceBranch = branchConfig.SourceBranch
	}

	if repo != nil {
		repoName = repo.Name
	}

	for _, link := range links {
		repoName, sourceBranch = applyBranchResLink(link, repoTfId, repoName, sourceBranch, repo)
	}

	if repoName == nil {
		panic("repository name is mandatory for branch config")
	}

	return &ghbranch.Config{
		ValueGenerator: valGen,
		Identifier:     identifier,
		Repository:     repoName,
		Branch:         &name,
		SourceBranch:   sourceBranch,
		SourceSha:      branchConfig.SourceSha,
	}
}

func MapToDefaultBranchRes(
	branchConfig *GhDefaultBranchConfig,
	valGen tfsig.ValueGenerator,
	repo *GhRepoConfig,
	repoTfId string,
	links ...MapperLink,
) *ghbranchdefault.Config {
	if branchConfig == nil {
		return nil
	}

	var repository *string

	branch := branchConfig.Name

	if repo != nil {
		repository = repo.Name
	}

	for _, link := range links {
		if link == LinkToRepository {
			// /!\ a branch can't be configured if repository doesn't exist
			// => Add an explicit dependency by using "github_repository.${repoTfId}.name"
			tmp := fmt.Sprintf("github_repository.%s.name", repoTfId)
			repository = &tmp
		} else if link == LinkToBranch {
			if branch != nil && repo.Branches != nil {
				if _, branchConfigExists := (*repo.Branches)[*branch]; branchConfigExists {
					// /!\ default branch can't be configured if related branch doesn't exist
					// => Add an explicit dependency by using "github_branch.${repoTfId}-${branchId}.branch"
					tmp := fmt.Sprintf("github_branch.%s-%s.branch", repoTfId, tfsig.ToTerraformIdentifier(*branch))
					branch = &tmp
				}
			}
		}
	}

	if repository == nil {
		panic("repository is mandatory for default branch config")
	}

	return &ghbranchdefault.Config{
		ValueGenerator: valGen,
		Identifier:     repoTfId,
		Repository:     repository,
		Branch:         branch,
	}
}

func MapDefaultBranchToBranchProtectionRes(
	branchConfig *GhDefaultBranchConfig,
	valGen tfsig.ValueGenerator,
	repo *GhRepoConfig,
	repoTfId string,
	links ...MapperLink,
) *ghbranchprotect.Config {
	if branchConfig == nil {
		return nil
	}

	res := MapBranchToBranchProtectionRes(
		branchConfig.Name,
		branchConfig.Protection,
		valGen,
		repo,
		repoTfId,
		links...,
	)
	if res != nil {
		res.Identifier = repoTfId + "-" + DefaultBranchIdentifier

		for _, v := range links {
			if v == LinkToBranch {
				tmp := fmt.Sprintf("github_branch_default.%s.branch", repoTfId)
				res.Pattern = &tmp

				break
			}
		}
	}

	return res
}

func MapBranchToBranchProtectionRes(
	pattern *string,
	protection *BaseGhBranchProtectionConfig,
	valGen tfsig.ValueGenerator,
	repo *GhRepoConfig,
	repoTfId string,
	links ...MapperLink,
) *ghbranchprotect.Config {
	if protection == nil {
		return nil
	}

	wrapper := &GhBranchProtectionConfig{
		Pattern:                      pattern,
		Forbid:                       &falseString,
		BaseGhBranchProtectionConfig: *newBaseGhBranchProtectionConfigFromBranchProtection(protection),
	}

	return MapToBranchProtectionRes(wrapper, valGen, repo, repoTfId, links...)
}

func MapToBranchProtectionRes(
	branchProtectionConfig *GhBranchProtectionConfig,
	valGen tfsig.ValueGenerator,
	repo *GhRepoConfig,
	repoTfId string,
	links ...MapperLink,
) *ghbranchprotect.Config {
	if branchProtectionConfig == nil {
		return nil
	}

	var (
		repoName             *string
		requiredStatusChecks *ghbranchprotect.RequiredStatusChecksConfig
		requiredPRReview     *ghbranchprotect.RequiredPRReviewConfig
		allowsForcePushes    *string
		pushRestrictions     *[]string
	)

	idEnd := "INVALID"
	pattern := branchProtectionConfig.Pattern

	if repo != nil {
		repoName = repo.Name
	}

	for _, v := range links {
		repoName, pattern = mapBranchProtectionResLink(v, repoTfId, repoName, pattern, repo)
	}

	if repoName == nil {
		panic("repository name is mandatory for branch protection config")
	}

	if branchProtectionConfig.Pattern != nil {
		idEnd = tfsig.ToTerraformIdentifier(patternToIdReplacer.Replace(*branchProtectionConfig.Pattern))
	}

	if branchProtectionConfig.StatusChecks != nil {
		requiredStatusChecks = &ghbranchprotect.RequiredStatusChecksConfig{
			ValueGenerator: valGen,
			Strict:         branchProtectionConfig.StatusChecks.Strict,
			Contexts:       branchProtectionConfig.StatusChecks.Required,
		}
	}

	if branchProtectionConfig.PullRequestReviews != nil {
		requiredPRReview = &ghbranchprotect.RequiredPRReviewConfig{
			ValueGenerator:               valGen,
			DismissStaleReviews:          nil,
			RestrictDismissals:           nil,
			DismissalRestrictions:        nil,
			RequireCodeOwnerReviews:      branchProtectionConfig.PullRequestReviews.CodeownerApprovals,
			RequiredApprovingReviewCount: branchProtectionConfig.PullRequestReviews.ApprovalCount,
		}

		if branchProtectionConfig.PullRequestReviews.Dismissals != nil {
			requiredPRReview.DismissStaleReviews = branchProtectionConfig.PullRequestReviews.Dismissals.Staled
			requiredPRReview.RestrictDismissals = branchProtectionConfig.PullRequestReviews.Dismissals.Restrict
			requiredPRReview.DismissalRestrictions = branchProtectionConfig.PullRequestReviews.Dismissals.RestrictTo
		}
	}

	if branchProtectionConfig.Pushes != nil {
		allowsForcePushes = branchProtectionConfig.Pushes.AllowsForcePushes
		pushRestrictions = branchProtectionConfig.Pushes.RestrictTo
	}

	return &ghbranchprotect.Config{
		ValueGenerator:        valGen,
		Identifier:            repoTfId + "-" + idEnd,
		RepositoryId:          repoName,
		Pattern:               pattern,
		EnforceAdmins:         branchProtectionConfig.EnforceAdmins,
		AllowsDeletions:       branchProtectionConfig.AllowDeletion,
		AllowsForcePushes:     allowsForcePushes,
		PushRestrictions:      pushRestrictions,
		RequiredLinearHistory: branchProtectionConfig.RequireLinearHistory,
		RequireSignedCommits:  branchProtectionConfig.RequireSignedCommits,
		RequiredStatusChecks:  requiredStatusChecks,
		RequiredPRReview:      requiredPRReview,
	}
}

func mapBranchProtectionResLink(
	link MapperLink,
	repoTfId string,
	repoName *string,
	pattern *string,
	repo *GhRepoConfig,
) (*string, *string) {
	if link == LinkToRepository {
		// /!\ a branch protection can't be configured if repository doesn't exist
		// => Add an explicit dependency by using "github_repository.${repoTfId}.node_id"
		tmp := fmt.Sprintf("github_repository.%s.node_id", repoTfId)
		repoName = &tmp
	} else if link == LinkToBranch {
		if pattern != nil && repo != nil && repo.Branches != nil {
			if _, branchConfigExists := (*repo.Branches)[*pattern]; branchConfigExists {
				tmp := fmt.Sprintf("github_branch.%s-%s.branch", repoTfId, tfsig.ToTerraformIdentifier(*pattern))
				pattern = &tmp
			}
		}
	}

	return repoName, pattern
}

func newBaseGhBranchProtectionConfigFromBranchProtection(
	protection *BaseGhBranchProtectionConfig,
) *BaseGhBranchProtectionConfig {
	return &BaseGhBranchProtectionConfig{
		ConfigTemplates:      protection.ConfigTemplates,
		EnforceAdmins:        protection.EnforceAdmins,
		AllowDeletion:        protection.AllowDeletion,
		RequireLinearHistory: protection.RequireLinearHistory,
		RequireSignedCommits: protection.RequireSignedCommits,
		Pushes:               protection.Pushes,
		StatusChecks:         protection.StatusChecks,
		PullRequestReviews:   protection.PullRequestReviews,
	}
}

func mapPullRequest(
	repoConfig *GhRepoConfig,
) (*string, *string, *string, *string, *string, *string, *string, *string, *string) {
	// PullRequest
	var (
		allowMergeCommit         *string
		allowRebaseMerge         *string
		allowSquashMerge         *string
		allowAutoMerge           *string
		mergeCommitTitle         *string
		mergeCommitMessage       *string
		squashMergeCommitTitle   *string
		squashMergeCommitMessage *string
		deleteBranchOnMerge      *string
	)
	// var suggestUpdate *string

	if repoConfig.PullRequests != nil {
		if repoConfig.PullRequests.MergeStrategy != nil {
			allowMergeCommit = repoConfig.PullRequests.MergeStrategy.AllowMerge
			allowRebaseMerge = repoConfig.PullRequests.MergeStrategy.AllowRebase
			allowSquashMerge = repoConfig.PullRequests.MergeStrategy.AllowSquash
			allowAutoMerge = repoConfig.PullRequests.MergeStrategy.AllowAutoMerge
		}

		if repoConfig.PullRequests.MergeCommit != nil {
			mergeCommitTitle = repoConfig.PullRequests.MergeCommit.Title
			mergeCommitMessage = repoConfig.PullRequests.MergeCommit.Message
		}

		if repoConfig.PullRequests.SquashCommit != nil {
			squashMergeCommitTitle = repoConfig.PullRequests.SquashCommit.Title
			squashMergeCommitMessage = repoConfig.PullRequests.SquashCommit.Message
		}

		if repoConfig.PullRequests.Branch != nil {
			// suggestUpdate = c.PullRequests.Branch.SuggestUpdate
			deleteBranchOnMerge = repoConfig.PullRequests.Branch.DeleteOnMerge
		}
	}

	//nolint:lll
	return allowMergeCommit, allowRebaseMerge, allowSquashMerge, allowAutoMerge, mergeCommitTitle, mergeCommitMessage, squashMergeCommitTitle, squashMergeCommitMessage, deleteBranchOnMerge
}

//nolint:nonamedreturns // Easier to understand as is
func mapMiscellaneous(repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator) (
	topics *[]string,
	autoInit *string,
	archived *string,
	homepageUrl *string,
	// isTemplate *string,
	// gitignoreTemplate *string,
	// licenseTemplate *string,
	hasIssues *string,
	hasProjects *string,
	hasWiki *string,
	hasDownloads *string,
	page *ghrepository.PagesConfig,
	template *ghrepository.TemplateConfig,
) {
	if repoConfig.Miscellaneous != nil {
		topics = repoConfig.Miscellaneous.Topics
		autoInit = repoConfig.Miscellaneous.AutoInit
		archived = repoConfig.Miscellaneous.Archived
		homepageUrl = repoConfig.Miscellaneous.HomepageUrl
		// isTemplate = c.Miscellaneous.IsTemplate
		hasIssues = repoConfig.Miscellaneous.HasIssues
		hasProjects = repoConfig.Miscellaneous.HasProjects
		hasWiki = repoConfig.Miscellaneous.HasWiki
		hasDownloads = repoConfig.Miscellaneous.HasDownloads

		// if c.Miscellaneous.FileTemplates != nil {
		// gitignoreTemplate = c.Miscellaneous.FileTemplates.Gitignore
		// licenseTemplate = c.Miscellaneous.FileTemplates.License
		// }

		template = mapTemplate(repoConfig, valGen)
		page = mapPage(repoConfig, valGen)
	}

	return topics, autoInit, archived, homepageUrl, hasIssues, hasProjects, hasWiki, hasDownloads, page, template
}

func mapPage(repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator) *ghrepository.PagesConfig {
	if repoConfig.Miscellaneous.Pages != nil {
		var source *ghrepository.PagesSourceConfig
		if repoConfig.Miscellaneous.Pages.SourcePath != nil || repoConfig.Miscellaneous.Pages.SourceBranch != nil {
			source = &ghrepository.PagesSourceConfig{
				ValueGenerator: valGen,
				Branch:         repoConfig.Miscellaneous.Pages.SourceBranch,
				Path:           repoConfig.Miscellaneous.Pages.SourcePath,
			}
		}

		return &ghrepository.PagesConfig{Source: source}
	}

	return nil
}

func mapTemplate(repoConfig *GhRepoConfig, valGen tfsig.ValueGenerator) *ghrepository.TemplateConfig {
	if repoConfig.Miscellaneous.Template != nil {
		template := &ghrepository.TemplateConfig{
			ValueGenerator: valGen,
			// IncludeAllBranches: c.Miscellaneous.Template.FullClone,
			Owner:      nil,
			Repository: nil,
		}

		if repoConfig.Miscellaneous.Template.Source != nil {
			sources := strings.Split(*repoConfig.Miscellaneous.Template.Source, "/")
			if len(sources) == 2 { //nolint:gomnd // Doesn't make sense here to wrap 2
				template.Owner = &sources[0]
				template.Repository = &sources[1]
			}
		}

		return template
	}

	return nil
}

func applyBranchResLink(
	link MapperLink,
	repoTfId string,
	repoName *string,
	sourceBranch *string,
	repo *GhRepoConfig,
) (*string, *string) {
	if link == LinkToRepository {
		// /!\ a branch can't be configured if repository doesn't exist
		// => Add an explicit dependency by using "github_repository.${repoTfId}.name"
		tmp := fmt.Sprintf("github_repository.%s.name", repoTfId)
		repoName = &tmp
	} else if link == LinkToBranch && sourceBranch != nil && repo != nil {
		sourceBranch = linkBranch(sourceBranch, repo, repoTfId)
	}

	return repoName, sourceBranch
}

func linkBranch(sourceBranch *string, repo *GhRepoConfig, repoTfId string) *string {
	// /!\ a branch can't be configured if source_branch branch doesn't exist
	// => Add an explicit dependency by using "github_branch.${repoTfId}-${branchId}.branch"
	// Or Add an explicit dependency by using "github_branch_default.${repoTfId}.branch"
	tmp := *sourceBranch

	if repo.Branches != nil {
		if _, sourceBranchConfigExists := (*repo.Branches)[*sourceBranch]; sourceBranchConfigExists {
			tmp = fmt.Sprintf("github_branch.%s-%s.branch", repoTfId, tfsig.ToTerraformIdentifier(*sourceBranch))
		}
	}
	// Look for default branch only if not already updated
	if tmp == *sourceBranch && repo.DefaultBranch != nil && *repo.DefaultBranch.Name == *sourceBranch {
		tmp = fmt.Sprintf("github_branch_default.%s.branch", repoTfId)
	}

	return &tmp
}
