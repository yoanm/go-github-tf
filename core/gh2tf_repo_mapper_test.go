package core_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/go-gh2tf"
	"github.com/yoanm/go-gh2tf/ghbranch"
	"github.com/yoanm/go-gh2tf/ghbranchdefault"
	"github.com/yoanm/go-gh2tf/ghbranchprotect"
	"github.com/yoanm/go-gh2tf/ghrepository"
	"github.com/yoanm/go-github-tf/core"
	"github.com/yoanm/go-tfsig"
	"github.com/yoanm/go-tfsig/testutils"
)

func TestMapToRepositoryRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	tfRepoId := "an_id"
	cases := map[string]struct {
		value    *core.GhRepoConfig
		expected *ghrepository.Config
	}{
		"nil": {
			nil,
			nil,
		},
		"empty": {
			&core.GhRepoConfig{},
			&ghrepository.Config{
				ValueGenerator:           valGen,
				Identifier:               tfRepoId,
				Name:                     nil,
				Visibility:               nil,
				Archived:                 nil,
				Description:              nil,
				HasIssues:                nil,
				HasProjects:              nil,
				HasWiki:                  nil,
				HasDownloads:             nil,
				HomepageUrl:              nil,
				Topics:                   nil,
				VulnerabilityAlerts:      nil,
				AllowMergeCommit:         nil,
				AllowRebaseMerge:         nil,
				AllowSquashMerge:         nil,
				AllowAutoMerge:           nil,
				MergeCommitTitle:         nil,
				MergeCommitMessage:       nil,
				SquashMergeCommitTitle:   nil,
				SquashMergeCommitMessage: nil,
				DeleteBranchOnMerge:      nil,
				ArchiveOnDestroy:         nil,
				Page:                     nil,
				Template:                 nil,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				actual := core.MapToRepositoryRes(tc.value, valGen, tfRepoId)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	branchName := "a_branch_name"
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.name", tfId)
	defaultBranchLink := fmt.Sprintf("github_branch_default.%s.branch", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := core.GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &core.GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &core.GhBranchesConfig{
			branch1Name: &core.GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		name     string
		value    *core.GhBranchConfig
		repo     *core.GhRepoConfig
		links    []core.MapperLink
		expected *ghbranch.Config
	}{
		"nil": {
			branchName,
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			branchName,
			&core.GhBranchConfig{},
			&repoConfig,
			nil,
			&ghbranch.Config{
				ValueGenerator: valGen,
				Identifier:     tfId + "-" + branchName,
				Repository:     &repoName,
				Branch:         &branchName,
			},
		},
		"sourceBranch same as current branch": {
			branchName,
			&core.GhBranchConfig{
				SourceBranch: &branchName,
			},
			&repoConfig,
			nil,
			&ghbranch.Config{
				ValueGenerator: valGen,
				Identifier:     tfId + "-" + branchName,
				Repository:     &repoName,
				Branch:         &branchName,
			},
		},
		"with all links - default branch link": {
			branchName,
			&core.GhBranchConfig{
				SourceBranch: &defaultBranchName,
			},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranch.Config{
				ValueGenerator: valGen,
				Identifier:     tfId + "-" + branchName,
				Repository:     &repoLink,
				Branch:         &branchName,
				SourceBranch:   &defaultBranchLink,
			},
		},
		"with all links - branch link": {
			branchName,
			&core.GhBranchConfig{
				SourceBranch: &branch1Name,
			},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranch.Config{
				ValueGenerator: valGen,
				Identifier:     tfId + "-" + branchName,
				Repository:     &repoLink,
				Branch:         &branchName,
				SourceBranch:   &branch1Link,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				actual := core.MapToBranchRes(tc.name, tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchRes_panic(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			core.MapToBranchRes("a-name", &core.GhBranchConfig{}, valGen, nil, tfId)
		},
		"repository name is mandatory for branch config",
	)
}

func TestMapToDefaultBranchRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.name", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := core.GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &core.GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &core.GhBranchesConfig{
			branch1Name: &core.GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		value    *core.GhDefaultBranchConfig
		repo     *core.GhRepoConfig
		links    []core.MapperLink
		expected *ghbranchdefault.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&core.GhDefaultBranchConfig{
				Name: &defaultBranchName,
			},
			&repoConfig,
			nil,
			&ghbranchdefault.Config{
				ValueGenerator: valGen,
				Identifier:     tfId,
				Repository:     &repoName,
				Branch:         &defaultBranchName,
			},
		},
		"with all links": {
			&core.GhDefaultBranchConfig{
				Name: &branch1Name,
			},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranchdefault.Config{
				ValueGenerator: valGen,
				Identifier:     tfId,
				Repository:     &repoLink,
				Branch:         &branch1Link,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				actual := core.MapToDefaultBranchRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToDefaultBranchRes_panic(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			core.MapToDefaultBranchRes(&core.GhDefaultBranchConfig{}, valGen, nil, tfId)
		},
		"repository is mandatory for default branch config",
	)
}

func TestMapDefaultBranchToBranchProtectionRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	repoName := "a_repo"
	pattern := "a-pattern"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	defaultBranchLink := fmt.Sprintf("github_branch_default.%s.branch", tfId)
	repoConfig := core.GhRepoConfig{
		Name: &repoName,
	}
	cases := map[string]struct {
		value    *core.GhDefaultBranchConfig
		repo     *core.GhRepoConfig
		links    []core.MapperLink
		expected *ghbranchprotect.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&core.GhDefaultBranchConfig{},
			&repoConfig,
			nil,
			nil,
		},
		"without protection": {
			&core.GhDefaultBranchConfig{
				&pattern,
				core.BaseGhBranchConfig{nil, nil},
			},
			&repoConfig,
			nil,
			nil,
		},
		"with links": {
			&core.GhDefaultBranchConfig{
				&pattern,
				core.BaseGhBranchConfig{
					nil,
					&core.BaseGhBranchProtectionConfig{
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
						nil,
					},
				},
			},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranchprotect.Config{
				ValueGenerator:        valGen,
				Identifier:            tfId + "-default",
				RepositoryId:          &repoLink,
				Pattern:               &defaultBranchLink,
				EnforceAdmins:         nil,
				AllowsDeletions:       nil,
				AllowsForcePushes:     nil,
				PushRestrictions:      nil,
				RequiredLinearHistory: nil,
				RequireSignedCommits:  nil,
				RequiredStatusChecks:  nil,
				RequiredPRReview:      nil,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				actual := core.MapDefaultBranchToBranchProtectionRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapBranchToBranchProtectionRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	branchName := "a_branch_name"
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := core.GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &core.GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &core.GhBranchesConfig{
			branch1Name: &core.GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		name     string
		value    *core.BaseGhBranchProtectionConfig
		repo     *core.GhRepoConfig
		links    []core.MapperLink
		expected *ghbranchprotect.Config
	}{
		"nil": {
			branchName,
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			branchName,
			&core.BaseGhBranchProtectionConfig{},
			&repoConfig,
			nil,
			&ghbranchprotect.Config{
				ValueGenerator:        valGen,
				Identifier:            tfId + "-" + branchName,
				RepositoryId:          &repoName,
				Pattern:               &branchName,
				EnforceAdmins:         nil,
				AllowsDeletions:       nil,
				AllowsForcePushes:     nil,
				PushRestrictions:      nil,
				RequiredLinearHistory: nil,
				RequireSignedCommits:  nil,
				RequiredStatusChecks:  nil,
				RequiredPRReview:      nil,
			},
		},
		"with all links": {
			branch1Name,
			&core.BaseGhBranchProtectionConfig{},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranchprotect.Config{
				ValueGenerator:        valGen,
				Identifier:            tfId + "-" + branch1Name,
				RepositoryId:          &repoLink,
				Pattern:               &branch1Link,
				EnforceAdmins:         nil,
				AllowsDeletions:       nil,
				AllowsForcePushes:     nil,
				PushRestrictions:      nil,
				RequiredLinearHistory: nil,
				RequireSignedCommits:  nil,
				RequiredStatusChecks:  nil,
				RequiredPRReview:      nil,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				branchNameTmp := tc.name
				actual := core.MapBranchToBranchProtectionRes(&branchNameTmp, tc.value, valGen, tc.repo, tfId, tc.links...)

				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchProtectionRes(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	repoName := "a_repo"
	pattern := "a-pattern"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	branchLink := fmt.Sprintf("github_branch.%s-%s.branch", tfId, pattern)
	repoConfig := core.GhRepoConfig{
		Name: &repoName,
		Branches: &core.GhBranchesConfig{
			pattern: &core.GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		value    *core.GhBranchProtectionConfig
		repo     *core.GhRepoConfig
		links    []core.MapperLink
		expected *ghbranchprotect.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&core.GhBranchProtectionConfig{},
			&repoConfig,
			nil,
			&ghbranchprotect.Config{
				ValueGenerator:        valGen,
				Identifier:            tfId + "-INVALID",
				RepositoryId:          &repoName,
				Pattern:               nil,
				EnforceAdmins:         nil,
				AllowsDeletions:       nil,
				AllowsForcePushes:     nil,
				PushRestrictions:      nil,
				RequiredLinearHistory: nil,
				RequireSignedCommits:  nil,
				RequiredStatusChecks:  nil,
				RequiredPRReview:      nil,
			},
		},
		"with links": {
			&core.GhBranchProtectionConfig{
				&pattern,
				nil,
				core.BaseGhBranchProtectionConfig{
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
					nil,
				},
			},
			&repoConfig,
			[]core.MapperLink{core.LinkToRepository, core.LinkToBranch},
			&ghbranchprotect.Config{
				ValueGenerator:        valGen,
				Identifier:            tfId + "-" + pattern,
				RepositoryId:          &repoLink,
				Pattern:               &branchLink,
				EnforceAdmins:         nil,
				AllowsDeletions:       nil,
				AllowsForcePushes:     nil,
				PushRestrictions:      nil,
				RequiredLinearHistory: nil,
				RequireSignedCommits:  nil,
				RequiredStatusChecks:  nil,
				RequiredPRReview:      nil,
			},
		},
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()

				actual := core.MapToBranchProtectionRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchProtectionRes_panic(t *testing.T) {
	t.Parallel()

	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			core.MapToBranchProtectionRes(&core.GhBranchProtectionConfig{}, valGen, nil, tfId)
		},
		"repository name is mandatory for branch protection config",
	)
}
