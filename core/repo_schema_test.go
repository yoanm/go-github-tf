package core_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/go-github-tf/core"
)

// ==> Repo

func updateGhRepoConfigHelper(c *core.GhRepoConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Name, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.Visibility, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Description, stringToCopy, updatePtr)

	if c.DefaultBranch == nil {
		c.DefaultBranch = &core.GhDefaultBranchConfig{}
	}

	updateGhDefaultBranchConfigHelper(c.DefaultBranch, stringToCopy, newSliceToCopy, updatePtr)

	/*
		// Do nothing with Branches as it's complicated to test it with that way
		// See TestGhBranchesConfig_Merge instead
		if c.Branches == nil {
			c.Branches = &core.GhBranchesConfig{}
		}
		updateGhBranchesConfigHelper(c.Branches, stringToCopy, newSliceToCopy, updatePtr)
	*/

	/*
		// Do nothing with BranchProtections as it's complicated to test it with that way
		// See TestGhBranchProtectionsConfig_Merge instead
		if c.BranchProtections == nil {
			c.BranchProtections = &core.GhBranchProtectionsConfig{}
		}
		updateGhBranchProtectionsConfigHelper(c.BranchProtections, stringToCopy, newSliceToCopy, updatePtr)
	*/

	if c.PullRequests == nil {
		c.PullRequests = &core.GhRepoPullRequestConfig{}
	}

	updateGhRepoPullRequestConfigHelper(c.PullRequests, stringToCopy, newSliceToCopy, updatePtr)

	if c.Security == nil {
		c.Security = &core.GhRepoSecurityConfig{}
	}

	updateGhRepoSecurityConfigHelper(c.Security, stringToCopy, newSliceToCopy, updatePtr)

	if c.Miscellaneous == nil {
		c.Miscellaneous = &core.GhRepoMiscellaneousConfig{}
	}

	updateGhRepoMiscellaneousConfigHelper(c.Miscellaneous, stringToCopy, newSliceToCopy, updatePtr)

	if c.Terraform == nil {
		c.Terraform = &core.GhRepoTerraformConfig{}
	}

	updateGhRepoTerraformConfigHelper(c.Terraform, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoConfig{},
		func(to, from *core.GhRepoConfig) {
			to.Merge(from)
		},
		updateGhRepoConfigHelper,
	)
}

func TestGhRepoConfig_Merge_2(t *testing.T) {
	t.Parallel()

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

	*fullMergeResult.DefaultBranch.Protection.Pushes.RestrictTo = append(
		*(full1.DefaultBranch.Protection.Pushes.RestrictTo),
		*(full2.DefaultBranch.Protection.Pushes.RestrictTo)...,
	)

	*fullMergeResult.DefaultBranch.Protection.StatusChecks.Required = append(
		*(full1.DefaultBranch.Protection.StatusChecks.Required),
		*(full2.DefaultBranch.Protection.StatusChecks.Required)...,
	)

	*fullMergeResult.DefaultBranch.Protection.PullRequestReviews.Bypassers = append(
		*(full1.DefaultBranch.Protection.PullRequestReviews.Bypassers),
		*(full2.DefaultBranch.Protection.PullRequestReviews.Bypassers)...,
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

	*(*fullMergeResult.Branches)["feature/branch2"].Protection.Pushes.RestrictTo = append(
		*((*full1.Branches)["feature/branch2"].Protection.Pushes.RestrictTo),
		*((*full2.Branches)["feature/branch2"].Protection.Pushes.RestrictTo)...,
	)

	*(*fullMergeResult.Branches)["feature/branch2"].Protection.StatusChecks.Required = append(
		*((*full1.Branches)["feature/branch2"].Protection.StatusChecks.Required),
		*((*full2.Branches)["feature/branch2"].Protection.StatusChecks.Required)...,
	)

	*(*fullMergeResult.Branches)["feature/branch2"].Protection.PullRequestReviews.Bypassers = append(
		*((*full1.Branches)["feature/branch2"].Protection.PullRequestReviews.Bypassers),
		*((*full2.Branches)["feature/branch2"].Protection.PullRequestReviews.Bypassers)...,
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
		value    *core.GhRepoConfig
		from     *core.GhRepoConfig
		expected *core.GhRepoConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->DefaultBranch

func updateGhDefaultBranchConfigHelper(c *core.GhDefaultBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Name, stringToCopy, updatePtr)

	updateBaseGhBranchConfigHelper(&c.BaseGhBranchConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhDefaultBranchConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhDefaultBranchConfig{},
		func(to, from *core.GhDefaultBranchConfig) {
			to.Merge(from)
		},
		updateGhDefaultBranchConfigHelper,
	)
}

func TestGhDefaultBranchConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := GetFullConfig(0).DefaultBranch
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *core.GhDefaultBranchConfig
		from     *core.GhDefaultBranchConfig
		expected *core.GhDefaultBranchConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches

func TestGhBranchesConfig_Merge(t *testing.T) {
	t.Parallel()
	// Init variables
	initialToStringValue := "initial_to_value"

	initialFromStringValue := "initial_from_value"

	// Init 'to' variables and create 'to'
	toSourceBranch := fmt.Sprintf("%s", initialToStringValue)                 //nolint:perfsprint // Because :p
	toSourceSha := fmt.Sprintf("%s", initialToStringValue)                    //nolint:perfsprint // Because :p
	toEnforceAdmins := fmt.Sprintf("%s", initialToStringValue)                //nolint:perfsprint // Because :p
	toAllowsDeletions := fmt.Sprintf("%s", initialToStringValue)              //nolint:perfsprint // Because :p
	toAllowsForcePushes := fmt.Sprintf("%s", initialToStringValue)            //nolint:perfsprint // Because :p
	toRequiredLinearHistory := fmt.Sprintf("%s", initialToStringValue)        //nolint:perfsprint // Because :p
	toRequireSignedCommits := fmt.Sprintf("%s", initialToStringValue)         //nolint:perfsprint // Because :p
	toStrict := fmt.Sprintf("%s", initialToStringValue)                       //nolint:perfsprint // Because :p
	toDismissStaleReviews := fmt.Sprintf("%s", initialToStringValue)          //nolint:perfsprint // Because :p
	toRequireCodeOwnerReviews := fmt.Sprintf("%s", initialToStringValue)      //nolint:perfsprint // Because :p
	toResolvedConversations := fmt.Sprintf("%s", initialToStringValue)        //nolint:perfsprint // Because :p
	toRequiredApprovingReviewCount := fmt.Sprintf("%s", initialToStringValue) //nolint:perfsprint // Because :p
	// Do nothing with PushRestrictions, see TestGhRepoConfig_Merge2
	// toPushRestrictions := append([]string{}, initialToSliceValue...)
	// Do nothing with Contexts, see TestGhRepoConfig_Merge2
	// toContext := append([]string{}, initialToSliceValue...)
	// Do nothing with DismissalRestrictions, see TestGhRepoConfig_Merge2
	// toDismissalRestrictions := append([]string{}, initialToSliceValue...)
	// Do nothing with ConfigTemplates, see TestGhRepoConfig_Merge2
	// toConfigTemplates := append([]string{}, initialFromSliceValue...)
	to := &core.GhBranchesConfig{
		"to_branch": &core.GhBranchConfig{
			SourceBranch: &toSourceBranch,
			SourceSha:    &toSourceSha,
			BaseGhBranchConfig: core.BaseGhBranchConfig{
				// ConfigTemplates: &toConfigTemplates,
				Protection: &core.BaseGhBranchProtectionConfig{
					// ConfigTemplates:       &toConfigTemplates,
					EnforceAdmins: &toEnforceAdmins,
					AllowDeletion: &toAllowsDeletions,
					Pushes: &core.GhBranchProtectPushesConfig{
						AllowsForcePushes: &toAllowsForcePushes,
						// PushRestrictions:      &toPushRestrictions,
					},
					RequireLinearHistory: &toRequiredLinearHistory,
					RequireSignedCommits: &toRequireSignedCommits,
					StatusChecks: &core.GhBranchProtectStatusChecksConfig{
						Strict: &toStrict,
						//	Contexts: &toContext,
					},
					PullRequestReviews: &core.GhBranchProtectPRReviewConfig{
						Dismissals: &core.GhBranchProtectPRReviewDismissalsConfig{
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
	fromSourceBranch := fmt.Sprintf("%s", initialFromStringValue)                 //nolint:perfsprint // Because :p
	fromSourceSha := fmt.Sprintf("%s", initialFromStringValue)                    //nolint:perfsprint // Because :p
	fromEnforceAdmins := fmt.Sprintf("%s", initialFromStringValue)                //nolint:perfsprint // Because :p
	fromAllowsDeletions := fmt.Sprintf("%s", initialFromStringValue)              //nolint:perfsprint // Because :p
	fromAllowsForcePushes := fmt.Sprintf("%s", initialFromStringValue)            //nolint:perfsprint // Because :p
	fromRequiredLinearHistory := fmt.Sprintf("%s", initialFromStringValue)        //nolint:perfsprint // Because :p
	fromRequireSignedCommits := fmt.Sprintf("%s", initialFromStringValue)         //nolint:perfsprint // Because :p
	fromStrict := fmt.Sprintf("%s", initialFromStringValue)                       //nolint:perfsprint // Because :p
	fromDismissStaleReviews := fmt.Sprintf("%s", initialFromStringValue)          //nolint:perfsprint // Because :p
	fromRequireCodeOwnerReviews := fmt.Sprintf("%s", initialFromStringValue)      //nolint:perfsprint // Because :p
	fromResolvedConversations := fmt.Sprintf("%s", initialToStringValue)          //nolint:perfsprint // Because :p
	fromRequiredApprovingReviewCount := fmt.Sprintf("%s", initialFromStringValue) //nolint:perfsprint // Because :p
	// Do nothing with PushRestrictions, see TestGhRepoConfig_Merge2
	// fromPushRestrictions := append([]string{}, initialFromSliceValue...)
	// Do nothing with Contexts, see TestGhRepoConfig_Merge2
	// fromContext := append([]string{}, initialFromSliceValue...)
	// Do nothing with DismissalRestrictions, see TestGhRepoConfig_Merge2
	// fromDismissalRestrictions := append([]string{}, initialFromSliceValue...)
	// Do nothing with ConfigTemplates, see TestGhRepoConfig_Merge2
	// fromConfigTemplates := append([]string{}, initialFromSliceValue...)
	from := &core.GhBranchesConfig{
		"from_branch": &core.GhBranchConfig{
			SourceBranch: &fromSourceBranch,
			SourceSha:    &fromSourceSha,
			BaseGhBranchConfig: core.BaseGhBranchConfig{
				// ConfigTemplates: &fromConfigTemplates,
				Protection: &core.BaseGhBranchProtectionConfig{
					// ConfigTemplates:       &fromConfigTemplates,
					EnforceAdmins: &fromEnforceAdmins,
					AllowDeletion: &fromAllowsDeletions,
					Pushes: &core.GhBranchProtectPushesConfig{
						AllowsForcePushes: &fromAllowsForcePushes,
						// PushRestrictions:      &fromPushRestrictions,
					},
					RequireLinearHistory: &fromRequiredLinearHistory,
					RequireSignedCommits: &fromRequireSignedCommits,
					StatusChecks: &core.GhBranchProtectStatusChecksConfig{
						Strict: &fromStrict,
						// Contexts: &fromContext,
					},
					PullRequestReviews: &core.GhBranchProtectPRReviewConfig{
						Dismissals: &core.GhBranchProtectPRReviewDismissalsConfig{
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
	expected := &core.GhBranchesConfig{
		"to_branch":   (*to)["to_branch"],
		"from_branch": (*from)["from_branch"],
	}

	to.Merge(from)

	if diff := cmp.Diff(expected, to); diff != "" {
		t.Fatalf("Config mismatch (-want +got):\n%s", diff)
	}
}

func TestGhBranchesConfig_Merge_2(t *testing.T) {
	t.Parallel()

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
		value    *core.GhBranchesConfig
		from     *core.GhBranchesConfig
		expected *core.GhBranchesConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]

func updateGhBranchConfigHelper(c *core.GhBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.SourceBranch, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourceSha, stringToCopy, updatePtr)

	updateBaseGhBranchConfigHelper(&c.BaseGhBranchConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchConfig{},
		func(to, from *core.GhBranchConfig) {
			to.Merge(from)
		},
		updateGhBranchConfigHelper,
	)
}

func TestGhBranchConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).Branches)["feature/branch0"]
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *core.GhBranchConfig
		from     *core.GhBranchConfig
		expected *core.GhBranchConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Branches[...]{BaseGhBranchConfig}

func updateBaseGhBranchConfigHelper(c *core.BaseGhBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)

	if c.Protection == nil {
		c.Protection = &core.BaseGhBranchProtectionConfig{}
	}

	updateBaseGhBranchProtectionConfigHelper(c.Protection, stringToCopy, newSliceToCopy, updatePtr)
}

func TestBaseGhBranchConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.BaseGhBranchConfig{},
		func(to, from *core.BaseGhBranchConfig) {
			to.Merge(from)
		},
		updateBaseGhBranchConfigHelper,
	)
}

func TestBaseGhBranchConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).Branches)["feature/branch0"].BaseGhBranchConfig
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Protection = nil

	cases := map[string]struct {
		value    *core.BaseGhBranchConfig
		from     *core.BaseGhBranchConfig
		expected *core.BaseGhBranchConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->BranchProtections.

func TestGhBranchProtectionsConfig_Merge(t *testing.T) {
	t.Parallel()
	// Init variables
	initialToStringValue := "initial_to_value"
	initialToSliceValue := []string{"initial_to_slice_value1", "initial_to_slice_value2"}

	initialFromStringValue := "initial_from_value"
	initialFromSliceValue := []string{"initial_from_slice_value1", "initial_from_slice_value2"}

	// Init 'to' variables and create 'to'
	toPattern := fmt.Sprintf("%s", initialToStringValue)                      //nolint:perfsprint // Because :p
	toForbid := fmt.Sprintf("%s", initialToStringValue)                       //nolint:perfsprint // Because :p
	toEnforceAdmins := fmt.Sprintf("%s", initialToStringValue)                //nolint:perfsprint // Because :p
	toAllowsDeletions := fmt.Sprintf("%s", initialToStringValue)              //nolint:perfsprint // Because :p
	toAllowsForcePushes := fmt.Sprintf("%s", initialToStringValue)            //nolint:perfsprint // Because :p
	toRequiredLinearHistory := fmt.Sprintf("%s", initialToStringValue)        //nolint:perfsprint // Because :p
	toRequireSignedCommits := fmt.Sprintf("%s", initialToStringValue)         //nolint:perfsprint // Because :p
	toStrict := fmt.Sprintf("%s", initialToStringValue)                       //nolint:perfsprint // Because :p
	toDismissStaleReviews := fmt.Sprintf("%s", initialToStringValue)          //nolint:perfsprint // Because :p
	toRestrict := fmt.Sprintf("%s", initialToStringValue)                     //nolint:perfsprint // Because :p
	toRequireCodeOwnerReviews := fmt.Sprintf("%s", initialToStringValue)      //nolint:perfsprint // Because :p
	toResolvedConversations := fmt.Sprintf("%s", initialToStringValue)        //nolint:perfsprint // Because :p
	toRequiredApprovingReviewCount := fmt.Sprintf("%s", initialToStringValue) //nolint:perfsprint // Because :p

	toBypasserList := append([]string{}, initialToSliceValue...)

	toPushRestrictions := append([]string{}, initialToSliceValue...)

	toContext := append([]string{}, initialToSliceValue...)

	toDismissalRestrictions := append([]string{}, initialToSliceValue...)

	toConfigTemplates := append([]string{}, initialFromSliceValue...)

	to := &core.GhBranchProtectionsConfig{
		{
			&toPattern,
			&toForbid,
			core.BaseGhBranchProtectionConfig{
				&toConfigTemplates,
				&toEnforceAdmins,
				&toAllowsDeletions,
				&toRequiredLinearHistory,
				&toRequireSignedCommits,
				&core.GhBranchProtectPushesConfig{
					&toAllowsForcePushes,
					&toPushRestrictions,
				},
				&core.GhBranchProtectStatusChecksConfig{
					&toStrict,
					&toContext,
				},
				&core.GhBranchProtectPRReviewConfig{
					&toBypasserList,
					&toRequireCodeOwnerReviews,
					&toResolvedConversations,
					&toRequiredApprovingReviewCount,
					&core.GhBranchProtectPRReviewDismissalsConfig{
						&toDismissStaleReviews,
						&toRestrict,
						&toDismissalRestrictions,
					},
				},
			},
		},
	}

	// Init 'from' variables and create 'from'
	fromPattern := fmt.Sprintf("%s", initialFromStringValue)                      //nolint:perfsprint // Because :p
	fromForbid := fmt.Sprintf("%s", initialFromStringValue)                       //nolint:perfsprint // Because :p
	fromEnforceAdmins := fmt.Sprintf("%s", initialFromStringValue)                //nolint:perfsprint // Because :p
	fromAllowsDeletions := fmt.Sprintf("%s", initialFromStringValue)              //nolint:perfsprint // Because :p
	fromAllowsForcePushes := fmt.Sprintf("%s", initialFromStringValue)            //nolint:perfsprint // Because :p
	fromRequiredLinearHistory := fmt.Sprintf("%s", initialFromStringValue)        //nolint:perfsprint // Because :p
	fromRequireSignedCommits := fmt.Sprintf("%s", initialFromStringValue)         //nolint:perfsprint // Because :p
	fromStrict := fmt.Sprintf("%s", initialFromStringValue)                       //nolint:perfsprint // Because :p
	fromDismissStaleReviews := fmt.Sprintf("%s", initialFromStringValue)          //nolint:perfsprint // Because :p
	fromRestrict := fmt.Sprintf("%s", initialFromStringValue)                     //nolint:perfsprint // Because :p
	fromRequireCodeOwnerReviews := fmt.Sprintf("%s", initialFromStringValue)      //nolint:perfsprint // Because :p
	fromResolvedConversations := fmt.Sprintf("%s", initialFromStringValue)        //nolint:perfsprint // Because :p
	fromRequiredApprovingReviewCount := fmt.Sprintf("%s", initialFromStringValue) //nolint:perfsprint // Because :p

	fromPushRestrictions := append([]string{}, initialFromSliceValue...)

	fromContext := append([]string{}, initialFromSliceValue...)

	fromBypasserList := append([]string{}, initialFromSliceValue...)

	fromDismissalRestrictions := append([]string{}, initialFromSliceValue...)

	fromConfigTemplates := append([]string{}, initialFromSliceValue...)

	from := &core.GhBranchProtectionsConfig{
		{
			&fromPattern,
			&fromForbid,
			core.BaseGhBranchProtectionConfig{
				&fromConfigTemplates,
				&fromEnforceAdmins,
				&fromAllowsDeletions,
				&fromRequiredLinearHistory,
				&fromRequireSignedCommits,
				&core.GhBranchProtectPushesConfig{
					&fromAllowsForcePushes,
					&fromPushRestrictions,
				},
				&core.GhBranchProtectStatusChecksConfig{
					&fromStrict,
					&fromContext,
				},
				&core.GhBranchProtectPRReviewConfig{
					&fromBypasserList,
					&fromRequireCodeOwnerReviews,
					&fromResolvedConversations,
					&fromRequiredApprovingReviewCount,
					&core.GhBranchProtectPRReviewDismissalsConfig{
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

	createCopyWithFn := func(stringVal *string, sliceVal *[]string) *core.GhBranchProtectionConfig {
		v := reflect.New(reflect.TypeOf(core.GhBranchProtectionConfig{})).Interface().(*core.GhBranchProtectionConfig)
		updateFn(v, stringVal, sliceVal, true)

		return v
	}

	// expectedInitialTo refers to the second item of 'to'
	// => following to.merge(from), it is equals to fromItem as merge append slice items
	expectedInitialTo := createCopyWithFn(&initialFromStringValue, &initialFromSliceValue)
	ensureNoOverflowBetweenToAndFrom(t, toItem, fromItem, createCopyWithFn, updateFn, expectedInitialTo)
}

func TestGhBranchProtectionsConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := GetFullConfig(0).BranchProtections
	(*toWithNilSlicesAndStruct)[0].ConfigTemplates = nil
	(*toWithNilSlicesAndStruct)[0].Pushes = nil
	(*toWithNilSlicesAndStruct)[0].StatusChecks = nil
	(*toWithNilSlicesAndStruct)[0].PullRequestReviews = nil

	expectedToWithNilSlicesAndStruct := GetFullConfig(1).BranchProtections
	*expectedToWithNilSlicesAndStruct = append(*toWithNilSlicesAndStruct, *expectedToWithNilSlicesAndStruct...)

	cases := map[string]struct {
		value    *core.GhBranchProtectionsConfig
		from     *core.GhBranchProtectionsConfig
		expected *core.GhBranchProtectionsConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->BranchProtections[...]

func updateGhBranchProtectionConfigHelper(c *core.GhBranchProtectionConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Pattern, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Forbid, stringToCopy, updatePtr)

	updateBaseGhBranchProtectionConfigHelper(&c.BaseGhBranchProtectionConfig, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectionConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchProtectionConfig{},
		func(to, from *core.GhBranchProtectionConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectionConfigHelper,
	)
}

func TestGhBranchProtectionConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0]
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Pushes = nil
	toWithNilSlicesAndStruct.StatusChecks = nil
	toWithNilSlicesAndStruct.PullRequestReviews = nil

	cases := map[string]struct {
		value    *core.GhBranchProtectionConfig
		from     *core.GhBranchProtectionConfig
		expected *core.GhBranchProtectionConfig
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
				t.Parallel()
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

func updateBaseGhBranchProtectionConfigHelper(c *core.BaseGhBranchProtectionConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.ConfigTemplates, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.EnforceAdmins, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowDeletion, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.RequireLinearHistory, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.RequireSignedCommits, stringToCopy, updatePtr)

	if c.Pushes == nil {
		c.Pushes = &core.GhBranchProtectPushesConfig{}
	}

	updateGhBranchProtectPushesConfigHelper(c.Pushes, stringToCopy, newSliceToCopy, updatePtr)

	if c.StatusChecks == nil {
		c.StatusChecks = &core.GhBranchProtectStatusChecksConfig{}
	}

	updateGhBranchProtectStatusChecksConfigHelper(c.StatusChecks, stringToCopy, newSliceToCopy, updatePtr)

	if c.PullRequestReviews == nil {
		c.PullRequestReviews = &core.GhBranchProtectPRReviewConfig{}
	}

	updateGhBranchProtectPRReviewConfigHelper(c.PullRequestReviews, stringToCopy, newSliceToCopy, updatePtr)
}

func TestBaseGhBranchProtectionConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.BaseGhBranchProtectionConfig{},
		func(to, from *core.BaseGhBranchProtectionConfig) {
			to.Merge(from)
		},
		updateBaseGhBranchProtectionConfigHelper,
	)
}

func TestBaseGhBranchProtectionConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].BaseGhBranchProtectionConfig
	toWithNilSlicesAndStruct.ConfigTemplates = nil
	toWithNilSlicesAndStruct.Pushes = nil
	toWithNilSlicesAndStruct.StatusChecks = nil
	toWithNilSlicesAndStruct.PullRequestReviews = nil

	cases := map[string]struct {
		value    *core.BaseGhBranchProtectionConfig
		from     *core.BaseGhBranchProtectionConfig
		expected *core.BaseGhBranchProtectionConfig
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
				t.Parallel()
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

func updateGhBranchProtectPushesConfigHelper(c *core.GhBranchProtectPushesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.AllowsForcePushes, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.RestrictTo, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPushesConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchProtectPushesConfig{},
		func(to, from *core.GhBranchProtectPushesConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPushesConfigHelper,
	)
}

func TestGhBranchProtectPushesConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].Pushes
	toWithNilSlicesAndStruct.RestrictTo = nil

	cases := map[string]struct {
		value    *core.GhBranchProtectPushesConfig
		from     *core.GhBranchProtectPushesConfig
		expected *core.GhBranchProtectPushesConfig
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
				t.Parallel()
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

func updateGhBranchProtectStatusChecksConfigHelper(c *core.GhBranchProtectStatusChecksConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Strict, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.Required, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectStatusChecksConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchProtectStatusChecksConfig{},
		func(to, from *core.GhBranchProtectStatusChecksConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectStatusChecksConfigHelper,
	)
}

func TestGhBranchProtectStatusChecksConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].StatusChecks
	toWithNilSlicesAndStruct.Required = nil

	cases := map[string]struct {
		value    *core.GhBranchProtectStatusChecksConfig
		from     *core.GhBranchProtectStatusChecksConfig
		expected *core.GhBranchProtectStatusChecksConfig
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
				t.Parallel()
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

func updateGhBranchProtectPRReviewConfigHelper(c *core.GhBranchProtectPRReviewConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateSlicePtrHelper(&c.Bypassers, newSliceToCopy, updatePtr)
	updateStringPtrHelper(&c.CodeownerApprovals, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.ResolvedConversations, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.ApprovalCount, stringToCopy, updatePtr)

	if c.Dismissals == nil {
		c.Dismissals = &core.GhBranchProtectPRReviewDismissalsConfig{}
	}

	updateGhBranchProtectPRReviewDismissalsConfigHelper(c.Dismissals, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPRReviewConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchProtectPRReviewConfig{},
		func(to, from *core.GhBranchProtectPRReviewConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPRReviewConfigHelper,
	)
}

func TestGhBranchProtectPRReviewConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].PullRequestReviews
	toWithNilSlicesAndStruct.Bypassers = nil
	toWithNilSlicesAndStruct.Dismissals = nil

	cases := map[string]struct {
		value    *core.GhBranchProtectPRReviewConfig
		from     *core.GhBranchProtectPRReviewConfig
		expected *core.GhBranchProtectPRReviewConfig
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
				t.Parallel()
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

func updateGhBranchProtectPRReviewDismissalsConfigHelper(c *core.GhBranchProtectPRReviewDismissalsConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Staled, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Restrict, stringToCopy, updatePtr)
	updateSlicePtrHelper(&c.RestrictTo, newSliceToCopy, updatePtr)
}

func TestGhBranchProtectPRReviewDismissalsConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhBranchProtectPRReviewDismissalsConfig{},
		func(to, from *core.GhBranchProtectPRReviewDismissalsConfig) {
			to.Merge(from)
		},
		updateGhBranchProtectPRReviewDismissalsConfigHelper,
	)
}

func TestGhBranchProtectPRReviewDismissalsConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := (*GetFullConfig(0).BranchProtections)[0].PullRequestReviews.Dismissals
	toWithNilSlicesAndStruct.RestrictTo = nil

	cases := map[string]struct {
		value    *core.GhBranchProtectPRReviewDismissalsConfig
		from     *core.GhBranchProtectPRReviewDismissalsConfig
		expected *core.GhBranchProtectPRReviewDismissalsConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests

func updateGhRepoPullRequestConfigHelper(c *core.GhRepoPullRequestConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	if c.MergeStrategy == nil {
		c.MergeStrategy = &core.GhRepoPRMergeStrategyConfig{}
	}

	updateGhRepoPRMergeStrategyConfigHelper(c.MergeStrategy, stringToCopy, newSliceToCopy, updatePtr)

	if c.MergeCommit == nil {
		c.MergeCommit = &core.GhRepoPRCommitConfig{}
	}

	updateGhRepoPRCommitConfigHelper(c.MergeCommit, stringToCopy, newSliceToCopy, updatePtr)

	if c.SquashCommit == nil {
		c.SquashCommit = &core.GhRepoPRCommitConfig{}
	}

	updateGhRepoPRCommitConfigHelper(c.SquashCommit, stringToCopy, newSliceToCopy, updatePtr)

	if c.Branch == nil {
		c.Branch = &core.GhRepoPRBranchConfig{}
	}

	updateGhRepoPRBranchConfigHelper(c.Branch, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoPullRequestConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoPullRequestConfig{},
		func(to, from *core.GhRepoPullRequestConfig) {
			to.Merge(from)
		},
		updateGhRepoPullRequestConfigHelper,
	)
}

func TestGhRepoPullRequestConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := GetFullConfig(0).PullRequests
	toWithNilSlicesAndStruct.MergeStrategy = nil
	toWithNilSlicesAndStruct.MergeCommit = nil
	toWithNilSlicesAndStruct.SquashCommit = nil
	toWithNilSlicesAndStruct.Branch = nil

	cases := map[string]struct {
		value    *core.GhRepoPullRequestConfig
		from     *core.GhRepoPullRequestConfig
		expected *core.GhRepoPullRequestConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->MergeStrategy

func updateGhRepoPRMergeStrategyConfigHelper(c *core.GhRepoPRMergeStrategyConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.AllowMerge, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowRebase, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowSquash, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.AllowAutoMerge, stringToCopy, updatePtr)
}

func TestGhRepoPRMergeStrategyConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoPRMergeStrategyConfig{},
		func(to, from *core.GhRepoPRMergeStrategyConfig) {
			to.Merge(from)
		},
		updateGhRepoPRMergeStrategyConfigHelper,
	)
}

func TestGhRepoPRMergeStrategyConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoPRMergeStrategyConfig
		from     *core.GhRepoPRMergeStrategyConfig
		expected *core.GhRepoPRMergeStrategyConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->MergeCommit and Repo->PullRequests->SquashCommit

func updateGhRepoPRCommitConfigHelper(c *core.GhRepoPRCommitConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Title, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.Message, stringToCopy, updatePtr)
}

func TestGhRepoPRCommitConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoPRCommitConfig{},
		func(to, from *core.GhRepoPRCommitConfig) {
			to.Merge(from)
		},
		updateGhRepoPRCommitConfigHelper,
	)
}

func TestGhRepoPRCommitConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoPRCommitConfig
		from     *core.GhRepoPRCommitConfig
		expected *core.GhRepoPRCommitConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->PullRequests->Branch

func updateGhRepoPRBranchConfigHelper(c *core.GhRepoPRBranchConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.SuggestUpdate, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.DeleteOnMerge, stringToCopy, updatePtr)
}

func TestGhRepoPRBranchConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoPRBranchConfig{},
		func(to, from *core.GhRepoPRBranchConfig) {
			to.Merge(from)
		},
		updateGhRepoPRBranchConfigHelper,
	)
}

func TestGhRepoPRBranchConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoPRBranchConfig
		from     *core.GhRepoPRBranchConfig
		expected *core.GhRepoPRBranchConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Security

func updateGhRepoSecurityConfigHelper(c *core.GhRepoSecurityConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.VulnerabilityAlerts, stringToCopy, updatePtr)
}

func TestGhRepoSecurityConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoSecurityConfig{},
		func(to, from *core.GhRepoSecurityConfig) {
			to.Merge(from)
		},
		updateGhRepoSecurityConfigHelper,
	)
}

func TestGhRepoSecurityConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoSecurityConfig
		from     *core.GhRepoSecurityConfig
		expected *core.GhRepoSecurityConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous

func updateGhRepoMiscellaneousConfigHelper(c *core.GhRepoMiscellaneousConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
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
		c.Template = &core.GhRepoTemplateConfig{}
	}

	updateGhRepoTemplateConfigHelper(c.Template, stringToCopy, newSliceToCopy, updatePtr)

	if c.Pages == nil {
		c.Pages = &core.GhRepoPagesConfig{}
	}

	updateGhRepoPagesConfigHelper(c.Pages, stringToCopy, newSliceToCopy, updatePtr)

	if c.FileTemplates == nil {
		c.FileTemplates = &core.GhRepoFileTemplatesConfig{}
	}

	updateGhRepoFileTemplatesConfigHelper(c.FileTemplates, stringToCopy, newSliceToCopy, updatePtr)
}

func TestGhRepoMiscellaneousConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoMiscellaneousConfig{},
		func(to, from *core.GhRepoMiscellaneousConfig) {
			to.Merge(from)
		},
		updateGhRepoMiscellaneousConfigHelper,
	)
}

func TestGhRepoMiscellaneousConfig_Merge_2(t *testing.T) {
	t.Parallel()

	toWithNilSlicesAndStruct := GetFullConfig(0).Miscellaneous
	toWithNilSlicesAndStruct.Topics = nil
	toWithNilSlicesAndStruct.Template = nil
	toWithNilSlicesAndStruct.Pages = nil
	toWithNilSlicesAndStruct.FileTemplates = nil

	cases := map[string]struct {
		value    *core.GhRepoMiscellaneousConfig
		from     *core.GhRepoMiscellaneousConfig
		expected *core.GhRepoMiscellaneousConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->Template

func updateGhRepoTemplateConfigHelper(c *core.GhRepoTemplateConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Source, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.FullClone, stringToCopy, updatePtr)
}

func TestGhRepoTemplateConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoTemplateConfig{},
		func(to, from *core.GhRepoTemplateConfig) {
			to.Merge(from)
		},
		updateGhRepoTemplateConfigHelper,
	)
}

func TestGhRepoTemplateConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoTemplateConfig
		from     *core.GhRepoTemplateConfig
		expected *core.GhRepoTemplateConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->FileTemplates

func updateGhRepoFileTemplatesConfigHelper(c *core.GhRepoFileTemplatesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Gitignore, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.License, stringToCopy, updatePtr)
}

func TestGhRepoFileTemplatesConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoFileTemplatesConfig{},
		func(to, from *core.GhRepoFileTemplatesConfig) {
			to.Merge(from)
		},
		updateGhRepoFileTemplatesConfigHelper,
	)
}

func TestGhRepoFileTemplatesConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoFileTemplatesConfig
		from     *core.GhRepoFileTemplatesConfig
		expected *core.GhRepoFileTemplatesConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Miscellaneous->Pages

func updateGhRepoPagesConfigHelper(c *core.GhRepoPagesConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.Domain, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourcePath, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.SourceBranch, stringToCopy, updatePtr)
}

func TestGhRepoPagesConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoPagesConfig{},
		func(to, from *core.GhRepoPagesConfig) {
			to.Merge(from)
		},
		updateGhRepoPagesConfigHelper,
	)
}

func TestGhRepoPagesConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoPagesConfig
		from     *core.GhRepoPagesConfig
		expected *core.GhRepoPagesConfig
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
				t.Parallel()
				tc.value.Merge(tc.from)

				if diff := cmp.Diff(tc.expected, tc.value); diff != "" {
					t.Errorf("Config mismatch (-want +got):\n%s", diff)
				}
			},
		)
	}
}

// ==> Repo->Terraform

func updateGhRepoTerraformConfigHelper(c *core.GhRepoTerraformConfig, stringToCopy *string, newSliceToCopy *[]string, updatePtr bool) {
	updateStringPtrHelper(&c.ArchiveOnDestroy, stringToCopy, updatePtr)
	updateStringPtrHelper(&c.IgnoreVulnerabilityAlertsDuringRead, stringToCopy, updatePtr)
}

func TestGhRepoTerraformConfig_Merge(t *testing.T) {
	t.Parallel()
	// Ensure updating from afterward doesn't affect result of merge and vice versa
	ensureMergeWithoutOverflowBetweenToAndFrom(
		t,
		core.GhRepoTerraformConfig{},
		func(to, from *core.GhRepoTerraformConfig) {
			to.Merge(from)
		},
		updateGhRepoTerraformConfigHelper,
	)
}

func TestGhRepoTerraformConfig_Merge_2(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value    *core.GhRepoTerraformConfig
		from     *core.GhRepoTerraformConfig
		expected *core.GhRepoTerraformConfig
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
				t.Parallel()
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
	t.Helper()

	initialToStringValue := "initial_to_value"
	initialToSliceValue := []string{"initial_to_slice_value1", "initial_to_slice_value2"}
	initialFromStringValue := "initial_from_value"
	initialFromSliceValue := []string{"initial_from_slice_value1", "initial_from_slice_value2"}

	initialFromStringValueCopy := fmt.Sprintf("%s", initialFromStringValue)                                                         //nolint:perfsprint // Because :p
	initialFromSliceValueCopy := []string{fmt.Sprintf("%s", initialFromSliceValue[0]), fmt.Sprintf("%s", initialFromSliceValue[1])} //nolint:perfsprint // Because :p

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
	expectedToSliceAfterMerge := append(initialToSliceValue, initialFromSliceValue...) //nolint:gocritic //expected to be merged into another slice
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
	t.Helper()

	updatedToStringValue := "updated_to_value"
	updatedToSliceValue := []string{"updated_to_slice_value1", "updated_to_slice_value2"}
	updatedToStringPtrValue := "updated_to_ptr_value"
	updatedToSlicePtrValue := []string{"updated_to_slice_ptr_value1", "updated_to_slice_ptr_value2"}

	updatedFromStringValue := "updated_from_value"
	updatedFromSliceValue := []string{"updated_from_slice_value1", "updated_from_slice_value2"}
	updatedFromStringPtrValue := "updated_from_ptr_value"
	updatedFromSlicePtrValue := []string{"updated_from_slice_ptr_value1", "updated_from_slice_ptr_value2"}

	updatedToStringValueCopy := fmt.Sprintf("%s", updatedToStringValue)     //nolint:perfsprint // Because :p
	updatedFromStringValueCopy := fmt.Sprintf("%s", updatedFromStringValue) //nolint:perfsprint // Because :p

	updatedToSliceValueCopy := []string{fmt.Sprintf("%s", updatedToSliceValue[0]), fmt.Sprintf("%s", updatedToSliceValue[1])}       //nolint:perfsprint // Because :p
	updatedFromSliceValueCopy := []string{fmt.Sprintf("%s", updatedFromSliceValue[0]), fmt.Sprintf("%s", updatedFromSliceValue[1])} //nolint:perfsprint // Because :p

	updatedToStringPtrValueCopy := fmt.Sprintf("%s", updatedToStringPtrValue) //nolint:perfsprint // Because :p
	updatedToStringPtrCopy := &updatedToStringPtrValueCopy
	updatedFromStringPtrValueCopy := fmt.Sprintf("%s", updatedFromStringPtrValue) //nolint:perfsprint // Because :p
	updatedFromStringPtrCopy := &updatedFromStringPtrValueCopy

	updatedToSlicePtrValueCopy := []string{fmt.Sprintf("%s", updatedToSlicePtrValue[0]), fmt.Sprintf("%s", updatedToSlicePtrValue[1])} //nolint:perfsprint // Because :p
	updatedToSlicePtrCopy := &updatedToSlicePtrValueCopy
	updatedFromSlicePtrValueCopy := []string{fmt.Sprintf("%s", updatedFromSlicePtrValue[0]), fmt.Sprintf("%s", updatedFromSlicePtrValue[1])} //nolint:perfsprint // Because :p
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
	newString := fmt.Sprintf("%s", *stringToCopy) //nolint:perfsprint //Be sure to copy the value rather than assigning the pointer !

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
