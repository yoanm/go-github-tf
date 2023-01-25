package core

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
)

func TestComputeRepoConfig(t *testing.T) {
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
	tplConfig := &TemplatesConfig{
		Repos: map[string]*GhRepoConfig{
			aTemplate: {Description: &description2, Miscellaneous: &GhRepoMiscellaneousConfig{Archived: &archived}},
			aTemplate2: {
				Branches: &GhBranchesConfig{
					"branch-name": &GhBranchConfig{
						BaseGhBranchConfig: BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					{&pattern, &forbid, BaseGhBranchProtectionConfig{EnforceAdmins: &enforceAdmins}},
				},
			},
		},
		BranchProtections: map[string]*GhBranchProtectionConfig{
			aTemplate: {
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					Pushes: &GhBranchProtectPushesConfig{
						AllowsForcePushes: &allowForcePushes,
					},
				},
			},
		},
		Branches: map[string]*GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
			aTemplate2: {
				BaseGhBranchConfig: BaseGhBranchConfig{
					Protection: &BaseGhBranchProtectionConfig{
						EnforceAdmins: &enforceAdmins,
					},
				},
			},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"repo template": {
			&GhRepoConfig{Name: &aName, Description: &description, ConfigTemplates: &[]string{aTemplate}},
			tplConfig,
			&GhRepoConfig{
				Name:        &aName,
				Description: &description,
				Miscellaneous: &GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			nil,
		},
		"default branch - branch template": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate2},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
							EnforceAdmins: &enforceAdmins,
						},
					},
				},
			},
			nil,
		},
		"default branch - branch protection template": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowsDeletion:  &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
							Pushes: &GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
			},
			nil,
		},
		"branch template": {
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						SourceBranch: &sourceBranch,
					},
				},
			},
			nil,
		},
		"branch - branch protection template": {
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							Protection: &BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
								AllowsDeletion:  &allowDeletions,
							},
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							Protection: &BaseGhBranchProtectionConfig{
								AllowsDeletion: &allowDeletions,
								Pushes: &GhBranchProtectPushesConfig{
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
			&GhRepoConfig{
				Name: &aName,
				BranchProtections: &GhBranchProtectionsConfig{
					{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowsDeletion:  &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
							Pushes: &GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
			},
			nil,
		},
		"repo + default branch (+protection) + branch (+protection) + branch protection templates": {
			&GhRepoConfig{
				Name:            &aName,
				ConfigTemplates: &[]string{aTemplate},
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate2},
						Protection: &BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowsDeletion:  &allowDeletions,
						},
					},
				},
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
							Protection: &BaseGhBranchProtectionConfig{
								ConfigTemplates: &[]string{aTemplate},
								AllowsDeletion:  &allowDeletions,
								Pushes: &GhBranchProtectPushesConfig{
									AllowsForcePushes: &allowForcePushes,
								},
							},
						},
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowsDeletion:  &allowDeletions,
						},
					},
				},
				Miscellaneous: &GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name:        &aName,
				Description: &description2,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
							EnforceAdmins:  &enforceAdmins,
							AllowsDeletion: &allowDeletions,
							Pushes: &GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
				Branches: &GhBranchesConfig{
					branchName: {
						SourceBranch: &sourceBranch,
						BaseGhBranchConfig: BaseGhBranchConfig{
							Protection: &BaseGhBranchProtectionConfig{
								AllowsDeletion: &allowDeletions,
								Pushes: &GhBranchProtectPushesConfig{
									AllowsForcePushes: &allowForcePushes,
								},
							},
						},
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
							Pushes: &GhBranchProtectPushesConfig{
								AllowsForcePushes: &allowForcePushes,
							},
						},
					},
				},
				Miscellaneous: &GhRepoMiscellaneousConfig{
					Archived: &archived,
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeRepoConfig(tc.value, tc.templates)
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
	tplConfig := &TemplatesConfig{
		Repos: map[string]*GhRepoConfig{
			aTemplate: {Description: &description2, Miscellaneous: &GhRepoMiscellaneousConfig{Archived: &archived}},
			aTemplate2: {
				Branches: &GhBranchesConfig{
					"branch-name": &GhBranchConfig{
						BaseGhBranchConfig: BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					{&pattern, &forbid, BaseGhBranchProtectionConfig{EnforceAdmins: &enforceAdmins}},
				},
			},
		},
		BranchProtections: map[string]*GhBranchProtectionConfig{
			aTemplate: {
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					Pushes: &GhBranchProtectPushesConfig{
						AllowsForcePushes: &allowForcePushes,
					},
				},
			},
		},
		Branches: map[string]*GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
			aTemplate2: {
				BaseGhBranchConfig: BaseGhBranchConfig{
					Protection: &BaseGhBranchProtectionConfig{
						EnforceAdmins: &enforceAdmins,
					},
				},
			},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"duplicated config due to templates": {
			&GhRepoConfig{
				Name:            &aName,
				ConfigTemplates: &[]string{aTemplate2},
				Branches: &GhBranchesConfig{
					"branch-name": &GhBranchConfig{
						SourceSha: &aSha,
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					"branch-name": &GhBranchConfig{
						SourceBranch: &sourceBranch,
						SourceSha:    &aSha,
					},
				},
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						Pattern: &pattern,
						Forbid:  &forbid,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
							EnforceAdmins:  &enforceAdmins,
						},
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeRepoConfig(tc.value, tc.templates)
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
	aName := "a_name"
	aTemplate := "a-template"
	branchName := "branch-name"
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"repo": {
			&GhRepoConfig{Name: &aName, ConfigTemplates: &[]string{aTemplate}},
			nil,
			nil,
			fmt.Errorf("unable to load repository template, no template available"),
		},
		"default branch - branch": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate},
					},
				},
			},
			nil,
			nil,
			fmt.Errorf("default branch: unable to load branch template, no template available"),
		},
		"default branch - branch protection": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
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
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							Protection: &BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				Name: &aName,
				BranchProtections: &GhBranchProtectionsConfig{
					{
						Pattern: &branchName,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
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
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeRepoConfig(tc.value, tc.templates)
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
	aName := "a_name"
	aTemplate := "a-template"
	emptyTplConfig := &TemplatesConfig{}
	branchName := "branch-name"
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"repo": {
			&GhRepoConfig{Name: &aName, ConfigTemplates: &[]string{aTemplate}},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown repository template a-template"),
		},
		"default branch - branch": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						ConfigTemplates: &[]string{aTemplate},
					},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("default branch: unknown branch template a-template"),
		},
		"default branch - branch protection": {
			&GhRepoConfig{
				Name: &aName,
				DefaultBranch: &GhDefaultBranchConfig{
					&branchName,
					BaseGhBranchConfig{
						Protection: &BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
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
			&GhRepoConfig{
				Name: &aName,
				Branches: &GhBranchesConfig{
					branchName: {
						BaseGhBranchConfig: BaseGhBranchConfig{
							Protection: &BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				Name: &aName,
				BranchProtections: &GhBranchProtectionsConfig{
					{
						Pattern: &branchName,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
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
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeRepoConfig(tc.value, tc.templates)
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
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no repo name": {
			&GhRepoConfig{ConfigTemplates: nil},
			nil,
			nil,
			fmt.Errorf("repository name is mandatory"),
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeRepoConfig(tc.value, tc.templates)
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
	aName := "a_name"
	aTemplate := "a-template"
	description := "my description"
	emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		Repos: map[string]*GhRepoConfig{
			aTemplate: {Description: &description},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&GhRepoConfig{ConfigTemplates: &[]string{aTemplate}},
			nil,
			nil,
			fmt.Errorf("unable to load repository template, no template available"),
		},
		"unknown template": {
			&GhRepoConfig{ConfigTemplates: &[]string{aTemplate}, Description: &description},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown repository template a-template"),
		},
		"no template provided": {
			&GhRepoConfig{Description: &description},
			emptyTplConfig,
			&GhRepoConfig{Description: &description},
			nil,
		},
		"base": {
			&GhRepoConfig{ConfigTemplates: &[]string{aTemplate}, Name: &aName},
			tplConfig,
			&GhRepoConfig{Description: &description, Name: &aName},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ApplyRepositoryTemplate(tc.value, tc.templates)
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
	aTemplate := "a-template"
	pattern := "my_pattern"
	allowDeletions := "true"
	emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		BranchProtections: map[string]*GhBranchProtectionConfig{
			aTemplate: {Pattern: &pattern},
		},
	}
	cases := map[string]struct {
		value     *GhBranchProtectionConfig
		templates *TemplatesConfig
		expected  *GhBranchProtectionConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			nil,
			nil,
			fmt.Errorf("unable to load branch protection template, no template available"),
		},
		"unknown template": {
			&GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown branch protection template a-template"),
		},
		"no template provided": {
			&GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					AllowsDeletion: &allowDeletions,
				},
			},
			emptyTplConfig,
			&GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					AllowsDeletion: &allowDeletions,
				},
			},
			nil,
		},
		"base": {
			&GhBranchProtectionConfig{
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					ConfigTemplates: &[]string{aTemplate},
					AllowsDeletion:  &allowDeletions,
				},
			},
			tplConfig,
			&GhBranchProtectionConfig{
				Pattern: &pattern,
				BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
					AllowsDeletion: &allowDeletions,
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ApplyBranchProtectionTemplate(tc.value, tc.templates)
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
	aTemplate := "a-template"
	pattern := "my_pattern"
	allowDeletions := "true"
	emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		BranchProtections: map[string]*GhBranchProtectionConfig{
			aTemplate: {Pattern: &pattern},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
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
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
						},
					},
				},
			},
			emptyTplConfig,
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
						},
					},
				},
			},
			nil,
		},
		"base": {
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							ConfigTemplates: &[]string{aTemplate},
							AllowsDeletion:  &allowDeletions,
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				BranchProtections: &GhBranchProtectionsConfig{
					&GhBranchProtectionConfig{
						Pattern: &pattern,
						BaseGhBranchProtectionConfig: BaseGhBranchProtectionConfig{
							AllowsDeletion: &allowDeletions,
						},
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				err := ApplyBranchProtectionsTemplate(tc.value, tc.templates)
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
	sourceSha := "source-sha"
	aTemplate := "a-template"
	sourceBranch := "source-branch"
	emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		Branches: map[string]*GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
		},
	}
	cases := map[string]struct {
		value     *GhBranchConfig
		templates *TemplatesConfig
		expected  *GhBranchConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&GhBranchConfig{
				nil,
				nil,
				BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			nil,
			nil,
			fmt.Errorf("unable to load branch template, no template available"),
		},
		"unknown template": {
			&GhBranchConfig{
				&sourceBranch,
				nil,
				BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			emptyTplConfig,
			nil,
			fmt.Errorf("unknown branch template a-template"),
		},
		"no template provided": {
			&GhBranchConfig{SourceBranch: &sourceBranch},
			emptyTplConfig,
			&GhBranchConfig{SourceBranch: &sourceBranch},
			nil,
		},
		"base": {
			&GhBranchConfig{
				nil,
				&sourceSha,
				BaseGhBranchConfig{
					ConfigTemplates: &[]string{aTemplate},
				},
			},
			tplConfig,
			&GhBranchConfig{
				SourceBranch: &sourceBranch,
				SourceSha:    &sourceSha,
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ApplyBranchTemplate(tc.value, tc.templates)
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
	branchName := "a-branch-name"
	sourceSha := "source-sha"
	aTemplate := "a-template"
	sourceBranch := "source-branch"
	emptyTplConfig := &TemplatesConfig{}
	tplConfig := &TemplatesConfig{
		Branches: map[string]*GhBranchConfig{
			aTemplate: {SourceBranch: &sourceBranch},
		},
	}
	cases := map[string]struct {
		value     *GhRepoConfig
		templates *TemplatesConfig
		expected  *GhRepoConfig
		error     error
	}{
		"nil": {
			nil,
			nil,
			nil,
			nil,
		},
		"no template available": {
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{
						nil,
						nil,
						BaseGhBranchConfig{
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
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{
						&sourceBranch,
						nil,
						BaseGhBranchConfig{
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
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{SourceBranch: &sourceBranch},
				},
			},
			emptyTplConfig,
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{SourceBranch: &sourceBranch},
				},
			},
			nil,
		},
		"base": {
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{
						nil,
						&sourceSha,
						BaseGhBranchConfig{
							ConfigTemplates: &[]string{aTemplate},
						},
					},
				},
			},
			tplConfig,
			&GhRepoConfig{
				Branches: &GhBranchesConfig{
					branchName: &GhBranchConfig{
						SourceBranch: &sourceBranch,
						SourceSha:    &sourceSha,
					},
				},
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				err := ApplyBranchesTemplate(tc.value, tc.templates)
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
