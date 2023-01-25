package core_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/github-tf/core"
)

func TestConfig_GetRepo(t *testing.T) {
	knownRepoName := "known-repo"
	unknownCase := core.NewConfig()
	knownCase := core.NewConfig()
	knownRepo := core.GhRepoConfig{Name: &knownRepoName}
	knownCase.Repos = append(knownCase.Repos, &knownRepo)
	cases := map[string]struct {
		value    *core.Config
		name     string
		expected *core.GhRepoConfig
	}{
		"unknown": {
			unknownCase,
			"unknown-repo",
			nil,
		},
		"known": {
			knownCase,
			knownRepoName,
			&knownRepo,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := tc.value.GetRepo(tc.name)
				if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				} else if fmt.Sprintf("%p", tc.expected) != fmt.Sprintf("%p", actual) {
					t.Errorf("Config mismatch want pointer to %p, got pointer to %p", tc.expected, actual)
				}
			},
		)
	}
}

func TestTemplatesConfig_GetRepo(t *testing.T) {
	knownTemplateName := "known-template"
	unknownCase := core.NewConfig().Templates
	knownCase := core.NewConfig().Templates
	knownTemplate := core.GhRepoConfig{Name: &knownTemplateName}
	knownCase.Repos[knownTemplateName] = &knownTemplate
	cases := map[string]struct {
		value    *core.TemplatesConfig
		name     string
		expected *core.GhRepoConfig
	}{
		"unknown": {
			unknownCase,
			"unknown-template",
			nil,
		},
		"known": {
			knownCase,
			knownTemplateName,
			&knownTemplate,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := tc.value.GetRepo(tc.name)
				if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				} else if fmt.Sprintf("%p", tc.expected) != fmt.Sprintf("%p", actual) {
					t.Errorf("Config mismatch want pointer to %p, got pointer to %p", tc.expected, actual)
				}
			},
		)
	}
}

func TestTemplatesConfig_GetBranch(t *testing.T) {
	knownTemplateName := "known-template"
	unknownCase := core.NewConfig().Templates
	knownCase := core.NewConfig().Templates
	knownTemplate := core.GhBranchConfig{}
	knownCase.Branches[knownTemplateName] = &knownTemplate
	cases := map[string]struct {
		value    *core.TemplatesConfig
		name     string
		expected *core.GhBranchConfig
	}{
		"unknown": {
			unknownCase,
			"unknown-template",
			nil,
		},
		"known": {
			knownCase,
			knownTemplateName,
			&knownTemplate,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := tc.value.GetBranch(tc.name)
				if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				} else if fmt.Sprintf("%p", tc.expected) != fmt.Sprintf("%p", actual) {
					t.Errorf("Config mismatch want pointer to %p, got pointer to %p", tc.expected, actual)
				}
			},
		)
	}
}

func TestTemplatesConfig_GetBranchProtection(t *testing.T) {
	knownTemplateName := "known-template"
	unknownCase := core.NewConfig().Templates
	knownCase := core.NewConfig().Templates
	knownTemplate := core.GhBranchProtectionConfig{}
	knownCase.BranchProtections[knownTemplateName] = &knownTemplate
	cases := map[string]struct {
		value    *core.TemplatesConfig
		name     string
		expected *core.GhBranchProtectionConfig
	}{
		"unknown": {
			unknownCase,
			"unknown-template",
			nil,
		},
		"known": {
			knownCase,
			knownTemplateName,
			&knownTemplate,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual := tc.value.GetBranchProtection(tc.name)
				if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				} else if fmt.Sprintf("%p", tc.expected) != fmt.Sprintf("%p", actual) {
					t.Errorf("Config mismatch want pointer to %p, got pointer to %p", tc.expected, actual)
				}
			},
		)
	}
}
