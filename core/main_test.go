package core

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
)

func EnsureConfigMatching(t *testing.T, expectedConf, actualConf interface{}, expectedError error, actualErr error) {
	if expectedError != nil {
		if actualErr == nil {
			t.Errorf("Expected an error but everything went well !")
		} else if actualErr.Error() != expectedError.Error() {
			t.Errorf("Error mismatch\n- expected\n+ actual\n\n%v", differ.LineDiff(expectedError.Error(), actualErr.Error()))
		}
	} else if expectedConf != nil {
		if actualErr != nil {
			t.Errorf("Error %s", actualErr)
		} else if diff := cmp.Diff(expectedConf, actualConf); diff != "" {
			t.Errorf("Config mismatch (-want +got):\n%s", diff)
		}
	} else {
		t.Errorf("No conf or error expected by the case !")
	}
}

func EnsureErrorMatching(t *testing.T, expectedErr error, actualErr error) {
	if actualErr == nil && expectedErr != nil {
		t.Errorf("Expected an error but everything went well !")
	} else if expectedErr != nil && actualErr != nil {
		if actualErr.Error() != expectedErr.Error() {
			t.Errorf("Error mismatch\n- expected\n+ actual\n\n%v", differ.LineDiff(expectedErr.Error(), actualErr.Error()))
		}
	} else if actualErr != nil {
		t.Errorf("Error %s", actualErr)
	}
}

func GetFullConfig(id int) *GhRepoConfig {
	bool1 := "false"
	bool2 := "true"
	if id%2 == 0 {
		bool1 = "true"
		bool2 = "false"
	}

	branchProtectionCount := 4
	approvalCount := (branchProtectionCount * id)
	// Repo
	name := fmt.Sprintf("repo%d", id)
	repoTemplate := fmt.Sprintf("a-repo-template%d", id)
	visibility := fmt.Sprintf("visibility%d", id)
	description := fmt.Sprintf("a description%d", id)
	// Repo->DefaultBranch
	defaultBranchBranchTemplate := fmt.Sprintf("default-branch-template%d", id)
	defaultBranchName := fmt.Sprintf("master%d", id)
	// Repo->DefaultBranch->Protection
	defaultBranchBranchProtectionTemplate := fmt.Sprintf("default-branch-branch-protection-template%d", id)
	defaultBranchBranchProtectionEnforceAdmins := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionAllowsDeletions := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionRequiredLinearHistory := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionRequireSignedCommits := fmt.Sprintf("%s", bool1)
	// Repo->DefaultBranch->Protection->Pushes
	defaultBranchBranchProtectionAllowsForcePushes := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionPushRestriction := fmt.Sprintf("default-branch-pushRestriction%d", id)
	// Repo->DefaultBranch->Protection->StatusChecks
	defaultBranchBranchProtectionStrict := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionContext := fmt.Sprintf("default-branch-context%d", id)
	// Repo->DefaultBranch->Protection->PullRequestReviews
	defaultBranchBranchProtectionBypasser := fmt.Sprintf("default-branch-bypasser%d", id)
	defaultBranchBranchProtectionRequiredApprovingReviewCount := fmt.Sprintf("%d", approvalCount%7)
	defaultBranchBranchProtectionRequireCodeOwnerReviews := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionResolvedConversations := fmt.Sprintf("%s", bool1)
	// Repo->DefaultBranch->Protection->PullRequestReviews->Dismissals
	defaultBranchBranchProtectionDismissStaleReviews := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionRestrictDismissal := fmt.Sprintf("%s", bool1)
	defaultBranchBranchProtectionDismissalRestriction := fmt.Sprintf("default-branch-dismissalRestriction%d", id)
	// Repo->Branches
	// Repo->Branches[0]
	branch1Name := fmt.Sprintf("feature/branch%d", id)
	branch1BranchTemplate := fmt.Sprintf("branch%d-branch-template%d", id, id)
	// branch1SourceBranch := fmt.Sprintf("branch%d-source-branch%d", id, id)
	// branch1SourceSha := fmt.Sprintf("branch%d-source-sha%d", id, id)
	// Repo->Branches[0]->Protection
	branch1BranchProtectionTemplate := fmt.Sprintf("branch%d-branch-protection-template%d", id, id)
	branch1BranchProtectionEnforceAdmins := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionAllowsDeletions := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionRequiredLinearHistory := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionRequireSignedCommits := fmt.Sprintf("%s", bool2)
	// Repo->Branches[0]->Protection->Pushes
	branch1BranchProtectionAllowsForcePushes := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionPushRestriction := fmt.Sprintf("branch%d-pushRestriction%d", id, id)
	// Repo->Branches[0]->Protection->StatusChecks
	branch1BranchProtectionStrict := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionContext := fmt.Sprintf("branch%d-context%d", id, id)
	// Repo->Branches[0]->Protection->PullRequestReviews
	branch1BranchProtectionBypasser := fmt.Sprintf("branch%d-bypasser%d", id, id)
	branch1BranchProtectionRequiredApprovingReviewCount := fmt.Sprintf("%d", (approvalCount+1)%7)
	branch1BranchProtectionRequireCodeOwnerReviews := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionResolvedConversations := fmt.Sprintf("%s", bool2)
	// Repo->Branches[0]->Protection->PullRequestReviews->Dismissals
	branch1BranchProtectionDismissStaleReviews := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionRestrictDismissal := fmt.Sprintf("%s", bool2)
	branch1BranchProtectionDismissalRestriction := fmt.Sprintf("branch%d-dismissalRestriction%d", id, id)
	// Repo->Branches[1]
	branch2Name := fmt.Sprintf("feature/branch%d", 1+id)
	branch2BranchTemplate := fmt.Sprintf("branch%d-branch-template%d", 1+id, id)
	branch2SourceBranch := fmt.Sprintf("branch%d-source-branch%d", 1+id, id)
	branch2SourceSha := fmt.Sprintf("branch%d-source-sha%d", 1+id, id)
	// Repo->Branches[1]->Protection
	branch2BranchProtectionTemplate := fmt.Sprintf("branch%d-branch-protection-template%d", 1+id, id)
	branch2BranchProtectionEnforceAdmins := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionAllowsDeletions := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionRequiredLinearHistory := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionRequireSignedCommits := fmt.Sprintf("%s", bool1)
	// Repo->Branches[1]->Protection->Pushes
	branch2BranchProtectionAllowsForcePushes := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionPushRestriction := fmt.Sprintf("branch%d-pushRestriction%d", 1+id, id)
	// Repo->Branches[1]->Protection->StatusChecks
	branch2BranchProtectionStrict := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionContext := fmt.Sprintf("branch%d-context%d", 1+id, id)
	// Repo->Branches[1]->Protection->PullRequestReviews
	branch2BranchProtectionBypasser := fmt.Sprintf("branch%d-bypasser%d", 1+id, id)
	branch2BranchProtectionRequiredApprovingReviewCount := fmt.Sprintf("%d", (approvalCount+2)%7)
	branch2BranchProtectionRequireCodeOwnerReviews := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionResolvedConversations := fmt.Sprintf("%s", bool1)
	// Repo->Branches[1]->Protection->PullRequestReviews->Dismissals
	branch2BranchProtectionDismissStaleReviews := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionRestrictDismissal := fmt.Sprintf("%s", bool1)
	branch2BranchProtectionDismissalRestriction := fmt.Sprintf("branch%d-dismissalRestriction%d", 1+id, id)
	// Repo->BranchProtections
	// Repo->BranchProtections[0]
	branchProtectionTemplate := fmt.Sprintf("branch-protection-template%d", id)
	pattern := fmt.Sprintf("a-pattern%d", id)
	forbid := fmt.Sprintf("%s", bool2)
	branchProtectionEnforceAdmins := fmt.Sprintf("%s", bool2)
	branchProtectionAllowsDeletions4 := fmt.Sprintf("%s", bool2)
	branchProtectionRequiredLinearHistory := fmt.Sprintf("%s", bool2)
	branchProtectionRequireSignedCommits := fmt.Sprintf("%s", bool2)
	// Repo->BranchProtections[0]->Pushes
	branchProtectionAllowsForcePushes := fmt.Sprintf("%s", bool2)
	branchProtectionPushRestriction := fmt.Sprintf("branch-protection-pushRestriction%d", id)
	// Repo->BranchProtections[0]->StatusChecks
	branchProtectionStrict := fmt.Sprintf("%s", bool2)
	branchProtectionContext := fmt.Sprintf("branch-protection-context%d", id)
	branchProtectionBypasser := fmt.Sprintf("branch-protection-bypasser%d", id)
	// Repo->BranchProtections[0]->PullRequestReviews
	branchProtectionRequiredApprovingReviewCount := fmt.Sprintf("%d", (approvalCount+3)%7)
	branchProtectionRequireCodeOwnerReviews := fmt.Sprintf("%s", bool2)
	branchProtectionResolvedConversations := fmt.Sprintf("%s", bool2)
	// Repo->BranchProtections[0]->PullRequestReviews->Dismissals
	branchProtectionDismissStaleReviews := fmt.Sprintf("%s", bool2)
	branchProtectionRestrictDismissal := fmt.Sprintf("%s", bool2)
	branchProtectionDismissalRestriction := fmt.Sprintf("branch-protection-dismissalRestriction%d", id)
	// Repo->PullRequests
	// Repo->PullRequests->MergeStrategy
	allowMergeCommit := fmt.Sprintf("%s", bool1)
	allowRebaseMerge := fmt.Sprintf("%s", bool1)
	allowSquashMerge := fmt.Sprintf("%s", bool1)
	allowAutoMerge := fmt.Sprintf("%s", bool1)
	// Repo->PullRequests->MergeCommit
	mergeCommitTitle := fmt.Sprintf("aMergeCommitTitle%d", id)
	mergeCommitMessage := fmt.Sprintf("aMergeCommitMessage%d", id)
	// Repo->PullRequests->SquashCommit
	squashMergeCommitTitle := fmt.Sprintf("aSquashMergeCommitTitle%d", id)
	squashMergeCommitMessage := fmt.Sprintf("aSquashMergeCommitMessage%d", id)
	// Repo->PullRequests->Branch
	suggestUpdate := fmt.Sprintf("%s", bool1)
	deleteBranchOnMerge := fmt.Sprintf("%s", bool1)
	// Repo->Security
	vulnerabilityAlerts := fmt.Sprintf("%s", bool2)
	// Repo->Misc
	archived := fmt.Sprintf("%s", bool1)
	topicCount := 2
	topic1 := fmt.Sprintf("topic%d", (topicCount * id))
	topic2 := fmt.Sprintf("topic%d", (topicCount*id)+1)
	autoInit := fmt.Sprintf("%s", bool1)
	isTemplate := fmt.Sprintf("%s", bool1)
	homepageUrl := fmt.Sprintf("http://localhost/%d", id)
	hasDownloads := fmt.Sprintf("%s", bool1)
	hasProjects := fmt.Sprintf("%s", bool1)
	hasWiki := fmt.Sprintf("%s", bool1)
	hasIssues := fmt.Sprintf("%s", bool1)
	// Repo->Misc->Pages
	domain := fmt.Sprintf("my.domain%d", id)
	branch := fmt.Sprintf("branch%d", id)
	path := fmt.Sprintf("path%d", id)
	// Repo->Misc->Template
	owner := fmt.Sprintf("owner%d", id)
	repository := fmt.Sprintf("repository%d", id)
	templateSource := owner + "/" + repository
	fullClone := fmt.Sprintf("%s", bool1)
	// Repo->Misc->FileTemplates
	gitignore := fmt.Sprintf("gitignore-tpl-name%d", id)
	license := fmt.Sprintf("license-tpl-name%d", id)
	// Repo->Terraform
	archiveOnDestroy := fmt.Sprintf("%s", bool2)
	ignoreVulnerabilityAlertsDuringRead := fmt.Sprintf("%s", bool2)

	return &GhRepoConfig{
		&name,
		&[]string{repoTemplate},
		&visibility,
		&description,
		&GhDefaultBranchConfig{
			&defaultBranchName,
			BaseGhBranchConfig{
				&[]string{defaultBranchBranchTemplate},
				&BaseGhBranchProtectionConfig{
					&[]string{defaultBranchBranchProtectionTemplate},
					&defaultBranchBranchProtectionEnforceAdmins,
					&defaultBranchBranchProtectionAllowsDeletions,
					&defaultBranchBranchProtectionRequiredLinearHistory,
					&defaultBranchBranchProtectionRequireSignedCommits,
					&GhBranchProtectPushesConfig{
						&defaultBranchBranchProtectionAllowsForcePushes,
						&[]string{defaultBranchBranchProtectionPushRestriction},
					},
					&GhBranchProtectStatusChecksConfig{
						&defaultBranchBranchProtectionStrict,
						&[]string{defaultBranchBranchProtectionContext},
					},
					&GhBranchProtectPRReviewConfig{
						&[]string{defaultBranchBranchProtectionBypasser},
						&defaultBranchBranchProtectionRequireCodeOwnerReviews,
						&defaultBranchBranchProtectionResolvedConversations,
						&defaultBranchBranchProtectionRequiredApprovingReviewCount,
						&GhBranchProtectPRReviewDismissalsConfig{
							&defaultBranchBranchProtectionDismissStaleReviews,
							&defaultBranchBranchProtectionRestrictDismissal,
							&[]string{defaultBranchBranchProtectionDismissalRestriction},
						},
					},
				},
			},
		},
		&GhBranchesConfig{
			branch1Name: {
				nil, // No base for first branch to be able to test that case when converting to terraform
				nil, // No base for first branch to be able to test that case when converting to terraform
				BaseGhBranchConfig{
					&[]string{branch1BranchTemplate},
					&BaseGhBranchProtectionConfig{
						&[]string{branch1BranchProtectionTemplate},
						&branch1BranchProtectionEnforceAdmins,
						&branch1BranchProtectionAllowsDeletions,
						&branch1BranchProtectionRequiredLinearHistory,
						&branch1BranchProtectionRequireSignedCommits,
						&GhBranchProtectPushesConfig{
							&branch1BranchProtectionAllowsForcePushes,
							&[]string{branch1BranchProtectionPushRestriction},
						},
						&GhBranchProtectStatusChecksConfig{
							&branch1BranchProtectionStrict,
							&[]string{branch1BranchProtectionContext},
						},
						&GhBranchProtectPRReviewConfig{
							&[]string{branch1BranchProtectionBypasser},
							&branch1BranchProtectionRequireCodeOwnerReviews,
							&branch1BranchProtectionResolvedConversations,
							&branch1BranchProtectionRequiredApprovingReviewCount,
							&GhBranchProtectPRReviewDismissalsConfig{
								&branch1BranchProtectionDismissStaleReviews,
								&branch1BranchProtectionRestrictDismissal,
								&[]string{branch1BranchProtectionDismissalRestriction},
							},
						},
					},
				},
			},
			branch2Name: {
				&branch2SourceBranch,
				&branch2SourceSha,
				BaseGhBranchConfig{
					&[]string{branch2BranchTemplate},
					&BaseGhBranchProtectionConfig{
						&[]string{branch2BranchProtectionTemplate},
						&branch2BranchProtectionEnforceAdmins,
						&branch2BranchProtectionAllowsDeletions,
						&branch2BranchProtectionRequiredLinearHistory,
						&branch2BranchProtectionRequireSignedCommits,
						&GhBranchProtectPushesConfig{
							&branch2BranchProtectionAllowsForcePushes,
							&[]string{branch2BranchProtectionPushRestriction},
						},
						&GhBranchProtectStatusChecksConfig{
							&branch2BranchProtectionStrict,
							&[]string{branch2BranchProtectionContext},
						},
						&GhBranchProtectPRReviewConfig{
							&[]string{branch2BranchProtectionBypasser},
							&branch2BranchProtectionRequireCodeOwnerReviews,
							&branch2BranchProtectionResolvedConversations,
							&branch2BranchProtectionRequiredApprovingReviewCount,
							&GhBranchProtectPRReviewDismissalsConfig{
								&branch2BranchProtectionDismissStaleReviews,
								&branch2BranchProtectionRestrictDismissal,
								&[]string{branch2BranchProtectionDismissalRestriction},
							},
						},
					},
				},
			},
		},
		&GhBranchProtectionsConfig{
			{
				&pattern,
				&forbid,
				BaseGhBranchProtectionConfig{
					&[]string{branchProtectionTemplate},
					&branchProtectionEnforceAdmins,
					&branchProtectionAllowsDeletions4,
					&branchProtectionRequiredLinearHistory,
					&branchProtectionRequireSignedCommits,
					&GhBranchProtectPushesConfig{
						&branchProtectionAllowsForcePushes,
						&[]string{branchProtectionPushRestriction},
					},
					&GhBranchProtectStatusChecksConfig{
						&branchProtectionStrict,
						&[]string{branchProtectionContext},
					},
					&GhBranchProtectPRReviewConfig{
						&[]string{branchProtectionBypasser},
						&branchProtectionRequireCodeOwnerReviews,
						&branchProtectionResolvedConversations,
						&branchProtectionRequiredApprovingReviewCount,
						&GhBranchProtectPRReviewDismissalsConfig{
							&branchProtectionDismissStaleReviews,
							&branchProtectionRestrictDismissal,
							&[]string{branchProtectionDismissalRestriction},
						},
					},
				},
			},
		},
		&GhRepoPullRequestConfig{
			&GhRepoPRMergeStrategyConfig{
				&allowMergeCommit,
				&allowRebaseMerge,
				&allowSquashMerge,
				&allowAutoMerge,
			},
			&GhRepoPRCommitConfig{
				&mergeCommitTitle,
				&mergeCommitMessage,
			},
			&GhRepoPRCommitConfig{
				&squashMergeCommitTitle,
				&squashMergeCommitMessage,
			},
			&GhRepoPRBranchConfig{
				&suggestUpdate,
				&deleteBranchOnMerge,
			},
		},
		&GhRepoSecurityConfig{&vulnerabilityAlerts},
		&GhRepoMiscellaneousConfig{
			&[]string{topic1, topic2},
			&autoInit,
			&archived,
			&isTemplate,
			&homepageUrl,
			&hasIssues,
			&hasWiki,
			&hasProjects,
			&hasDownloads,
			&GhRepoTemplateConfig{
				&templateSource,
				&fullClone,
			},
			&GhRepoPagesConfig{
				&domain,
				&branch,
				&path,
			},
			&GhRepoFileTemplatesConfig{
				&gitignore,
				&license,
			},
		},
		&GhRepoTerraformConfig{&archiveOnDestroy, &ignoreVulnerabilityAlertsDuringRead},
	}
}
