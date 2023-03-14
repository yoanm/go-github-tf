package core_test

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/github-tf/core"
)

func TestComputeRepoConfig(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	pattern := "my_pattern"
	allowDeletions := "true"
	archived := "true"
	forbid := "true"
	allowForcePushes := "true"
	enforceAdmins := "false"
	sourceBranch := "a-source-branch"
	aTemplate := "a-template"
	aTemplate2 := "a-template2"
	description := "my description"
	description2 := "my description2"
	branchName := "branch-name"
	tplConfig := &core.TemplatesConfig{
		Repos: map[string]*core.GhRepoConfig{
			aTemplate: {Description: &description2, Miscellaneous: &core.GhRepoMiscellaneousConfig{Archived: &archived}},
			aTemplate2: {
				Branches: &core.GhBranchesConfig{
					"branch-name": &core.GhBranchConfig{
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					{&pattern, &forbid, core.BaseGhBranchProtectionConfig{EnforceAdmins: &enforceAdmins}},
				},
			},
		},
		BranchProtections: map[string]*core.GhBranchProtectionConfig{
			aTemplate: {
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					Pushes: &core.GhBranchProtectPushesConfig{
						AllowsForcePushes: &allowForcePushes,
					},
				},
			},
		},
		Branches: map[string]*core.GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
			aTemplate2: {
				BaseGhBranchConfig: core.BaseGhBranchConfig{
					Protection: &core.BaseGhBranchProtectionConfig{
						EnforceAdmins: &enforceAdmins,
					},
				},
			},
		},
	}
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"repo template": {
			&core.GhRepoConfig{Name: &aName, Description: &description, ConfigTemplates: &[]string{aTemplate}},
			tplConfig,
			&core.GhRepoConfig{
				Name:        &aName,
				Description: &description,
				Miscellaneous: &core.GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			nil,
		},
		"default branch - branch template": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate2},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							EnforceAdmins: &enforceAdmins,
						},
					},
				},
			},
			nil,
		},
		"default branch - branch protection template": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowDeletion:   &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
							Pushes: &core.GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
			},
			nil,
		},
		"branch template": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						SourceBranch: &sourceBranch,
					},
				},
			},
			nil,
		},
		"branch - branch protection template": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							Protection: &core.BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
								AllowDeletion:   &allowDeletions,
							},
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							Protection: &core.BaseGhBranchProtectionConfig{
								AllowDeletion: &allowDeletions,
								Pushes: &core.GhBranchProtectPushesConfig{
									AllowsForcePushes: &allowForcePushes,
								},
							},
						},
					},
				},
			},
			nil,
		},
		"branch protection template": {
			&core.GhRepoConfig{
				Name: &aName,
				BranchProtections: &core.GhBranchProtectionsConfig{
					{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowDeletion:   &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
							Pushes: &core.GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
			},
			nil,
		},
		"repo + default branch (+protection) + branch (+protection) + branch protection templates": {
			&core.GhRepoConfig{
				Name:            &aName,
				ConfigTemplates: &[]string{aTemplate},
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate2},
						Protection: &core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowDeletion:   &allowDeletions,
						},
					},
				},
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
							Protection: &core.BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
								AllowDeletion:   &allowDeletions,
								Pushes: &core.GhBranchProtectPushesConfig{
									AllowsForcePushes: &allowForcePushes,
								},
							},
						},
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowDeletion:   &allowDeletions,
						},
					},
				},
				Miscellaneous: &core.GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name:        &aName,
				Description: &description2,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							EnforceAdmins: &enforceAdmins,
							AllowDeletion: &allowDeletions,
							Pushes: &core.GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
				Branches: &core.GhBranchesConfig{
					branchName: {
						SourceBranch: &sourceBranch,
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							Protection: &core.BaseGhBranchProtectionConfig{
								AllowDeletion: &allowDeletions,
								Pushes: &core.GhBranchProtectPushesConfig{
									AllowsForcePushes: &allowForcePushes,
								},
							},
						},
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
							Pushes: &core.GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
				Miscellaneous: &core.GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeRepoConfig(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestComputeRepoConfig_edgeCases(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	pattern := "my_pattern"
	allowDeletions := "true"
	archived := "true"
	forbid := "true"
	allowForcePushes := "true"
	enforceAdmins := "false"
	sourceBranch := "a-source-branch"
	aTemplate := "a-template"
	aTemplate2 := "a-template2"
	description2 := "my description2"
	aSha := "a-sha"
	tplConfig := &core.TemplatesConfig{
		Repos: map[string]*core.GhRepoConfig{
			aTemplate: {Description: &description2, Miscellaneous: &core.GhRepoMiscellaneousConfig{Archived: &archived}},
			aTemplate2: {
				Branches: &core.GhBranchesConfig{
					"branch-name": &core.GhBranchConfig{
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					{&pattern, &forbid, core.BaseGhBranchProtectionConfig{EnforceAdmins: &enforceAdmins}},
				},
			},
		},
		BranchProtections: map[string]*core.GhBranchProtectionConfig{
			aTemplate: {
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					Pushes: &core.GhBranchProtectPushesConfig{
						AllowsForcePushes: &allowForcePushes,
					},
				},
			},
		},
		Branches: map[string]*core.GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
			aTemplate2: {
				BaseGhBranchConfig: core.BaseGhBranchConfig{
					Protection: &core.BaseGhBranchProtectionConfig{
						EnforceAdmins: &enforceAdmins,
					},
				},
			},
		},
	}
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"duplicated config due to templates": {
			&core.GhRepoConfig{
				Name:            &aName,
				ConfigTemplates: &[]string{aTemplate2},
				Branches: &core.GhBranchesConfig{
					"branch-name": &core.GhBranchConfig{
						SourceSha: &aSha,
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					"branch-name": &core.GhBranchConfig{
						SourceBranch: &sourceBranch,
						SourceSha:    &aSha,
					},
				},
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						Pattern: &pattern,
						Forbid:  &forbid,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
							EnforceAdmins: &enforceAdmins,
						},
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeRepoConfig(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestComputeRepoConfig_noTemplateAvailable(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	aTemplate := "a-template"
	branchName := "branch-name"
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"repo": {
			&core.GhRepoConfig{Name: &aName, ConfigTemplates: &[]string{aTemplate}},
			nil,
			nil,
			fmt.Errorf("unable to load repository template, no template available"),
		},
		"default branch - branch": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("default branch: unable to load branch template, no template available"),
		},
		"default branch - branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("default branch: unable to load branch protection template, no template available"),
		},
		"branch": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("branch branch-name: unable to load branch template, no template available"),
		},
		"branch - branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							Protection: &core.BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
							},
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("branch branch-name: unable to load branch protection template, no template available"),
		},
		"branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				BranchProtections: &core.GhBranchProtectionsConfig{
					{
						Pattern: &branchName,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("branch protection #0: unable to load branch protection template, no template available"),
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeRepoConfig(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestComputeRepoConfig_unknownTemplate(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	aTemplate := "a-template"
	emptyTplConfig := &core.TemplatesConfig{}
	branchName := "branch-name"
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"repo": {
			&core.GhRepoConfig{Name: &aName, ConfigTemplates: &[]string{aTemplate}},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown repository template a-template"),
		},
		"default branch - branch": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("default branch: unknown branch template a-template"),
		},
		"default branch - branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				DefaultBranch: &core.GhDefaultBranchConfig{
					&branchName,
					core.BaseGhBranchConfig{
						Protection: &core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("default branch: unknown branch protection template a-template"),
		},
		"branch": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("branch branch-name: unknown branch template a-template"),
		},
		"branch - branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				Branches: &core.GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: core.BaseGhBranchConfig{
							Protection: &core.BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
							},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("branch branch-name: unknown branch protection template a-template"),
		},
		"branch protection": {
			&core.GhRepoConfig{
				Name: &aName,
				BranchProtections: &core.GhBranchProtectionsConfig{
					{
						Pattern: &branchName,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("branch protection #0: unknown branch protection template a-template"),
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeRepoConfig(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestComputeRepoConfig_validationError(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no repo name": {
			&core.GhRepoConfig{ConfigTemplates: nil},
			nil,
			nil,
			fmt.Errorf("repository name is mandatory"),
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeRepoConfig(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestApplyRepositoryTemplate(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	aTemplate := "a-template"
	description := "my description"
	emptyTplConfig := &core.TemplatesConfig{}
	tplConfig := &core.TemplatesConfig{
		Repos: map[string]*core.GhRepoConfig{
			aTemplate: {Description: &description},
		},
	}
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&core.GhRepoConfig{ConfigTemplates: &[]string{aTemplate}},
			nil,
			nil,
			fmt.Errorf("unable to load repository template, no template available"),
		},
		"unknown template": {
			&core.GhRepoConfig{ConfigTemplates: &[]string{aTemplate}, Description: &description},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown repository template a-template"),
		},
		"no template provided": {
			&core.GhRepoConfig{Description: &description},
			emptyTplConfig,
			&core.GhRepoConfig{Description: &description},
			nil,
		},
		"base": {
			&core.GhRepoConfig{ConfigTemplates: &[]string{aTemplate}, Name: &aName},
			tplConfig,
			&core.GhRepoConfig{Description: &description, Name: &aName},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ApplyRepositoryTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestApplyBranchProtectionTemplate(t *testing.T) {
	t.Parallel()

	aTemplate := "a-template"
	pattern := "my_pattern"
	allowDeletions := "true"
	emptyTplConfig := &core.TemplatesConfig{}
	tplConfig := &core.TemplatesConfig{
		BranchProtections: map[string]*core.GhBranchProtectionConfig{
			aTemplate: {Pattern: &pattern},
		},
	}
	cases := map[string]struct {
		value     *core.GhBranchProtectionConfig
		templates *core.TemplatesConfig
		expected  *core.GhBranchProtectionConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&core.GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			nil,
			nil,
			fmt.Errorf("unable to load branch protection template, no template available"),
		},
		"unknown template": {
			&core.GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown branch protection template a-template"),
		},
		"no template provided": {
			&core.GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					AllowDeletion: &allowDeletions,
				},
			},
			emptyTplConfig,
			&core.GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					AllowDeletion: &allowDeletions,
				},
			},
			nil,
		},
		"base": {
			&core.GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
					AllowDeletion:   &allowDeletions,
				},
			},
			tplConfig,
			&core.GhBranchProtectionConfig{
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
					AllowDeletion: &allowDeletions,
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ApplyBranchProtectionTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestApplyBranchProtectionsTemplate(t *testing.T) {
	t.Parallel()

	aTemplate := "a-template"
	pattern := "my_pattern"
	allowDeletions := "true"
	emptyTplConfig := &core.TemplatesConfig{}
	tplConfig := &core.TemplatesConfig{
		BranchProtections: map[string]*core.GhBranchProtectionConfig{
			aTemplate: {Pattern: &pattern},
		},
	}
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("branch protection #0: unable to load branch protection template, no template available"),
		},
		"unknown template": {
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("branch protection #0: unknown branch protection template a-template"),
		},
		"no template provided": {
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
						},
					},
				},
			},
			emptyTplConfig,
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
						},
					},
				},
			},
			nil,
		},
		"base": {
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowDeletion:   &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				BranchProtections: &core.GhBranchProtectionsConfig{
					&core.GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: core.BaseGhBranchProtectionConfig{
							AllowDeletion: &allowDeletions,
						},
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				err := core.ApplyBranchProtectionsTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestApplyBranchTemplate(t *testing.T) {
	t.Parallel()

	sourceSha := "source-sha"
	aTemplate := "a-template"
	sourceBranch := "source-branch"
	emptyTplConfig := &core.TemplatesConfig{}
	tplConfig := &core.TemplatesConfig{
		Branches: map[string]*core.GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
		},
	}
	cases := map[string]struct {
		value     *core.GhBranchConfig
		templates *core.TemplatesConfig
		expected  *core.GhBranchConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&core.GhBranchConfig{
				nil,
				nil,
				core.BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			nil,
			nil,
			fmt.Errorf("unable to load branch template, no template available"),
		},
		"unknown template": {
			&core.GhBranchConfig{
				&sourceBranch,
				nil,
				core.BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown branch template a-template"),
		},
		"no template provided": {
			&core.GhBranchConfig{SourceBranch: &sourceBranch},
			emptyTplConfig,
			&core.GhBranchConfig{SourceBranch: &sourceBranch},
			nil,
		},
		"base": {
			&core.GhBranchConfig{
				nil,
				&sourceSha,
				core.BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			tplConfig,
			&core.GhBranchConfig{
				SourceBranch: &sourceBranch,
				SourceSha:    &sourceSha,
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ApplyBranchTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, actual); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

func TestApplyBranchesTemplate(t *testing.T) {
	t.Parallel()

	branchName := "a-branch-name"
	sourceSha := "source-sha"
	aTemplate := "a-template"
	sourceBranch := "source-branch"
	emptyTplConfig := &core.TemplatesConfig{}
	tplConfig := &core.TemplatesConfig{
		Branches: map[string]*core.GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
		},
	}
	cases := map[string]struct {
		value     *core.GhRepoConfig
		templates *core.TemplatesConfig
		expected  *core.GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{
						nil,
						nil,
						core.BaseGhBranchConfig{
							&[]string{aTemplate},
							nil,
						},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("branch a-branch-name: unable to load branch template, no template available"),
		},
		"unknown template": {
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{
						&sourceBranch,
						nil,
						core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("branch a-branch-name: unknown branch template a-template"),
		},
		"no template provided": {
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{SourceBranch: &sourceBranch},
				},
			},
			emptyTplConfig,
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{SourceBranch: &sourceBranch},
				},
			},
			nil,
		},
		"base": {
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{
						nil,
						&sourceSha,
						core.BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			tplConfig,
			&core.GhRepoConfig{
				Branches: &core.GhBranchesConfig{
					branchName: &core.GhBranchConfig{
						SourceBranch: &sourceBranch,
						SourceSha:    &sourceSha,
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				err := core.ApplyBranchesTemplate(tc.value, tc.templates)
				if tc.error != nil {
					if err == nil {
						t.Errorf("Case %q: expected an error but everything went well", tcname)
					} else if err.Error() != tc.error.Error() {
						t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, differ.LineDiff(tc.error.Error(), err.Error()))
					}
				} else if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}
