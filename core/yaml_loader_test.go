package core_test

import (
	"fmt"
	"testing"

	"github.com/yoanm/github-tf/core"
)

func TestLoadRepositoryFromFile(t *testing.T) {
	t.Parallel()

	anchorDir := "testdata/yaml-anchors"
	core.YamlAnchorDirectory = &anchorDir

	name := "my-repo-name-anchor"
	desc := "my-repo-description-anchor"

	cases := map[string]struct {
		filename string
		expected *core.GhRepoConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/repos/repo.unexpected-property.yml",
			nil,
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/repos/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo.full.yml",
			GetFullConfig(1),
			nil,
		},
		"Working with anchors": {
			"testdata/repo.with-anchor.yml",
			&core.GhRepoConfig{
				&name, nil, nil, &desc, nil, nil, nil,
				nil, nil, nil, nil,
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
				conf, err := core.LoadRepositoryFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadRepositoriesFromFile(t *testing.T) {
	t.Parallel()

	anchorDir := "testdata/yaml-anchors"
	core.YamlAnchorDirectory = &anchorDir

	cases := map[string]struct {
		filename string
		expected []*core.GhRepoConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/repos/repos.unexpected-property.yml",
			nil,
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/repos/repos.unexpected-property.yml: /0/unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repos.full.yml",
			[]*core.GhRepoConfig{GetFullConfig(1), GetFullConfig(2)},
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
				conf, err := core.LoadRepositoriesFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadRepositoryTemplateFromFile(t *testing.T) {
	t.Parallel()

	full := GetFullConfig(1)
	// Template can't have a Name
	full.Name = nil
	cases := map[string]struct {
		filename string
		expected *core.GhRepoConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/templates/repo.unexpected-property.yml",
			nil,
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/templates/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo-template.full.yml",
			full,
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
				conf, err := core.LoadRepositoryTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadBranchTemplateFromFile(t *testing.T) {
	t.Parallel()

	fullRepo := GetFullConfig(1)
	full := (*fullRepo.Branches)["feature/branch1"]
	cases := map[string]struct {
		filename string
		expected *core.GhBranchConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/templates/branch.unexpected-property.yml",
			nil,
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/templates/branch.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-template.full.yml",
			full,
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
				conf, err := core.LoadBranchTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadBranchProtectionTemplateFromFile(t *testing.T) {
	t.Parallel()

	fullRepo := GetFullConfig(1)
	full := (*fullRepo.BranchProtections)[0]
	cases := map[string]struct {
		filename string
		expected *core.GhBranchProtectionConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml",
			nil,
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
			full,
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
				conf, err := core.LoadBranchProtectionTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoConfigFromFile(t *testing.T) {
	t.Parallel()

	full := GetFullConfig(1)
	repoName := "repo-name"
	cases := map[string]struct {
		filename string
		expected *core.GhRepoConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/repos/repo.unexpected-property.yml",
			&core.GhRepoConfig{
				&repoName, nil, nil, nil, nil, nil, nil,
				nil, nil, nil, nil,
			},
			nil,
		},
		"Working": {
			"testdata/repo.full.yml",
			full,
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
				conf, err := core.LoadGhRepoConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoConfigListFromFile(t *testing.T) {
	t.Parallel()

	full1, full2 := GetFullConfig(1), GetFullConfig(2)
	repoName := "repo-name"
	cases := map[string]struct {
		filename string
		expected []*core.GhRepoConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/repos/repos.unexpected-property.yml",
			[]*core.GhRepoConfig{
				{
					&repoName, nil, nil, nil, nil, nil, nil,
					nil, nil, nil, nil,
				},
			},
			nil,
		},
		"Working": {
			"testdata/repos.full.yml",
			[]*core.GhRepoConfig{full1, full2},
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
				conf, err := core.LoadGhRepoConfigListFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoBranchConfigFromFile(t *testing.T) {
	t.Parallel()

	fullRepo := GetFullConfig(1)
	full := (*fullRepo.Branches)["feature/branch1"]
	sourceBranch := "branch1-source-branch1"
	cases := map[string]struct {
		filename string
		expected *core.GhBranchConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/templates/branch.unexpected-property.yml",
			&core.GhBranchConfig{SourceBranch: &sourceBranch},
			nil,
		},
		"Working": {
			"testdata/branch-template.full.yml",
			full,
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
				conf, err := core.LoadGhRepoBranchConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoBranchProtectionConfigFromFile(t *testing.T) {
	t.Parallel()

	fullRepo := GetFullConfig(1)
	full := (*fullRepo.BranchProtections)[0]
	pattern := "master"
	cases := map[string]struct {
		filename string
		expected *core.GhBranchProtectionConfig
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			nil,
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			nil,
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unexpected property": {
			"testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml",
			&core.GhBranchProtectionConfig{Pattern: &pattern},
			nil,
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
			full,
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
				conf, err := core.LoadGhRepoBranchProtectionConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}
