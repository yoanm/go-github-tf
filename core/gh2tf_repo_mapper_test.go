package core

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/go-tfsig"
	"github.com/yoanm/go-tfsig/testutils"

	"github.com/yoanm/go-gh2tf"
	"github.com/yoanm/go-gh2tf/ghbranch"
	"github.com/yoanm/go-gh2tf/ghbranchdefault"
	"github.com/yoanm/go-gh2tf/ghbranchprotect"
	"github.com/yoanm/go-gh2tf/ghrepository"
)

func TestMapToRepositoryRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	tfRepoId := "an_id"
	cases := map[string]struct {
		value    *GhRepoConfig
		expected *ghrepository.Config
	}{
		"nil": {
			nil,
			nil,
		},
		"empty": {
			&GhRepoConfig{},
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
				actual := MapToRepositoryRes(tc.value, valGen, tfRepoId)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	branchName := "a_branch_name"
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.name", tfId)
	defaultBranchLink := fmt.Sprintf("github_branch_default.%s.branch", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &GhBranchesConfig{
			branch1Name: &GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		name     string
		value    *GhBranchConfig
		repo     *GhRepoConfig
		links    []MapperLink
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
			&GhBranchConfig{},
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
			&GhBranchConfig{
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
			&GhBranchConfig{
				SourceBranch: &defaultBranchName,
			},
			&repoConfig,
			[]MapperLink{LinkToRepository, LinkToBranch},
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
			&GhBranchConfig{
				SourceBranch: &branch1Name,
			},
			&repoConfig,
			[]MapperLink{LinkToRepository, LinkToBranch},
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
				actual := MapToBranchRes(tc.name, tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchRes_panic(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			MapToBranchRes("a-name", &GhBranchConfig{}, valGen, nil, tfId)
		},
		"repository name is mandatory for branch config",
	)
}

func TestMapToDefaultBranchRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.name", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &GhBranchesConfig{
			branch1Name: &GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		value    *GhDefaultBranchConfig
		repo     *GhRepoConfig
		links    []MapperLink
		expected *ghbranchdefault.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&GhDefaultBranchConfig{
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
			&GhDefaultBranchConfig{
				Name: &branch1Name,
			},
			&repoConfig,
			[]MapperLink{LinkToRepository, LinkToBranch},
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
				actual := MapToDefaultBranchRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToDefaultBranchRes_panic(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			MapToDefaultBranchRes(&GhDefaultBranchConfig{}, valGen, nil, tfId)
		},
		"repository is mandatory for default branch config",
	)
}

func TestMapDefaultBranchToBranchProtectionRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	repoName := "a_repo"
	pattern := "a-pattern"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	defaultBranchLink := fmt.Sprintf("github_branch_default.%s.branch", tfId)
	repoConfig := GhRepoConfig{
		Name: &repoName,
	}
	cases := map[string]struct {
		value    *GhDefaultBranchConfig
		repo     *GhRepoConfig
		links    []MapperLink
		expected *ghbranchprotect.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&GhDefaultBranchConfig{},
			&repoConfig,
			nil,
			nil,
		},
		"without protection": {
			&GhDefaultBranchConfig{
				&pattern,
				BaseGhBranchConfig{nil, nil},
			},
			&repoConfig,
			nil,
			nil,
		},
		"with links": {
			&GhDefaultBranchConfig{
				&pattern,
				BaseGhBranchConfig{
					nil,
					&BaseGhBranchProtectionConfig{
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
			[]MapperLink{LinkToRepository, LinkToBranch},
			&ghbranchprotect.Config{
				valGen,
				tfId + "-default",
				&repoLink,
				&defaultBranchLink,
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
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}
	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := MapDefaultBranchToBranchProtectionRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapBranchToBranchProtectionRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	branchName := "a_branch_name"
	defaultBranchName := "default_branch_name"
	branch1Name := "branch1_name"
	repoName := "a_repo"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	branch1Link := fmt.Sprintf("github_branch.%s-%s.branch", tfId, branch1Name)
	repoConfig := GhRepoConfig{
		Name: &repoName,
		DefaultBranch: &GhDefaultBranchConfig{
			Name: &defaultBranchName,
		},
		Branches: &GhBranchesConfig{
			branch1Name: &GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		name     string
		value    *GhBranchConfig
		repo     *GhRepoConfig
		links    []MapperLink
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
			&GhBranchConfig{},
			&repoConfig,
			nil,
			nil,
		},
		"without protection": {
			branchName,
			&GhBranchConfig{
				SourceBranch: &branchName,
			},
			&repoConfig,
			nil,
			nil,
		},
		"with all links": {
			branch1Name,
			&GhBranchConfig{
				SourceBranch: &branch1Name,
				BaseGhBranchConfig: BaseGhBranchConfig{
					Protection: &BaseGhBranchProtectionConfig{},
				},
			},
			&repoConfig,
			[]MapperLink{LinkToRepository, LinkToBranch},
			&ghbranchprotect.Config{
				valGen,
				tfId + "-" + branch1Name,
				&repoLink,
				&branch1Link,
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
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}
	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := MapBranchToBranchProtectionRes(tc.name, tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchProtectionRes(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	repoName := "a_repo"
	pattern := "a-pattern"
	tfId := "an_id"
	repoLink := fmt.Sprintf("github_repository.%s.node_id", tfId)
	branchLink := fmt.Sprintf("github_branch.%s-%s.branch", tfId, pattern)
	repoConfig := GhRepoConfig{
		Name: &repoName,
		Branches: &GhBranchesConfig{
			pattern: &GhBranchConfig{},
		},
	}
	cases := map[string]struct {
		value    *GhBranchProtectionConfig
		repo     *GhRepoConfig
		links    []MapperLink
		expected *ghbranchprotect.Config
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"empty": {
			&GhBranchProtectionConfig{},
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
			&GhBranchProtectionConfig{
				&pattern,
				nil,
				BaseGhBranchProtectionConfig{
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
			[]MapperLink{LinkToRepository, LinkToBranch},
			&ghbranchprotect.Config{
				valGen,
				tfId + "-" + pattern,
				&repoLink,
				&branchLink,
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
	}

	diffOpts := []cmp.Option{cmp.AllowUnexported(tfsig.ValueGenerator{}), cmp.AllowUnexported(tfsig.IdentTokenMatcher{})}
	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := MapToBranchProtectionRes(tc.value, valGen, tc.repo, tfId, tc.links...)
				if diff := cmp.Diff(tc.expected, actual, diffOpts...); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestMapToBranchProtectionRes_panic(t *testing.T) {
	valGen := gh2tf.NewValueGenerator()
	tfId := "an_id"

	testutils.ExpectPanic(
		t,
		"Basic",
		func() {
			MapToBranchProtectionRes(&GhBranchProtectionConfig{}, valGen, nil, tfId)
		},
		"repository name is mandatory for branch protection config",
	)
}
