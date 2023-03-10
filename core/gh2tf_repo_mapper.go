package core

import (
	"fmt"
	"strings"

	"github.com/yoanm/go-gh2tf/ghbranch"
	"github.com/yoanm/go-gh2tf/ghbranchdefault"
	"github.com/yoanm/go-gh2tf/ghbranchprotect"
	"github.com/yoanm/go-tfsig"

	"github.com/yoanm/go-gh2tf/ghrepository"
)

type MapperLink int

const (
	LinkToRepository MapperLink = iota
	LinkToBranch
)

const DefaultBranchIdentifier = "default"

func MapToRepositoryRes(c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) *ghrepository.Config {
	if c == nil {
		return nil
	}

	// Miscellaneous
	var topics *[]string
	var autoInit *string
	var archived *string
	var homepageUrl *string
	// var isTemplate *string
	// var gitignoreTemplate *string
	// var licenseTemplate *string
	var hasIssues *string
	var hasProjects *string
	var hasWiki *string
	var hasDownloads *string
	// Miscellaneous -> Page
	var page *ghrepository.PagesConfig
	// Miscellaneous -> Template
	var template *ghrepository.TemplateConfig

	if c.Miscellaneous != nil {
		topics = c.Miscellaneous.Topics
		autoInit = c.Miscellaneous.AutoInit
		archived = c.Miscellaneous.Archived
		homepageUrl = c.Miscellaneous.HomepageUrl
		// isTemplate = c.Miscellaneous.IsTemplate
		hasIssues = c.Miscellaneous.HasIssues
		hasProjects = c.Miscellaneous.HasProjects
		hasWiki = c.Miscellaneous.HasWiki
		hasDownloads = c.Miscellaneous.HasDownloads

		if c.Miscellaneous.Template != nil {
			template = &ghrepository.TemplateConfig{
				ValueGenerator: valGen,
				// IncludeAllBranches: c.Miscellaneous.Template.FullClone,
				Owner:      nil,
				Repository: nil,
			}
			if c.Miscellaneous.Template.Source != nil {
				sources := strings.Split(*c.Miscellaneous.Template.Source, "/")
				if len(sources) == 2 {
					template.Owner = &sources[0]
					template.Repository = &sources[1]
				}
			}
		}
		if c.Miscellaneous.Pages != nil {
			var Source *ghrepository.PagesSourceConfig
			if c.Miscellaneous.Pages.SourcePath != nil || c.Miscellaneous.Pages.SourceBranch != nil {
				Source = &ghrepository.PagesSourceConfig{
					ValueGenerator: valGen,
					Branch:         c.Miscellaneous.Pages.SourceBranch,
					Path:           c.Miscellaneous.Pages.SourcePath,
				}
			}
			page = &ghrepository.PagesConfig{Source: Source}
		}
		if c.Miscellaneous.FileTemplates != nil {
			// gitignoreTemplate = c.Miscellaneous.FileTemplates.Gitignore
			// licenseTemplate = c.Miscellaneous.FileTemplates.License
		}
	}
	// PullRequest
	var allowMergeCommit *string
	var allowRebaseMerge *string
	var allowSquashMerge *string
	var allowAutoMerge *string
	var mergeCommitTitle *string
	var mergeCommitMessage *string
	var squashMergeCommitTitle *string
	var squashMergeCommitMessage *string
	var deleteBranchOnMerge *string
	// var suggestUpdate *string
	if c.PullRequests != nil {
		if c.PullRequests.MergeStrategy != nil {
			allowMergeCommit = c.PullRequests.MergeStrategy.AllowMerge
			allowRebaseMerge = c.PullRequests.MergeStrategy.AllowRebase
			allowSquashMerge = c.PullRequests.MergeStrategy.AllowSquash
			allowAutoMerge = c.PullRequests.MergeStrategy.AllowAutoMerge
		}
		if c.PullRequests.MergeCommit != nil {
			mergeCommitTitle = c.PullRequests.MergeCommit.Title
			mergeCommitMessage = c.PullRequests.MergeCommit.Message
		}
		if c.PullRequests.SquashCommit != nil {
			squashMergeCommitTitle = c.PullRequests.SquashCommit.Title
			squashMergeCommitMessage = c.PullRequests.SquashCommit.Message
		}
		if c.PullRequests.Branch != nil {
			deleteBranchOnMerge = c.PullRequests.Branch.DeleteBranchOnMerge
			// suggestUpdate = c.PullRequests.Branch.SuggestUpdate
		}
	}
	// Security
	var vulnerabilityAlerts *string
	if c.Security != nil {
		vulnerabilityAlerts = c.Security.VulnerabilityAlerts
	}

	// Terraform
	var archiveOnDestroy *string
	if c.Terraform != nil {
		archiveOnDestroy = c.Terraform.ArchiveOnDestroy
	}

	return &ghrepository.Config{
		ValueGenerator: valGen,
		Identifier:     repoTfId,

		Name:                     c.Name,
		Visibility:               c.Visibility,
		Archived:                 archived,
		Description:              c.Description,
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

func MapToBranchRes(name string, c *GhBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranch.Config {
	if c == nil {
		return nil
	}

	identifier := fmt.Sprintf("%s-%s", repoTfId, tfsig.ToTerraformIdentifier(name))
	var sourceBranch *string
	if c.SourceBranch != nil && *c.SourceBranch != name {
		sourceBranch = c.SourceBranch
	}
	var repoName *string
	if repo != nil {
		repoName = repo.Name
	}
	for _, v := range links {
		if v == LinkToRepository {
			// /!\ a branch can't be configured if repository doesn't exist
			// => Add an explicit dependency by using "github_repository.${repoTfId}.name"
			tmp := fmt.Sprintf("github_repository.%s.name", repoTfId)
			repoName = &tmp
		} else if v == LinkToBranch {
			if sourceBranch != nil {
				if repo != nil {
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
					sourceBranch = &tmp
				}
			}
		}
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
		SourceSha:      c.SourceSha,
	}
}

func MapToDefaultBranchRes(c *GhDefaultBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchdefault.Config {
	if c == nil {
		return nil
	}
	branch := c.Name
	var repository *string
	if repo != nil {
		repository = repo.Name
	}
	for _, v := range links {
		if v == LinkToRepository {
			// /!\ a branch can't be configured if repository doesn't exist
			// => Add an explicit dependency by using "github_repository.${repoTfId}.name"
			tmp := fmt.Sprintf("github_repository.%s.name", repoTfId)
			repository = &tmp
		} else if v == LinkToBranch {
			if branch != nil {
				if repo.Branches != nil {
					if _, branchConfigExists := (*repo.Branches)[*branch]; branchConfigExists {
						// /!\ default branch can't be configured if related branch doesn't exist
						// => Add an explicit dependency by using "github_branch.${repoTfId}-${branchId}.branch"
						tmp := fmt.Sprintf("github_branch.%s-%s.branch", repoTfId, tfsig.ToTerraformIdentifier(*branch))
						branch = &tmp
					}
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

func MapDefaultBranchToBranchProtectionRes(c *GhDefaultBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config {
	if c == nil || c.Protection == nil {
		return nil
	}

	wrapper := &GhBranchProtectionConfig{
		Pattern: c.Name,
		BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
			ConfigTemplates:      c.Protection.ConfigTemplates,
			EnforceAdmins:        c.Protection.EnforceAdmins,
			AllowsDeletion:       c.Protection.AllowsDeletion,
			RequireLinearHistory: c.Protection.RequireLinearHistory,
			RequireSignedCommits: c.Protection.RequireSignedCommits,
			Pushes:               c.Protection.Pushes,
			StatusChecks:         c.Protection.StatusChecks,
			PullRequestReviews:   c.Protection.PullRequestReviews,
		},
	}

	res := MapToBranchProtectionRes(wrapper, valGen, repo, repoTfId, links...)
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

func MapBranchToBranchProtectionRes(name string, c *GhBranchConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config {
	if c == nil || c.Protection == nil {
		return nil
	}

	wrapper := &GhBranchProtectionConfig{
		Pattern: &name,
		BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
			ConfigTemplates:      c.Protection.ConfigTemplates,
			EnforceAdmins:        c.Protection.EnforceAdmins,
			AllowsDeletion:       c.Protection.AllowsDeletion,
			RequireLinearHistory: c.Protection.RequireLinearHistory,
			RequireSignedCommits: c.Protection.RequireSignedCommits,
			Pushes:               c.Protection.Pushes,
			StatusChecks:         c.Protection.StatusChecks,
			PullRequestReviews:   c.Protection.PullRequestReviews,
		},
	}

	return MapToBranchProtectionRes(wrapper, valGen, repo, repoTfId, links...)
}

func MapToBranchProtectionRes(c *GhBranchProtectionConfig, valGen tfsig.ValueGenerator, repo *GhRepoConfig, repoTfId string, links ...MapperLink) *ghbranchprotect.Config {
	if c == nil {
		return nil
	}

	pattern := c.Pattern
	var repoName *string
	if repo != nil {
		repoName = repo.Name
	}
	for _, v := range links {
		if v == LinkToRepository {
			// /!\ a branch protection can't be configured if repository doesn't exist
			// => Add an explicit dependency by using "github_repository.${repoTfId}.node_id"
			tmp := fmt.Sprintf("github_repository.%s.node_id", repoTfId)
			repoName = &tmp
		} else if v == LinkToBranch {
			if pattern != nil && repo != nil && repo.Branches != nil {
				if _, branchConfigExists := (*repo.Branches)[*pattern]; branchConfigExists {
					tmp := fmt.Sprintf("github_branch.%s-%s.branch", repoTfId, tfsig.ToTerraformIdentifier(*pattern))
					pattern = &tmp
				}
			}
		}
	}
	if repoName == nil {
		panic("repository name is mandatory for branch protection config")
	}
	idEnd := "INVALID"
	if c.Pattern != nil {
		idEnd = tfsig.ToTerraformIdentifier(*c.Pattern)
	}
	var requiredStatusChecks *ghbranchprotect.RequiredStatusChecksConfig
	if c.StatusChecks != nil {
		requiredStatusChecks = &ghbranchprotect.RequiredStatusChecksConfig{
			ValueGenerator: valGen,
			Strict:         c.StatusChecks.Strict,
			Contexts:       c.StatusChecks.Required,
		}
	}
	var requiredPRReview *ghbranchprotect.RequiredPRReviewConfig
	if c.PullRequestReviews != nil {
		requiredPRReview = &ghbranchprotect.RequiredPRReviewConfig{
			ValueGenerator:               valGen,
			DismissStaleReviews:          nil,
			RestrictDismissals:           nil,
			DismissalRestrictions:        nil,
			RequireCodeOwnerReviews:      c.PullRequestReviews.CodeownerApprovals,
			RequiredApprovingReviewCount: c.PullRequestReviews.ApprovalCount,
		}
		if c.PullRequestReviews.Dismissals != nil {
			requiredPRReview.DismissStaleReviews = c.PullRequestReviews.Dismissals.Staled
			requiredPRReview.RestrictDismissals = c.PullRequestReviews.Dismissals.Restrict
			requiredPRReview.DismissalRestrictions = c.PullRequestReviews.Dismissals.RestrictTo
		}
	}
	var allowsForcePushes *string
	var pushRestrictions *[]string
	if c.Pushes != nil {
		allowsForcePushes = c.Pushes.AllowsForcePushes
		pushRestrictions = c.Pushes.PushRestrictions
	}

	return &ghbranchprotect.Config{
		ValueGenerator:        valGen,
		Identifier:            repoTfId + "-" + idEnd,
		RepositoryId:          repoName,
		Pattern:               pattern,
		EnforceAdmins:         c.EnforceAdmins,
		AllowsDeletions:       c.AllowsDeletion,
		AllowsForcePushes:     allowsForcePushes,
		PushRestrictions:      pushRestrictions,
		RequiredLinearHistory: c.RequireLinearHistory,
		RequireSignedCommits:  c.RequireSignedCommits,
		RequiredStatusChecks:  requiredStatusChecks,
		RequiredPRReview:      requiredPRReview,
	}
}
