package core

import (
	"fmt"
	"testing"
)

func TestValidateRepositoryConfig(t *testing.T) {
	cases := map[string]struct {
		filename string
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unwanted property": {
			"testdata/invalid-config-files/repos/repo.unexpected-property.yml",
			fmt.Errorf("file testdata/invalid-config-files/repos/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo.full.yml",
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				EnsureErrorMatching(t, tc.error, ValidateRepositoryConfig(tc.filename))
			},
		)
	}
}

func TestValidateRepositoryConfigs(t *testing.T) {
	cases := map[string]struct {
		filename string
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unwanted property": {
			"testdata/invalid-config-files/repos/repos.unexpected-property.yml",
			fmt.Errorf("file testdata/invalid-config-files/repos/repos.unexpected-property.yml: /0/unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repos.full.yml",
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				EnsureErrorMatching(t, tc.error, ValidateRepositoryConfigs(tc.filename))
			},
		)
	}
}

func TestValidateRepositoryTemplateConfig(t *testing.T) {
	full := GetFullConfig(0)
	// Template can't have a Name
	full.Name = nil
	cases := map[string]struct {
		filename string
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unwanted property": {
			"testdata/invalid-config-files/templates/repo.unexpected-property.yml",
			fmt.Errorf("file testdata/invalid-config-files/templates/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo-template.full.yml",
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				EnsureErrorMatching(t, tc.error, ValidateRepositoryTemplateConfig(tc.filename))
			},
		)
	}
}

func TestValidateBranchProtectionTemplateConfig(t *testing.T) {
	// Reset YamlAnchorDirectory, so it's certain to cover getYamlValidatorDecoderOptions default return
	YamlAnchorDirectory = nil

	cases := map[string]struct {
		filename string
		error    error
	}{
		"Not found file": {
			"an_unknown_file",
			fmt.Errorf("open an_unknown_file: no such file or directory"),
		},
		"Empty": {
			"testdata/invalid-config-files/empty.yml",
			fmt.Errorf("file testdata/invalid-config-files/empty.yml: EOF"),
		},
		"Unwanted property": {
			"testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml",
			fmt.Errorf("file testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				EnsureErrorMatching(t, tc.error, ValidateBranchProtectionTemplateConfig(tc.filename))
			},
		)
	}
}
