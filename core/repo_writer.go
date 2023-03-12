package core

import (
	"sort"

	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/yoanm/go-gh2tf/ghbranch"
	"github.com/yoanm/go-gh2tf/ghbranchdefault"
	"github.com/yoanm/go-gh2tf/ghbranchprotect"
	"github.com/yoanm/go-gh2tf/ghrepository"
	"github.com/yoanm/go-tfsig"
)

/** Public **/

func NewHclRepository(repoTfId string, c *GhRepoConfig, valGen tfsig.ValueGenerator) *hclwrite.File {
	hclFile := hclwrite.NewEmptyFile()

	appendRepositoryResource(hclFile.Body(), c, valGen, repoTfId)
	appendBranchDefaultResources(hclFile.Body(), c, valGen, repoTfId)
	appendBranchResources(hclFile.Body(), c, valGen, repoTfId)
	appendBranchProtectionResourceContent(hclFile.Body(), c, valGen, repoTfId)

	return hclFile
}

/** Private **/

func appendRepositoryResource(body *hclwrite.Body, c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) {
	tfsig.AppendBlockIfNotNil(
		body,
		ghrepository.New(
			MapToRepositoryRes(c, valGen, repoTfId),
		),
	)
}

func appendBranchDefaultResources(body *hclwrite.Body, c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) {
	tfsig.AppendNewLineAndBlockIfNotNil(
		body,
		ghbranchdefault.New(
			// /!\ use LinkToRepository, so underlying repository will have to be created before creating default branch
			MapToDefaultBranchRes(c.DefaultBranch, valGen, c, repoTfId, LinkToRepository, LinkToBranch),
		),
	)
}

func appendBranchResources(body *hclwrite.Body, c *GhRepoConfig, valGen tfsig.ValueGenerator, repoTfId string) {
	if c.Branches != nil {
		hasDefaultBranch := c.DefaultBranch != nil && c.DefaultBranch.Name != nil
		// sort branches to always get a predictable output (for tests mostly)
		keys := make([]string, 0, len(*c.Branches))
		for k := range *c.Branches {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			v := (*c.Branches)[k]

			appendBranchResourceSignature(
				body,
				c,
				// /!\ use LinkToRepository, so underlying repository will have to be created before creating the branch
				// /!\ use LinkToBranch, so if a source branch is configured, it will have to be created before creating current one
				ghbranch.NewSignature(MapToBranchRes(k, v, valGen, c, repoTfId, LinkToRepository, LinkToBranch)),
				hasDefaultBranch,
				k,
				v,
			)
		}
	}
}

func appendBranchResourceSignature(
	body *hclwrite.Body,
	c *GhRepoConfig,
	sig *tfsig.BlockSignature,
	hasDefaultBranch bool,
	k string,
	v *GhBranchConfig,
) {
	if sig != nil {
		// In case
		//  - it's the default branch config
		ignoreSourceBranch := hasDefaultBranch && *c.DefaultBranch.Name == k
		//  - or source_branch is the default branch
		ignoreSourceBranch = ignoreSourceBranch || (hasDefaultBranch && v.SourceBranch != nil &&
			*c.DefaultBranch.Name == *v.SourceBranch)
		//  - or no source branch configured (which means the same as current default branch)
		ignoreSourceBranch = ignoreSourceBranch || v.SourceBranch == nil
		// => append ignore_changes directive on source_branch
		// It's useful mostly when switching default branch to another one from outside of terraform,
		// or from terraform when following those steps:
		// Step 1 - add both current (if not already there) and new branch config(s) with empty
		// 				configuration (brName: {})
		//        - run github-tf to re-generate terraform config
		// Step 3 - if current default branch had no config previously:
		// 				run terraform import github_branch.${repo}-${oldBranch} ${repo}:${oldBranch}
		//        - run terraform apply to add the new branch
		// Step 4 - switch default branch name
		//        - If old default branch wasn't protected and you don't want to keep old default branch config
		//        		(branches->BRANCH_NAME), jump to Step 5
		//        - run github-tf to re-generate terraform config
		//        - run terraform apply
		// Step 5 [Optional, but usually] - remove old branch config
		// 	      - If you don't want to keep new default branch config (branches->BRANCH_NAME), jump to step 6
		//        - run github-tf to re-generate terraform config
		//        - run terraform apply
		// Step 6 [Optional] - Remove new branch config
		//        - run github-tf to re-generate terraform config
		//        - run terraform state rm github_branch.${repo}-${newBranch}
		//        - in case step 5 has been done: run terraform apply (will delete old branch)
		if ignoreSourceBranch {
			sig.Lifecycle(
				tfsig.LifecycleConfig{
					CreateBeforeDestroy: nil,
					PreventDestroy:      nil,
					IgnoreChanges:       []string{"source_branch"},
					ReplaceTriggeredBy:  nil,
					Precondition:        nil,
					Postcondition:       nil,
				},
			)
		}

		tfsig.AppendNewLineAndBlockIfNotNil(body, sig.Build())
	}
}

func appendBranchProtectionResourceContent(
	body *hclwrite.Body,
	c *GhRepoConfig,
	valGen tfsig.ValueGenerator,
	repoTfId string,
) {
	tfsig.AppendNewLineAndBlockIfNotNil(
		body,
		ghbranchprotect.New(
			// /!\ use LinkToRepository, so underlying repository will have to be created before
			// creating the branch protection
			// /!\ use LinkToBranch, to explicitly bind the protection to the branch, so if something
			// change on default_branch resource,
			// the protection will be impacted only if the change on default_branch resource succeeded
			MapDefaultBranchToBranchProtectionRes(c.DefaultBranch, valGen, c, repoTfId, LinkToRepository, LinkToBranch),
		),
	)

	if c.Branches != nil {
		// sort branches to always get a predictable output (for tests mostly)
		keys := make([]string, 0, len(*c.Branches))

		for k := range *c.Branches {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			v := (*c.Branches)[k]

			tfsig.AppendNewLineAndBlockIfNotNil(
				body,
				ghbranchprotect.New(
					// /!\ use LinkToRepository, so underlying repository will have to be created before
					// creating the branch protection
					// /!\ do not use LinkToBranch, else branch protection will be created only after branch is created
					// (which is useless and can be done in parallel)
					MapBranchToBranchProtectionRes(k, v, valGen, c, repoTfId, LinkToRepository),
				),
			)
		}
	}

	if c.BranchProtections != nil {
		for _, branchProtectionConfig := range *c.BranchProtections {
			tfsig.AppendNewLineAndBlockIfNotNil(
				body,
				ghbranchprotect.New(
					// /!\ use LinkToRepository, so underlying repository will have to be created before
					// creating the branch protection
					// /!\ do not use LinkToBranch, else branch protection will be created only after
					// branch is created (which is useless and can be done in parallel) and in many cases related
					// branch config doesn't exist anyway (else it's simpler to move the protection config under
					// 'protection' attribute of the related branch)
					MapToBranchProtectionRes(branchProtectionConfig, valGen, c, repoTfId, LinkToRepository),
				),
			)
		}
	}
}
