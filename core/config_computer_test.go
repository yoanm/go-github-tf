package core

import (
	"fmt"
	"testing"

	differ "github.com/andreyvit/diff"
	"github.com/google/go-cmp/cmp"
)

func TestComputeConfig(t *testing.T) {
	aName := "a_name"
	cases := map[string]struct {
		value    *Config
		expected *Config
		error    error
	}{
		"nil": {
			nil,
			nil,
			nil,
		},
		"empty": {
			&Config{Templates: nil, Repos: nil},
			NewConfig(),
			nil,
		},
		"base": {
			&Config{
				Templates: nil,
				Repos: []*GhRepoConfig{
					{
						&aName, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			&Config{
				Templates: &TemplatesConfig{
					Repos:             map[string]*GhRepoConfig{},
					Branches:          map[string]*GhBranchConfig{},
					BranchProtections: map[string]*GhBranchProtectionConfig{},
				},
				Repos: []*GhRepoConfig{
					{
						&aName, nil, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			nil,
		},
		"Repo without name": {
			&Config{
				Templates: nil,
				Repos: []*GhRepoConfig{
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
			fmt.Errorf("error during computation:\n\t - repository name is missing for repo #0\n\t - repository name is missing for repo #1"),
		},
		"Underlying computation error": {
			&Config{
				Templates: nil,
				Repos: []*GhRepoConfig{
					{
						&aName, &[]string{aName}, nil, nil, nil, nil,
						nil, nil, nil, nil, nil,
					},
				},
			},
			nil,
			fmt.Errorf("error during computation:\n\t - repository a_name: unable to load repository template, no template available"),
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				actual, err := ComputeConfig(tc.value)
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
