package core

import (
	"fmt"
	"testing"
)

func TestLoadRepositoryFromFile(t *testing.T) {
	anchorDir := "testdata/yaml-anchors"
	YamlAnchorDirectory = &anchorDir

	name := "my-repo-name-anchor"
	desc := "my-repo-description-anchor"

	cases := map[string]struct {
		filename string
		expected *GhRepoConfig
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
			fmt.Errorf("file testdata/invalid-config-files/repos/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo.full.yml",
			GetFullConfig(1),
			nil,
		},
		"Working with anchors": {
			"testdata/repo.with-anchor.yml",
			&GhRepoConfig{
				&name, nil, nil, &desc, nil, nil, nil,
				nil, nil, nil, nil,
			},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadRepositoryFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadRepositoriesFromFile(t *testing.T) {
	anchorDir := "testdata/yaml-anchors"
	YamlAnchorDirectory = &anchorDir

	cases := map[string]struct {
		filename string
		expected []*GhRepoConfig
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
			fmt.Errorf("file testdata/invalid-config-files/repos/repos.unexpected-property.yml: /0/unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repos.full.yml",
			[]*GhRepoConfig{GetFullConfig(1), GetFullConfig(2)},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadRepositoriesFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadRepositoryTemplateFromFile(t *testing.T) {
	full := GetFullConfig(1)
	// Template can't have a Name
	full.Name = nil
	cases := map[string]struct {
		filename string
		expected *GhRepoConfig
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
			fmt.Errorf("file testdata/invalid-config-files/templates/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo-template.full.yml",
			full,
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadRepositoryTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadBranchTemplateFromFile(t *testing.T) {
	fullRepo := GetFullConfig(1)
	full := (*fullRepo.Branches)["feature/branch1"]
	cases := map[string]struct {
		filename string
		expected *GhBranchConfig
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
			fmt.Errorf("file testdata/invalid-config-files/templates/branch.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-template.full.yml",
			full,
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadBranchTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadBranchProtectionTemplateFromFile(t *testing.T) {
	fullRepo := GetFullConfig(1)
	full := (*fullRepo.BranchProtections)[0]
	cases := map[string]struct {
		filename string
		expected *GhBranchProtectionConfig
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
			fmt.Errorf("file testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
			full,
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadBranchProtectionTemplateFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoConfigFromFile(t *testing.T) {
	full := GetFullConfig(1)
	repoName := "repo-name"
	cases := map[string]struct {
		filename string
		expected *GhRepoConfig
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
			&GhRepoConfig{
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
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadGhRepoConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoConfigListFromFile(t *testing.T) {
	full1, full2 := GetFullConfig(1), GetFullConfig(2)
	repoName := "repo-name"
	cases := map[string]struct {
		filename string
		expected []*GhRepoConfig
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
			[]*GhRepoConfig{
				&GhRepoConfig{
					&repoName, nil, nil, nil, nil, nil, nil,
					nil, nil, nil, nil,
				},
			},
			nil,
		},
		"Working": {
			"testdata/repos.full.yml",
			[]*GhRepoConfig{full1, full2},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadGhRepoConfigListFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoBranchConfigFromFile(t *testing.T) {
	fullRepo := GetFullConfig(1)
	full := (*fullRepo.Branches)["feature/branch1"]
	sourceBranch := "branch1-source-branch1"
	cases := map[string]struct {
		filename string
		expected *GhBranchConfig
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
			&GhBranchConfig{SourceBranch: &sourceBranch},
			nil,
		},
		"Working": {
			"testdata/branch-template.full.yml",
			full,
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadGhRepoBranchConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}

func TestLoadGhRepoBranchProtectionConfigFromFile(t *testing.T) {
	fullRepo := GetFullConfig(1)
	full := (*fullRepo.BranchProtections)[0]
	pattern := "master"
	cases := map[string]struct {
		filename string
		expected *GhBranchProtectionConfig
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
			&GhBranchProtectionConfig{Pattern: &pattern},
			nil,
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
			full,
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				conf, err := LoadGhRepoBranchProtectionConfigFromFile(tc.filename)
				EnsureConfigMatching(t, tc.expected, conf, tc.error, err)
			},
		)
	}
}
