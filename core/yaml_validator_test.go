package core_test

import (
	"fmt"
	"testing"

	"github.com/yoanm/github-tf/core"
)

func TestValidateRepositoryConfig(t *testing.T) {
	t.Parallel()

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
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/repos/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo.full.yml",
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
				EnsureErrorMatching(t, tc.error, core.ValidateRepositoryConfig(tc.filename))
			},
		)
	}
}

func TestValidateRepositoryConfigs(t *testing.T) {
	t.Parallel()

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
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/repos/repos.unexpected-property.yml: /0/unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repos.full.yml",
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
				EnsureErrorMatching(t, tc.error, core.ValidateRepositoryConfigs(tc.filename))
			},
		)
	}
}

func TestValidateRepositoryTemplateConfig(t *testing.T) {
	t.Parallel()

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
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/templates/repo.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/repo-template.full.yml",
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
				EnsureErrorMatching(t, tc.error, core.ValidateRepositoryTemplateConfig(tc.filename))
			},
		)
	}
}

func TestValidateBranchProtectionTemplateConfig(t *testing.T) {
	t.Parallel()

	// Reset YamlAnchorDirectory, so it's certain to cover getYamlValidatorDecoderOptions default return
	core.YamlAnchorDirectory = nil

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
			fmt.Errorf("schema validation error: file testdata/invalid-config-files/templates/branch-protection.unexpected-property.yml: /unexpected-property not allowed"),
		},
		"Working": {
			"testdata/branch-protection-template.full.yml",
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
				EnsureErrorMatching(t, tc.error, core.ValidateBranchProtectionTemplateConfig(tc.filename))
			},
		)
	}
}
