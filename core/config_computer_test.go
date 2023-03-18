package core_test

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"

	"github.com/yoanm/go-github-tf/core"
)

func TestComputeConfig(t *testing.T) {
	t.Parallel()

	aName := "a_name"
	cases := map[string]struct {
		value    *core.Config
		expected *core.Config
		error    error
	}{
		"nil": {
			nil,
			nil,
			nil,
		},
		"empty": {
			&core.Config{Templates: nil, Repos: nil},
			core.NewConfig(),
			nil,
		},
		"base": {
			&core.Config{
				Templates: nil,
				Repos: []*core.GhRepoConfig{
					{
						&aName, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			&core.Config{
				Templates: &core.TemplatesConfig{
					Repos:             map[string]*core.GhRepoConfig{},
					Branches:          map[string]*core.GhBranchConfig{},
					BranchProtections: map[string]*core.GhBranchProtectionConfig{},
				},
				Repos: []*core.GhRepoConfig{
					{
						&aName, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			nil,
		},
		"Repo without name": {
			&core.Config{
				Templates: nil,
				Repos: []*core.GhRepoConfig{
					{
						nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
					{
						nil, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			nil,
			fmt.Errorf("error during computation:\n\t - repo #0: repository name is mandatory\n\t - repo #1: repository name is mandatory"),
		},
		"Underlying computation error": {
			&core.Config{
				Templates: nil,
				Repos: []*core.GhRepoConfig{
					{
						&aName, &[]string{aName}, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			nil,
			fmt.Errorf("error during computation:\n\t - repository a_name: repository template not found as none available"),
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel tests
		tc := tc         // Reinit var for parallel tests

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				actual, err := core.ComputeConfig(tc.value)
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
