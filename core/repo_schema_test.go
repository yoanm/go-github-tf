package core

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// ==> Repo

func updateGhRepoConfigHelper(c *GhRepoConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Name, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.Visibility, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Description, stringToCopy, updatePtr)

	if c.DefaultBranch == nil {
		c.DefaultBranch = &GhDefaultBranchConfig{}
	}
	updateGhDefaultBranchConfigHelper(c.DefaultBranch, stringToCopy, newSliceToCopy, updatePtr)

	/*
		// Do nothing with Branches as it's complicated to test it with that way
		// See TestGhBranchesConfig_Merge instead
		if c.Branches == nil {
			c.Branches = &GhBranchesConfig{}
		}
		updateGhBranchesConfigHelper(c.Branches, stringToCopy, newSliceToCopy, updatePtr)
	*/

	/*
		// Do nothing with BranchProtections as it's complicated to test it with that way
		// See TestGhBranchProtectionsConfig_Merge instead
		if c.BranchProtections == nil {
			c.BranchProtections = &GhBranchProtectionsConfig{}
		}
		updateGhBranchProtectionsConfigHelper(c.BranchProtections, stringToCopy, newSliceToCopy, updatePtr)
	*/

	if c.PullRequests == nil {
		c.PullRequests = &GhRepoPullRequestConfig{}
	}
	updateGhRepoPullRequestConfigHelper(c.PullRequests, stringToCopy, newSliceToCopy, updatePtr)

	if c.Security == nil {
		c.Security = &GhRepoSecurityConfig{}
	}
	updateGhRepoSecurityConfigHelper(c.Security, stringToCopy, newSliceToCopy, updatePtr)

	if c.Miscellaneous == nil {
		c.Miscellaneous = &GhRepoMiscellaneousConfig{}
	}
	updateGhRepoMiscellaneousConfigHelper(c.Miscellaneous, stringToCopy, newSliceToCopy, updatePtr)

	if c.Terraform == nil {
		c.Terraform = &GhRepoTerraformConfig{}
	}
	updateGhRepoTerraformConfigHelper(c.Terraform, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoConfig{},
		func(to, from *GhRepoConfig) {
			to.Merge(from)
		},
		updateGhRepoConfigHelper,
	)
}

func TestGhRepoConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0)
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.DefaultBranch = nil
	toWithNilSlicesAndStruct.Branches = nil
	toWithNilSlicesAndStruct.BranchProtections = nil
	toWithNilSlicesAndStruct.PullRequests = nil
	toWithNilSlicesAndStruct.Security = nil
	toWithNilSlicesAndStruct.Miscellaneous = nil
	toWithNilSlicesAndStruct.Terraform = nil
	full1 := GetFullConfig(1)
	full2 := GetFullConfig(2)
	// manually generate result of full2 into full1
	fullMergeResult := GetFullConfig(2)
	*fullMergeResult.ConfigTemplates = append(
		*(full1.ConfigTemplates),
		*(full2.ConfigTemplates)...,
	)
	*fullMergeResult.DefaultBranch.ConfigTemplates = append(
		*(full1.DefaultBranch.ConfigTemplates),
		*(full2.DefaultBranch.ConfigTemplates)...,
	)
	*fullMergeResult.DefaultBranch.Protection.ConfigTemplates = append(
		*(full1.DefaultBranch.Protection.ConfigTemplates),
		*(full2.DefaultBranch.Protection.ConfigTemplates)...,
	)
	*fullMergeResult.DefaultBranch.Protection.Pushes.PushRestrictions = append(
		*(full1.DefaultBranch.Protection.Pushes.PushRestrictions),
		*(full2.DefaultBranch.Protection.Pushes.PushRestrictions)...,
	)
	*fullMergeResult.DefaultBranch.Protection.StatusChecks.Required = append(
		*(full1.DefaultBranch.Protection.StatusChecks.Required),
		*(full2.DefaultBranch.Protection.StatusChecks.Required)...,
	)
	*fullMergeResult.DefaultBranch.Protection.PullRequestReviews.BypasserList = append(
		*(full1.DefaultBranch.Protection.PullRequestReviews.BypasserList),
		*(full2.DefaultBranch.Protection.PullRequestReviews.BypasserList)...,
	)
	*fullMergeResult.DefaultBranch.Protection.PullRequestReviews.Dismissals.RestrictTo = append(
		*(full1.DefaultBranch.Protection.PullRequestReviews.Dismissals.RestrictTo),
		*(full2.DefaultBranch.Protection.PullRequestReviews.Dismissals.RestrictTo)...,
	)
	(*fullMergeResult.Branches)["feature/branch1"] = (*full1.Branches)["feature/branch1"]
	(*fullMergeResult.Branches)["feature/branch2"].SourceBranch = (*full1.Branches)["feature/branch2"].SourceBranch
	*(*fullMergeResult.Branches)["feature/branch2"].ConfigTemplates = append(
		*((*full1.Branches)["feature/branch2"].ConfigTemplates),
		*((*full2.Branches)["feature/branch2"].ConfigTemplates)...,
	)
	(*fullMergeResult.Branches)["feature/branch2"].SourceSha = (*full1.Branches)["feature/branch2"].SourceSha
	*(*fullMergeResult.Branches)["feature/branch2"].Protection.ConfigTemplates = append(
		*((*full1.Branches)["feature/branch2"].Protection.ConfigTemplates),
		*((*full2.Branches)["feature/branch2"].Protection.ConfigTemplates)...,
	)
	*(*fullMergeResult.Branches)["feature/branch2"].Protection.Pushes.PushRestrictions = append(
		*((*full1.Branches)["feature/branch2"].Protection.Pushes.PushRestrictions),
		*((*full2.Branches)["feature/branch2"].Protection.Pushes.PushRestrictions)...,
	)
	*(*fullMergeResult.Branches)["feature/branch2"].Protection.StatusChecks.Required = append(
		*((*full1.Branches)["feature/branch2"].Protection.StatusChecks.Required),
		*((*full2.Branches)["feature/branch2"].Protection.StatusChecks.Required)...,
	)
	*(*fullMergeResult.Branches)["feature/branch2"].Protection.PullRequestReviews.BypasserList = append(
		*((*full1.Branches)["feature/branch2"].Protection.PullRequestReviews.BypasserList),
		*((*full2.Branches)["feature/branch2"].Protection.PullRequestReviews.BypasserList)...,
	)
	*(*fullMergeResult.Branches)["feature/branch2"].Protection.PullRequestReviews.Dismissals.RestrictTo = append(
		*((*full1.Branches)["feature/branch2"].Protection.PullRequestReviews.Dismissals.RestrictTo),
		*((*full2.Branches)["feature/branch2"].Protection.PullRequestReviews.Dismissals.RestrictTo)...,
	)
	*fullMergeResult.BranchProtections = append(
		*(full1.BranchProtections),
		*(full2.BranchProtections)...,
	)
	*fullMergeResult.Miscellaneous.Topics = append(
		*(full1.Miscellaneous.Topics),
		*(full2.Miscellaneous.Topics)...,
	)
	cases := map[string]struct {
		value    *GhRepoConfig
		from     *GhRepoConfig
		expected *GhRepoConfig
	}{
		"full": {
			full1,
			full2,
			fullMergeResult,
		},
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1),
			GetFullConfig(1),
		},
		"from nil": {
			GetFullConfig(0),
			nil,
			GetFullConfig(0),
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->DefaultBranch

func updateGhDefaultBranchConfigHelper(c *GhDefaultBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Name, stringToCopy, updatePtr)

	updateBaseGhBranchConfigHelper(&c.BaseGhBranchConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhDefaultBranchConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhDefaultBranchConfig{},
		func(to, from *GhDefaultBranchConfig) {
			to.Merge(from)
		},
		updateGhDefaultBranchConfigHelper,
	)
}

func TestGhDefaultBranchConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0).DefaultBranch
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *GhDefaultBranchConfig
		from     *GhDefaultBranchConfig
		expected *GhDefaultBranchConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1).DefaultBranch,
			GetFullConfig(1).DefaultBranch,
		},
		"from nil": {
			GetFullConfig(0).DefaultBranch,
			nil,
			GetFullConfig(0).DefaultBranch,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches

func updateGhBranchesConfigHelper(c *GhBranchesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	v := &GhBranchConfig{}
	updateGhBranchConfigHelper(v, stringToCopy, newSliceToCopy, updatePtr)
	// If nil pointer or nil map
	if c == nil || *c == nil {
		*c = GhBranchesConfig{}
	}
	(*c)[*stringToCopy] = v
}

func TestGhBranchesConfig_Merge(t *testing.T) {
	// Init variables
	initialToStringValue := "initial_to_value"

	initialFromStringValue := "initial_from_value"

	// Init 'to' variables and create 'to'
	toSourceBranch := fmt.Sprintf("%s", initialToStringValue)
	toSourceSha := fmt.Sprintf("%s", initialToStringValue)
	toEnforceAdmins := fmt.Sprintf("%s", initialToStringValue)
	toAllowsDeletions := fmt.Sprintf("%s", initialToStringValue)
	toAllowsForcePushes := fmt.Sprintf("%s", initialToStringValue)
	toRequiredLinearHistory := fmt.Sprintf("%s", initialToStringValue)
	toRequireSignedCommits := fmt.Sprintf("%s", initialToStringValue)
	toStrict := fmt.Sprintf("%s", initialToStringValue)
	toDismissStaleReviews := fmt.Sprintf("%s", initialToStringValue)
	toRequireCodeOwnerReviews := fmt.Sprintf("%s", initialToStringValue)
	toResolvedConversations := fmt.Sprintf("%s", initialToStringValue)
	toRequiredApprovingReviewCount := fmt.Sprintf("%s", initialToStringValue)
	// Do nothing with PushRestrictions, see TestGhRepoConfig_Merge2
	// toPushRestrictions := append([]string{}, initialToSliceValue...)
	// Do nothing with Contexts, see TestGhRepoConfig_Merge2
	// toContext := append([]string{}, initialToSliceValue...)
	// Do nothing with DismissalRestrictions, see TestGhRepoConfig_Merge2
	// toDismissalRestrictions := append([]string{}, initialToSliceValue...)
	// Do nothing with ConfigTemplates, see TestGhRepoConfig_Merge2
	// toConfigTemplates := append([]string{}, initialFromSliceValue...)
	to := &GhBranchesConfig{
		"to_branch": &GhBranchConfig{
			SourceBranch: &toSourceBranch,
			SourceSha:    &toSourceSha,
			BaseGhBranchConfig: BaseGhBranchConfig{
				// ConfigTemplates: &toConfigTemplates,
				Protection: &BaseGhBranchProtectionConfig{
					// ConfigTemplates:       &toConfigTemplates,
					EnforceAdmins:  &toEnforceAdmins,
					AllowsDeletion: &toAllowsDeletions,
					Pushes: &GhBranchProtectPushesConfig{
						AllowsForcePushes: &toAllowsForcePushes,
						// PushRestrictions:      &toPushRestrictions,
					},
					RequireLinearHistory: &toRequiredLinearHistory,
					RequireSignedCommits: &toRequireSignedCommits,
					StatusChecks: &GhBranchProtectStatusChecksConfig{
						Strict: &toStrict,
						//	Contexts: &toContext,
					},
					PullRequestReviews: &GhBranchProtectPRReviewConfig{
						Dismissals: &GhBranchProtectPRReviewDismissalsConfig{
							Staled: &toDismissStaleReviews,
							// RestrictTo:  &toDismissalRestrictions,
						},
						CodeownerApprovals:    &toRequireCodeOwnerReviews,
						ResolvedConversations: &toResolvedConversations,
						ApprovalCount:         &toRequiredApprovingReviewCount,
					},
				},
			},
		},
	}

	// Init 'from' variables and create 'from'
	fromSourceBranch := fmt.Sprintf("%s", initialFromStringValue)
	fromSourceSha := fmt.Sprintf("%s", initialFromStringValue)
	fromEnforceAdmins := fmt.Sprintf("%s", initialFromStringValue)
	fromAllowsDeletions := fmt.Sprintf("%s", initialFromStringValue)
	fromAllowsForcePushes := fmt.Sprintf("%s", initialFromStringValue)
	fromRequiredLinearHistory := fmt.Sprintf("%s", initialFromStringValue)
	fromRequireSignedCommits := fmt.Sprintf("%s", initialFromStringValue)
	fromStrict := fmt.Sprintf("%s", initialFromStringValue)
	fromDismissStaleReviews := fmt.Sprintf("%s", initialFromStringValue)
	fromRequireCodeOwnerReviews := fmt.Sprintf("%s", initialFromStringValue)
	fromResolvedConversations := fmt.Sprintf("%s", initialToStringValue)
	fromRequiredApprovingReviewCount := fmt.Sprintf("%s", initialFromStringValue)
	// Do nothing with PushRestrictions, see TestGhRepoConfig_Merge2
	// fromPushRestrictions := append([]string{}, initialFromSliceValue...)
	// Do nothing with Contexts, see TestGhRepoConfig_Merge2
	// fromContext := append([]string{}, initialFromSliceValue...)
	// Do nothing with DismissalRestrictions, see TestGhRepoConfig_Merge2
	// fromDismissalRestrictions := append([]string{}, initialFromSliceValue...)
	// Do nothing with ConfigTemplates, see TestGhRepoConfig_Merge2
	// fromConfigTemplates := append([]string{}, initialFromSliceValue...)
	from := &GhBranchesConfig{
		"from_branch": &GhBranchConfig{
			SourceBranch: &fromSourceBranch,
			SourceSha:    &fromSourceSha,
			BaseGhBranchConfig: BaseGhBranchConfig{
				// ConfigTemplates: &fromConfigTemplates,
				Protection: &BaseGhBranchProtectionConfig{
					// ConfigTemplates:       &fromConfigTemplates,
					EnforceAdmins:  &fromEnforceAdmins,
					AllowsDeletion: &fromAllowsDeletions,
					Pushes: &GhBranchProtectPushesConfig{
						AllowsForcePushes: &fromAllowsForcePushes,
						// PushRestrictions:      &fromPushRestrictions,
					},
					RequireLinearHistory: &fromRequiredLinearHistory,
					RequireSignedCommits: &fromRequireSignedCommits,
					StatusChecks: &GhBranchProtectStatusChecksConfig{
						Strict: &fromStrict,
						// Contexts: &fromContext,
					},
					PullRequestReviews: &GhBranchProtectPRReviewConfig{
						Dismissals: &GhBranchProtectPRReviewDismissalsConfig{
							Staled: &fromDismissStaleReviews,
							// RestrictTo: &fromDismissalRestrictions,
						},
						CodeownerApprovals:    &fromRequireCodeOwnerReviews,
						ResolvedConversations: &fromResolvedConversations,
						ApprovalCount:         &fromRequiredApprovingReviewCount,
					},
				},
			},
		},
	}
	expected := &GhBranchesConfig{
		"to_branch":   (*to)["to_branch"],
		"from_branch": (*from)["from_branch"],
	}

	to.Merge(from)
	if diff := cmp.Diff(expected, to); diff != "" {
		t.Fatalf("Config mismatch (-want +got):\n%s", diff)
	}
}

func TestGhBranchesConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0).Branches
	(*toWithNilSlicesAndStruct)["feature/branch0"].ConfigTemplates = nil
	(*toWithNilSlicesAndStruct)["feature/branch0"].Protection = nil
	(*toWithNilSlicesAndStruct)["feature/branch1"].ConfigTemplates = nil
	(*toWithNilSlicesAndStruct)["feature/branch1"].Protection = nil

	expectedToWithNilSlicesAndStruct := GetFullConfig(1).Branches
	(*expectedToWithNilSlicesAndStruct)["feature/branch0"] = (*toWithNilSlicesAndStruct)["feature/branch0"]
	(*expectedToWithNilSlicesAndStruct)["feature/branch1"].SourceBranch = (*toWithNilSlicesAndStruct)["feature/branch1"].SourceBranch
	(*expectedToWithNilSlicesAndStruct)["feature/branch1"].SourceSha = (*toWithNilSlicesAndStruct)["feature/branch1"].SourceSha

	cases := map[string]struct {
		value    *GhBranchesConfig
		from     *GhBranchesConfig
		expected *GhBranchesConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1).Branches,
			expectedToWithNilSlicesAndStruct,
		},
		"from nil": {
			GetFullConfig(0).Branches,
			nil,
			GetFullConfig(0).Branches,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]

func updateGhBranchConfigHelper(c *GhBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.SourceBranch, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourceSha, stringToCopy, updatePtr)

	updateBaseGhBranchConfigHelper(&c.BaseGhBranchConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchConfig{},
		func(to, from *GhBranchConfig) {
			to.Merge(from)
		},
		updateGhBranchConfigHelper,
	)
}

func TestGhBranchConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).Branches)["feature/branch0"]
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *GhBranchConfig
		from     *GhBranchConfig
		expected *GhBranchConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).Branches)["feature/branch1"],
			(*GetFullConfig(1).Branches)["feature/branch1"],
		},
		"from nil": {
			(*GetFullConfig(0).Branches)["feature/branch1"],
			nil,
			(*GetFullConfig(0).Branches)["feature/branch1"],
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig}

func updateBaseGhBranchConfigHelper(c *BaseGhBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)

	if c.Protection == nil {
		c.Protection = &BaseGhBranchProtectionConfig{}
	}
	updateBaseGhBranchProtectionConfigHelper(c.Protection, stringToCopy, newSliceToCopy, updatePtr)
}

func TestBaseGhBranchConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		BaseGhBranchConfig{},
		func(to, from *BaseGhBranchConfig) {
			to.Merge(from)
		},
		updateBaseGhBranchConfigHelper,
	)
}

func TestBaseGhBranchConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).Branches)["feature/branch0"].BaseGhBranchConfig
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *BaseGhBranchConfig
		from     *BaseGhBranchConfig
		expected *BaseGhBranchConfig
	}{
		"to has nil slices and struct": {
			&toWithNilSlicesAndStruct,
			&(*GetFullConfig(1).Branches)["feature/branch1"].BaseGhBranchConfig,
			&(*GetFullConfig(1).Branches)["feature/branch1"].BaseGhBranchConfig,
		},
		"from nil": {
			&(*GetFullConfig(0).Branches)["feature/branch1"].BaseGhBranchConfig,
			nil,
			&(*GetFullConfig(0).Branches)["feature/branch1"].BaseGhBranchConfig,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->BranchProtections.
func updateGhBranchProtectionsConfigHelper(c *GhBranchProtectionsConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	v := &GhBranchProtectionConfig{}
	updateGhBranchProtectionConfigHelper(v, stringToCopy, newSliceToCopy, updatePtr)
	// If nil pointer or nil slice
	if c == nil || *c == nil {
		*c = GhBranchProtectionsConfig{}
	}
	*c = append(*c, v)
}

func TestGhBranchProtectionsConfig_Merge(t *testing.T) {
	// Init variables
	initialToStringValue := "initial_to_value"
	initialToSliceValue := []string{"initial_to_slice_value1", "initial_to_slice_value2"}

	initialFromStringValue := "initial_from_value"
	initialFromSliceValue := []string{"initial_from_slice_value1", "initial_from_slice_value2"}

	// Init 'to' variables and create 'to'
	toPattern := fmt.Sprintf("%s", initialToStringValue)
	toForbid := fmt.Sprintf("%s", initialToStringValue)
	toEnforceAdmins := fmt.Sprintf("%s", initialToStringValue)
	toAllowsDeletions := fmt.Sprintf("%s", initialToStringValue)
	toAllowsForcePushes := fmt.Sprintf("%s", initialToStringValue)
	toRequiredLinearHistory := fmt.Sprintf("%s", initialToStringValue)
	toRequireSignedCommits := fmt.Sprintf("%s", initialToStringValue)
	toStrict := fmt.Sprintf("%s", initialToStringValue)
	toDismissStaleReviews := fmt.Sprintf("%s", initialToStringValue)
	toRestrict := fmt.Sprintf("%s", initialToStringValue)
	toRequireCodeOwnerReviews := fmt.Sprintf("%s", initialToStringValue)
	toResolvedConversations := fmt.Sprintf("%s", initialToStringValue)
	toRequiredApprovingReviewCount := fmt.Sprintf("%s", initialToStringValue)
	toBypasserList := append([]string{}, initialToSliceValue...)
	toPushRestrictions := append([]string{}, initialToSliceValue...)
	toContext := append([]string{}, initialToSliceValue...)
	toDismissalRestrictions := append([]string{}, initialToSliceValue...)
	toConfigTemplates := append([]string{}, initialFromSliceValue...)
	to := &GhBranchProtectionsConfig{
		{
			&toPattern,
			&toForbid,
			BaseGhBranchProtectionConfig{
				&toConfigTemplates,
				&toEnforceAdmins,
				&toAllowsDeletions,
				&toRequiredLinearHistory,
				&toRequireSignedCommits,
				&GhBranchProtectPushesConfig{
					&toAllowsForcePushes,
					&toPushRestrictions,
				},
				&GhBranchProtectStatusChecksConfig{
					&toStrict,
					&toContext,
				},
				&GhBranchProtectPRReviewConfig{
					&toBypasserList,
					&toRequireCodeOwnerReviews,
					&toResolvedConversations,
					&toRequiredApprovingReviewCount,
					&GhBranchProtectPRReviewDismissalsConfig{
						&toDismissStaleReviews,
						&toRestrict,
						&toDismissalRestrictions,
					},
				},
			},
		},
	}

	// Init 'from' variables and create 'from'
	fromPattern := fmt.Sprintf("%s", initialFromStringValue)
	fromForbid := fmt.Sprintf("%s", initialFromStringValue)
	fromEnforceAdmins := fmt.Sprintf("%s", initialFromStringValue)
	fromAllowsDeletions := fmt.Sprintf("%s", initialFromStringValue)
	fromAllowsForcePushes := fmt.Sprintf("%s", initialFromStringValue)
	fromRequiredLinearHistory := fmt.Sprintf("%s", initialFromStringValue)
	fromRequireSignedCommits := fmt.Sprintf("%s", initialFromStringValue)
	fromStrict := fmt.Sprintf("%s", initialFromStringValue)
	fromDismissStaleReviews := fmt.Sprintf("%s", initialFromStringValue)
	fromRestrict := fmt.Sprintf("%s", initialFromStringValue)
	fromRequireCodeOwnerReviews := fmt.Sprintf("%s", initialFromStringValue)
	fromResolvedConversations := fmt.Sprintf("%s", initialFromStringValue)
	fromRequiredApprovingReviewCount := fmt.Sprintf("%s", initialFromStringValue)
	fromPushRestrictions := append([]string{}, initialFromSliceValue...)
	fromContext := append([]string{}, initialFromSliceValue...)
	fromBypasserList := append([]string{}, initialFromSliceValue...)
	fromDismissalRestrictions := append([]string{}, initialFromSliceValue...)
	fromConfigTemplates := append([]string{}, initialFromSliceValue...)
	from := &GhBranchProtectionsConfig{
		{
			&fromPattern,
			&fromForbid,
			BaseGhBranchProtectionConfig{
				&fromConfigTemplates,
				&fromEnforceAdmins,
				&fromAllowsDeletions,
				&fromRequiredLinearHistory,
				&fromRequireSignedCommits,
				&GhBranchProtectPushesConfig{
					&fromAllowsForcePushes,
					&fromPushRestrictions,
				},
				&GhBranchProtectStatusChecksConfig{
					&fromStrict,
					&fromContext,
				},
				&GhBranchProtectPRReviewConfig{
					&fromBypasserList,
					&fromRequireCodeOwnerReviews,
					&fromResolvedConversations,
					&fromRequiredApprovingReviewCount,
					&GhBranchProtectPRReviewDismissalsConfig{
						&fromDismissStaleReviews,
						&fromRestrict,
						&fromDismissalRestrictions,
					},
				},
			},
		},
	}

	to.Merge(from)
	// Ensure there is now two BranchProtectionConfigs
	if len(*to) != 2 {
		t.Fatalf("Expected 2 BranchProtectionConfig item, got %d", len(*to))
	}

	// Ensure added item to 'to' is not linked anymore to 'from' item
	fromItem := (*from)[0]
	toItem := (*to)[1]

	updateFn := updateGhBranchProtectionConfigHelper

	createCopyWithFn := func(stringVal *string, sliceVal *[]string) *GhBranchProtectionConfig {
		v := reflect.New(reflect.TypeOf(GhBranchProtectionConfig{})).Interface().(*GhBranchProtectionConfig)
		updateFn(v, stringVal, sliceVal, true)

		return v
	}

	// expectedInitialTo refers to the second item of 'to'
	// => following to.merge(from), it is equals to fromItem as merge append slice items
	expectedInitialTo := createCopyWithFn(&initialFromStringValue, &initialFromSliceValue)
	ensureNoOverflowBetweenToAndFrom(t, toItem, fromItem, createCopyWithFn, updateFn, expectedInitialTo)
}

func TestGhBranchProtectionsConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0).BranchProtections
	(*toWithNilSlicesAndStruct)[0].ConfigTemplates = nil
	(*toWithNilSlicesAndStruct)[0].Pushes = nil
	(*toWithNilSlicesAndStruct)[0].StatusChecks = nil
	(*toWithNilSlicesAndStruct)[0].PullRequestReviews = nil

	expectedToWithNilSlicesAndStruct := GetFullConfig(1).BranchProtections
	*expectedToWithNilSlicesAndStruct = append(*toWithNilSlicesAndStruct, *expectedToWithNilSlicesAndStruct...)

	cases := map[string]struct {
		value    *GhBranchProtectionsConfig
		from     *GhBranchProtectionsConfig
		expected *GhBranchProtectionsConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1).BranchProtections,
			expectedToWithNilSlicesAndStruct,
		},
		"from nil": {
			GetFullConfig(0).BranchProtections,
			nil,
			GetFullConfig(0).BranchProtections,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->BranchProtections[...]

func updateGhBranchProtectionConfigHelper(c *GhBranchProtectionConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Pattern, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Forbid, stringToCopy, updatePtr)

	updateBaseGhBranchProtectionConfigHelper(&c.BaseGhBranchProtectionConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectionConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchProtectionConfig{},
		func(to, from *GhBranchProtectionConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectionConfigHelper,
	)
}

func TestGhBranchProtectionConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0]
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Pushes = nil
	toWithNilSlicesAndStruct.StatusChecks = nil
	toWithNilSlicesAndStruct.PullRequestReviews = nil

	cases := map[string]struct {
		value    *GhBranchProtectionConfig
		from     *GhBranchProtectionConfig
		expected *GhBranchProtectionConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).BranchProtections)[0],
			(*GetFullConfig(1).BranchProtections)[0],
		},
		"from nil": {
			(*GetFullConfig(0).BranchProtections)[0],
			nil,
			(*GetFullConfig(0).BranchProtections)[0],
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig->Protection}
// and Repo->BranchProtection[...]{BaseGhBranchProtectionConfig}

func updateBaseGhBranchProtectionConfigHelper(c *BaseGhBranchProtectionConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.EnforceAdmins, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowsDeletion, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.RequireLinearHistory, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.RequireSignedCommits, stringToCopy, updatePtr)

	if c.Pushes == nil {
		c.Pushes = &GhBranchProtectPushesConfig{}
	}
	updateGhBranchProtectPushesConfigHelper(c.Pushes, stringToCopy, newSliceToCopy, updatePtr)

	if c.StatusChecks == nil {
		c.StatusChecks = &GhBranchProtectStatusChecksConfig{}
	}
	updateGhBranchProtectStatusChecksConfigHelper(c.StatusChecks, stringToCopy, newSliceToCopy, updatePtr)

	if c.PullRequestReviews == nil {
		c.PullRequestReviews = &GhBranchProtectPRReviewConfig{}
	}
	updateGhBranchProtectPRReviewConfigHelper(c.PullRequestReviews, stringToCopy, newSliceToCopy, updatePtr)
}

func TestBaseGhBranchProtectionConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		BaseGhBranchProtectionConfig{},
		func(to, from *BaseGhBranchProtectionConfig) {
			to.Merge(from)
		},
		updateBaseGhBranchProtectionConfigHelper,
	)
}

func TestBaseGhBranchProtectionConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].BaseGhBranchProtectionConfig
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Pushes = nil
	toWithNilSlicesAndStruct.StatusChecks = nil
	toWithNilSlicesAndStruct.PullRequestReviews = nil

	cases := map[string]struct {
		value    *BaseGhBranchProtectionConfig
		from     *BaseGhBranchProtectionConfig
		expected *BaseGhBranchProtectionConfig
	}{
		"to has nil slices and struct": {
			&toWithNilSlicesAndStruct,
			&(*GetFullConfig(1).BranchProtections)[0].BaseGhBranchProtectionConfig,
			&(*GetFullConfig(1).BranchProtections)[0].BaseGhBranchProtectionConfig,
		},
		"from nil": {
			&(*GetFullConfig(0).BranchProtections)[0].BaseGhBranchProtectionConfig,
			nil,
			&(*GetFullConfig(0).BranchProtections)[0].BaseGhBranchProtectionConfig,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig->Protection->Pushes}
// and Repo->BranchProtection[...]{BaseGhBranchProtectionConfig->Pushes}

func updateGhBranchProtectPushesConfigHelper(c *GhBranchProtectPushesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.AllowsForcePushes, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.PushRestrictions, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPushesConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchProtectPushesConfig{},
		func(to, from *GhBranchProtectPushesConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPushesConfigHelper,
	)
}

func TestGhBranchProtectPushesConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].Pushes
	toWithNilSlicesAndStruct.PushRestrictions = nil

	cases := map[string]struct {
		value    *GhBranchProtectPushesConfig
		from     *GhBranchProtectPushesConfig
		expected *GhBranchProtectPushesConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).BranchProtections)[0].Pushes,
			(*GetFullConfig(1).BranchProtections)[0].Pushes,
		},
		"from nil": {
			(*GetFullConfig(0).BranchProtections)[0].Pushes,
			nil,
			(*GetFullConfig(0).BranchProtections)[0].Pushes,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig->Protection->StatusChecks}
// and Repo->BranchProtection[...]{BaseGhBranchProtectionConfig->StatusChecks}

func updateGhBranchProtectStatusChecksConfigHelper(c *GhBranchProtectStatusChecksConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Strict, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.Required, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectStatusChecksConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchProtectStatusChecksConfig{},
		func(to, from *GhBranchProtectStatusChecksConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectStatusChecksConfigHelper,
	)
}

func TestGhBranchProtectStatusChecksConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].StatusChecks
	toWithNilSlicesAndStruct.Required = nil

	cases := map[string]struct {
		value    *GhBranchProtectStatusChecksConfig
		from     *GhBranchProtectStatusChecksConfig
		expected *GhBranchProtectStatusChecksConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).BranchProtections)[0].StatusChecks,
			(*GetFullConfig(1).BranchProtections)[0].StatusChecks,
		},
		"from nil": {
			(*GetFullConfig(0).BranchProtections)[0].StatusChecks,
			nil,
			(*GetFullConfig(0).BranchProtections)[0].StatusChecks,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig->Protection->PullRequestReviews}
// and Repo->BranchProtection[...]{BaseGhBranchProtectionConfig->PullRequestReviews}

func updateGhBranchProtectPRReviewConfigHelper(c *GhBranchProtectPRReviewConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.BypasserList, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.CodeownerApprovals, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.ResolvedConversations, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.ApprovalCount, stringToCopy, updatePtr)

	if c.Dismissals == nil {
		c.Dismissals = &GhBranchProtectPRReviewDismissalsConfig{}
	}
	updateGhBranchProtectPRReviewDismissalsConfigHelper(c.Dismissals, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPRReviewConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchProtectPRReviewConfig{},
		func(to, from *GhBranchProtectPRReviewConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPRReviewConfigHelper,
	)
}

func TestGhBranchProtectPRReviewConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].PullRequestReviews
	toWithNilSlicesAndStruct.BypasserList = nil
	toWithNilSlicesAndStruct.Dismissals = nil

	cases := map[string]struct {
		value    *GhBranchProtectPRReviewConfig
		from     *GhBranchProtectPRReviewConfig
		expected *GhBranchProtectPRReviewConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).BranchProtections)[0].PullRequestReviews,
			(*GetFullConfig(1).BranchProtections)[0].PullRequestReviews,
		},
		"from nil": {
			(*GetFullConfig(0).BranchProtections)[0].PullRequestReviews,
			nil,
			(*GetFullConfig(0).BranchProtections)[0].PullRequestReviews,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig->Protection->{BaseGhBranchProtectionConfig->PullRequestReviews->Dismissals}}
// and Repo->BranchProtection[...]{BaseGhBranchProtectionConfig->PullRequestReviews->Dismissals}

func updateGhBranchProtectPRReviewDismissalsConfigHelper(c *GhBranchProtectPRReviewDismissalsConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Staled, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Restrict, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.RestrictTo, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPRReviewDismissalsConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhBranchProtectPRReviewDismissalsConfig{},
		func(to, from *GhBranchProtectPRReviewDismissalsConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPRReviewDismissalsConfigHelper,
	)
}

func TestGhBranchProtectPRReviewDismissalsConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].PullRequestReviews.Dismissals
	toWithNilSlicesAndStruct.RestrictTo = nil

	cases := map[string]struct {
		value    *GhBranchProtectPRReviewDismissalsConfig
		from     *GhBranchProtectPRReviewDismissalsConfig
		expected *GhBranchProtectPRReviewDismissalsConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			(*GetFullConfig(1).BranchProtections)[0].PullRequestReviews.Dismissals,
			(*GetFullConfig(1).BranchProtections)[0].PullRequestReviews.Dismissals,
		},
		"from nil": {
			(*GetFullConfig(0).BranchProtections)[0].PullRequestReviews.Dismissals,
			nil,
			(*GetFullConfig(0).BranchProtections)[0].PullRequestReviews.Dismissals,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests

func updateGhRepoPullRequestConfigHelper(c *GhRepoPullRequestConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	if c.MergeStrategy == nil {
		c.MergeStrategy = &GhRepoPRMergeStrategyConfig{}
	}
	updateGhRepoPRMergeStrategyConfigHelper(c.MergeStrategy, stringToCopy, newSliceToCopy, updatePtr)
	if c.MergeCommit == nil {
		c.MergeCommit = &GhRepoPRCommitConfig{}
	}
	updateGhRepoPRCommitConfigHelper(c.MergeCommit, stringToCopy, newSliceToCopy, updatePtr)
	if c.SquashCommit == nil {
		c.SquashCommit = &GhRepoPRCommitConfig{}
	}
	updateGhRepoPRCommitConfigHelper(c.SquashCommit, stringToCopy, newSliceToCopy, updatePtr)
	if c.Branch == nil {
		c.Branch = &GhRepoPRBranchConfig{}
	}
	updateGhRepoPRBranchConfigHelper(c.Branch, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoPullRequestConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoPullRequestConfig{},
		func(to, from *GhRepoPullRequestConfig) {
			to.Merge(from)
		},
		updateGhRepoPullRequestConfigHelper,
	)
}

func TestGhRepoPullRequestConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0).PullRequests
	toWithNilSlicesAndStruct.MergeStrategy = nil
	toWithNilSlicesAndStruct.MergeCommit = nil
	toWithNilSlicesAndStruct.SquashCommit = nil
	toWithNilSlicesAndStruct.Branch = nil

	cases := map[string]struct {
		value    *GhRepoPullRequestConfig
		from     *GhRepoPullRequestConfig
		expected *GhRepoPullRequestConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1).PullRequests,
			GetFullConfig(1).PullRequests,
		},
		"from nil": {
			GetFullConfig(0).PullRequests,
			nil,
			GetFullConfig(0).PullRequests,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->MergeStrategy

func updateGhRepoPRMergeStrategyConfigHelper(c *GhRepoPRMergeStrategyConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.AllowMerge, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowRebase, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowSquash, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowAutoMerge, stringToCopy, updatePtr)
}

func TestGhRepoPRMergeStrategyConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoPRMergeStrategyConfig{},
		func(to, from *GhRepoPRMergeStrategyConfig) {
			to.Merge(from)
		},
		updateGhRepoPRMergeStrategyConfigHelper,
	)
}

func TestGhRepoPRMergeStrategyConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoPRMergeStrategyConfig
		from     *GhRepoPRMergeStrategyConfig
		expected *GhRepoPRMergeStrategyConfig
	}{
		"from nil": {
			GetFullConfig(0).PullRequests.MergeStrategy,
			nil,
			GetFullConfig(0).PullRequests.MergeStrategy,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->MergeCommit and Repo->PullRequests->SquashCommit

func updateGhRepoPRCommitConfigHelper(c *GhRepoPRCommitConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Title, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Message, stringToCopy, updatePtr)
}

func TestGhRepoPRCommitConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoPRCommitConfig{},
		func(to, from *GhRepoPRCommitConfig) {
			to.Merge(from)
		},
		updateGhRepoPRCommitConfigHelper,
	)
}

func TestGhRepoPRCommitConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoPRCommitConfig
		from     *GhRepoPRCommitConfig
		expected *GhRepoPRCommitConfig
	}{
		"from nil": {
			GetFullConfig(0).PullRequests.MergeCommit,
			nil,
			GetFullConfig(0).PullRequests.MergeCommit,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->Branch

func updateGhRepoPRBranchConfigHelper(c *GhRepoPRBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.SuggestUpdate, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.DeleteBranchOnMerge, stringToCopy, updatePtr)
}

func TestGhRepoPRBranchConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoPRBranchConfig{},
		func(to, from *GhRepoPRBranchConfig) {
			to.Merge(from)
		},
		updateGhRepoPRBranchConfigHelper,
	)
}

func TestGhRepoPRBranchConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoPRBranchConfig
		from     *GhRepoPRBranchConfig
		expected *GhRepoPRBranchConfig
	}{
		"from nil": {
			GetFullConfig(0).PullRequests.Branch,
			nil,
			GetFullConfig(0).PullRequests.Branch,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Security

func updateGhRepoSecurityConfigHelper(c *GhRepoSecurityConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.VulnerabilityAlerts, stringToCopy, updatePtr)
}

func TestGhRepoSecurityConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoSecurityConfig{},
		func(to, from *GhRepoSecurityConfig) {
			to.Merge(from)
		},
		updateGhRepoSecurityConfigHelper,
	)
}

func TestGhRepoSecurityConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoSecurityConfig
		from     *GhRepoSecurityConfig
		expected *GhRepoSecurityConfig
	}{
		"from nil": {
			GetFullConfig(0).Security,
			nil,
			GetFullConfig(0).Security,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous

func updateGhRepoMiscellaneousConfigHelper(c *GhRepoMiscellaneousConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.Topics, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.AutoInit, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Archived, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.IsTemplate, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.HomepageUrl, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.HasIssues, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.HasWiki, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.HasProjects, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.HasDownloads, stringToCopy, updatePtr)

	if c.Template == nil {
		c.Template = &GhRepoTemplateConfig{}
	}
	updateGhRepoTemplateConfigHelper(c.Template, stringToCopy, newSliceToCopy, updatePtr)

	if c.Pages == nil {
		c.Pages = &GhRepoPagesConfig{}
	}

	updateGhRepoPagesConfigHelper(c.Pages, stringToCopy, newSliceToCopy, updatePtr)

	if c.FileTemplates == nil {
		c.FileTemplates = &GhRepoFileTemplatesConfig{}
	}
	updateGhRepoFileTemplatesConfigHelper(c.FileTemplates, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoMiscellaneousConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoMiscellaneousConfig{},
		func(to, from *GhRepoMiscellaneousConfig) {
			to.Merge(from)
		},
		updateGhRepoMiscellaneousConfigHelper,
	)
}

func TestGhRepoMiscellaneousConfig_Merge_2(t *testing.T) {
	toWithNilSlicesAndStruct := GetFullConfig(0).Miscellaneous
	toWithNilSlicesAndStruct.Topics = nil
	toWithNilSlicesAndStruct.Template = nil
	toWithNilSlicesAndStruct.Pages = nil
	toWithNilSlicesAndStruct.FileTemplates = nil

	cases := map[string]struct {
		value    *GhRepoMiscellaneousConfig
		from     *GhRepoMiscellaneousConfig
		expected *GhRepoMiscellaneousConfig
	}{
		"to has nil slices and struct": {
			toWithNilSlicesAndStruct,
			GetFullConfig(1).Miscellaneous,
			GetFullConfig(1).Miscellaneous,
		},
		"from nil": {
			GetFullConfig(0).Miscellaneous,
			nil,
			GetFullConfig(0).Miscellaneous,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->Template

func updateGhRepoTemplateConfigHelper(c *GhRepoTemplateConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Source, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.FullClone, stringToCopy, updatePtr)
}

func TestGhRepoTemplateConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoTemplateConfig{},
		func(to, from *GhRepoTemplateConfig) {
			to.Merge(from)
		},
		updateGhRepoTemplateConfigHelper,
	)
}

func TestGhRepoTemplateConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoTemplateConfig
		from     *GhRepoTemplateConfig
		expected *GhRepoTemplateConfig
	}{
		"from nil": {
			GetFullConfig(0).Miscellaneous.Template,
			nil,
			GetFullConfig(0).Miscellaneous.Template,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->FileTemplates

func updateGhRepoFileTemplatesConfigHelper(c *GhRepoFileTemplatesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Gitignore, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.License, stringToCopy, updatePtr)
}

func TestGhRepoFileTemplatesConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoFileTemplatesConfig{},
		func(to, from *GhRepoFileTemplatesConfig) {
			to.Merge(from)
		},
		updateGhRepoFileTemplatesConfigHelper,
	)
}

func TestGhRepoFileTemplatesConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoFileTemplatesConfig
		from     *GhRepoFileTemplatesConfig
		expected *GhRepoFileTemplatesConfig
	}{
		"from nil": {
			GetFullConfig(0).Miscellaneous.FileTemplates,
			nil,
			GetFullConfig(0).Miscellaneous.FileTemplates,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->Pages

func updateGhRepoPagesConfigHelper(c *GhRepoPagesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Domain, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourcePath, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourceBranch, stringToCopy, updatePtr)
}

func TestGhRepoPagesConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoPagesConfig{},
		func(to, from *GhRepoPagesConfig) {
			to.Merge(from)
		},
		updateGhRepoPagesConfigHelper,
	)
}

func TestGhRepoPagesConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoPagesConfig
		from     *GhRepoPagesConfig
		expected *GhRepoPagesConfig
	}{
		"from nil": {
			GetFullConfig(0).Miscellaneous.Pages,
			nil,
			GetFullConfig(0).Miscellaneous.Pages,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Terraform

func updateGhRepoTerraformConfigHelper(c *GhRepoTerraformConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.ArchiveOnDestroy, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.IgnoreVulnerabilityAlertsDuringRead, stringToCopy, updatePtr)
}

func TestGhRepoTerraformConfig_Merge(t *testing.T) {
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		GhRepoTerraformConfig{},
		func(to, from *GhRepoTerraformConfig) {
			to.Merge(from)
		},
		updateGhRepoTerraformConfigHelper,
	)
}

func TestGhRepoTerraformConfig_Merge_2(t *testing.T) {
	cases := map[string]struct {
		value    *GhRepoTerraformConfig
		from     *GhRepoTerraformConfig
		expected *GhRepoTerraformConfig
	}{
		"from nil": {
			GetFullConfig(0).Terraform,
			nil,
			GetFullConfig(0).Terraform,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				tc.value.Merge(tc.from)
				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

/** Internal **/

type (
	updateStructForTest[T any]  func(c *T, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool)
	createCopyWithFnType[T any] func(stringVal *string, sliceVal *[]string) *T
)

func ensureMergeWithoutOverflowBetweenToAndFrom[T any](t *testing.T, empty T, mergeFn func(to, from *T), updateFn updateStructForTest[T]) {
	initialToStringValue := "initial_to_value"
	initialToSliceValue := []string{"initial_to_slice_value1", "initial_to_slice_value2"}
	initialFromStringValue := "initial_from_value"
	initialFromSliceValue := []string{"initial_from_slice_value1", "initial_from_slice_value2"}

	initialFromStringValueCopy := fmt.Sprintf("%s", initialFromStringValue)
	initialFromSliceValueCopy := []string{fmt.Sprintf("%s", initialFromSliceValue[0]), fmt.Sprintf("%s", initialFromSliceValue[1])}

	createCopyWithFn := func(stringVal *string, sliceVal *[]string) *T {
		v := reflect.New(reflect.TypeOf(empty)).Interface().(*T)
		updateFn(v, stringVal, sliceVal, true)

		return v
	}

	to := createCopyWithFn(&initialToStringValue, &initialToSliceValue)
	// Check that merging from nil has no impact
	expectedInitialTo := createCopyWithFn(&initialToStringValue, &initialToSliceValue)
	mergeFn(to, nil)
	if diff := cmp.Diff(expectedInitialTo, to); diff != "" {
		t.Fatalf("'to' merged with nil must not impact 'to': Config mismatch (-want +got):\n%s", diff)
	}

	// Check that merging from non nil do have an impact
	from := createCopyWithFn(&initialFromStringValue, &initialFromSliceValue)
	mergeFn(to, from)
	// All values must be equals to 'from' values and slice must be a combination of both - Use copy in case original values are updated !
	expectedToSliceAfterMerge := append(initialToSliceValue, initialFromSliceValue...)
	fmt.Printf("expectedToSliceAfterMerge: %#v\n", expectedToSliceAfterMerge)
	expectedToAfterMerge := createCopyWithFn(&initialFromStringValueCopy, &expectedToSliceAfterMerge)
	// Initial check - ensure merge worked as expected
	if diff := cmp.Diff(expectedToAfterMerge, to); diff != "" {
		t.Fatalf("'to' merged with 'from' must impact 'to': Config mismatch (-want +got):\n%s", diff)
	}
	// Must be the same as original 'from' as merge must not impact from
	expectedFromAfterMerge := createCopyWithFn(&initialFromStringValueCopy, &initialFromSliceValueCopy)
	if diff := cmp.Diff(expectedFromAfterMerge, from); diff != "" {
		t.Fatalf("'to' merged with 'from' must not impact 'from': Config mismatch (-want +got):\n%s", diff)
	}

	// expectedInitialTo is not anymore the same since the merge => must have all 'from' values
	expectedInitialTo = createCopyWithFn(&initialFromStringValueCopy, &expectedToSliceAfterMerge)

	ensureNoOverflowBetweenToAndFrom(t, to, from, createCopyWithFn, updateFn, expectedInitialTo)
}

func ensureNoOverflowBetweenToAndFrom[T any](t *testing.T, to *T, from *T, createCopyWithFn createCopyWithFnType[T], updateFn updateStructForTest[T], expectedInitialTo *T) {
	updatedToStringValue := "updated_to_value"
	updatedToSliceValue := []string{"updated_to_slice_value1", "updated_to_slice_value2"}
	updatedToStringPtrValue := "updated_to_ptr_value"
	updatedToSlicePtrValue := []string{"updated_to_slice_ptr_value1", "updated_to_slice_ptr_value2"}

	updatedFromStringValue := "updated_from_value"
	updatedFromSliceValue := []string{"updated_from_slice_value1", "updated_from_slice_value2"}
	updatedFromStringPtrValue := "updated_from_ptr_value"
	updatedFromSlicePtrValue := []string{"updated_from_slice_ptr_value1", "updated_from_slice_ptr_value2"}

	updatedToStringValueCopy := fmt.Sprintf("%s", updatedToStringValue)
	updatedFromStringValueCopy := fmt.Sprintf("%s", updatedFromStringValue)

	updatedToSliceValueCopy := []string{fmt.Sprintf("%s", updatedToSliceValue[0]), fmt.Sprintf("%s", updatedToSliceValue[1])}
	updatedFromSliceValueCopy := []string{fmt.Sprintf("%s", updatedFromSliceValue[0]), fmt.Sprintf("%s", updatedFromSliceValue[1])}

	updatedToStringPtrValueCopy := fmt.Sprintf("%s", updatedToStringPtrValue)
	updatedToStringPtrCopy := &updatedToStringPtrValueCopy
	updatedFromStringPtrValueCopy := fmt.Sprintf("%s", updatedFromStringPtrValue)
	updatedFromStringPtrCopy := &updatedFromStringPtrValueCopy

	updatedToSlicePtrValueCopy := []string{fmt.Sprintf("%s", updatedToSlicePtrValue[0]), fmt.Sprintf("%s", updatedToSlicePtrValue[1])}
	updatedToSlicePtrCopy := &updatedToSlicePtrValueCopy
	updatedFromSlicePtrValueCopy := []string{fmt.Sprintf("%s", updatedFromSlicePtrValue[0]), fmt.Sprintf("%s", updatedFromSlicePtrValue[1])}
	updatedFromSlicePtrCopy := &updatedFromSlicePtrValueCopy

	// All values must be equals to 'to' updated values - Use copy in case original values are updated !
	expectedToAfterToUpdate := createCopyWithFn(&updatedToStringValueCopy, &updatedToSliceValueCopy)
	// All values must be equals to 'to' updated pointers values - Use copy in case original values are updated !
	expectedToAfterToPointersUpdate := createCopyWithFn(updatedToStringPtrCopy, updatedToSlicePtrCopy)
	// All values must be equals to 'from' updated values - Use copy in case original values are updated !
	expectedFromAfterFromValuesUpdate := createCopyWithFn(&updatedFromStringValueCopy, &updatedFromSliceValueCopy)
	// All values must be equals to 'from' updated pointers values - Use copy in case original values are updated !
	expectedFromAfterFromPointersUpdate := createCopyWithFn(updatedFromStringPtrCopy, updatedFromSlicePtrCopy)

	// Case 1: Updating 'from' *values must not impact 'to' *values
	updateFn(from, &updatedFromStringValue, &updatedFromSliceValue, false)
	// case1Fn()
	if diff := cmp.Diff(expectedFromAfterFromValuesUpdate, from); diff != "" {
		t.Fatalf("updating 'from' *values must impact 'from': Config mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedInitialTo, to); diff != "" {
		t.Fatalf("updating 'from' *values must not impact 'to': Config mismatch (-want +got):\n%s", diff)
	}

	// Case 2: Updating 'from' pointers must not impact 'to' pointers
	updateFn(from, &updatedFromStringPtrValue, &updatedFromSlicePtrValue, true)
	if diff := cmp.Diff(expectedFromAfterFromPointersUpdate, from); diff != "" {
		t.Fatalf("updating 'from' pointers must impact 'from': Config mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedInitialTo, to); diff != "" {
		t.Fatalf("updating 'from' pointers must not impact 'to': Config mismatch (-want +got):\n%s", diff)
	}

	// Case 3: Updating 'to' *values must not impact 'from' *values
	updateFn(to, &updatedToStringValue, &updatedToSliceValue, false)
	if diff := cmp.Diff(expectedToAfterToUpdate, to); diff != "" {
		t.Fatalf("updating 'to' *values must impact 'to': Config mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedFromAfterFromPointersUpdate, from); diff != "" {
		t.Fatalf("updating 'to' *values must not impact 'from': Config mismatch (-want +got):\n%s", diff)
	}

	// Case 4: Updating 'to' pointers must not impact 'from' pointers
	updateFn(to, &updatedToStringPtrValue, &updatedToSlicePtrValue, true)
	if diff := cmp.Diff(expectedToAfterToPointersUpdate, to); diff != "" {
		t.Fatalf("updating 'to' pointers must impact 'to': Config mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(expectedFromAfterFromPointersUpdate, from); diff != "" {
		t.Fatalf("updating 'to' pointers must not impact 'from': Config mismatch (-want +got):\n%s", diff)
	}
}

func updateStringPtrHelper(to **string, stringToCopy *string, updatePtr bool) {
	newString := fmt.Sprintf("%s", *stringToCopy)
	if updatePtr {
		*to = &newString
	} else {
		**to = newString
	}
}

func updateSlicePtrHelper(to **[]string, newSliceToCopy *[]string, updatePtr bool) {
	newSlice := append([]string{}, *newSliceToCopy...)
	if updatePtr {
		*to = &newSlice
	} else {
		**to = newSlice
	}
}
