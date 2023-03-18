package core_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/yoanm/github-tf/core"
	"github.com/yoanm/go-tfsig/testutils"
)

func TestGenerateHclRepoFiles(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		value         []*core.GhRepoConfig
		expectedFiles map[string]string
	}{
		"Full": {
			[]*core.GhRepoConfig{GetFullConfig(1), GetFullConfig(2)},
			map[string]string{
				"repo.repo1.tf": "repo1.full",
				"repo.repo2.tf": "repo2.full",
			},
		},
		"nil": {
			nil,
			nil,
		},
		"empty": {
			[]*core.GhRepoConfig{},
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
				files, err := core.GenerateHclRepoFiles(tc.value)
				if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else {
					for fname, goldenfile := range tc.expectedFiles {
						tffile, exists := files[fname]
						if !exists {
							t.Errorf("Case %q: expected file %s doesn't exist !", tcname, fname)
						} else {
							if err2 := testutils.EnsureFileEqualsGoldenFile(tffile, goldenfile); err2 != nil {
								t.Errorf("Case %q file %s: %v", tcname, fname, err2)
							}
						}
					}
					if !t.Failed() && len(files) != len(tc.expectedFiles) {
						t.Errorf("Case %q: expected %d files, got %d", tcname, len(tc.expectedFiles), len(files))
					}
				}
			},
		)
	}
}

func TestWriteTerraformFiles(t *testing.T) {
	t.Parallel()

	file1 := hclwrite.NewEmptyFile()
	file1.Body().AppendBlock(hclwrite.NewBlock("type1", []string{"label1"}))
	file2 := hclwrite.NewEmptyFile()
	file2.Body().AppendBlock(hclwrite.NewBlock("type2", []string{"label2"}))

	cases := map[string]struct {
		value         map[string]*hclwrite.File
		expectedFiles map[string]string
	}{
		"Full": {
			map[string]*hclwrite.File{
				"repo.repo1.tf": file1,
				"repo.repo2.tf": file2,
			},
			map[string]string{
				"repo.repo1.tf": "type1 \"label1\" {\n}\n",
				"repo.repo2.tf": "type2 \"label2\" {\n}\n",
			},
		},
		"Nil": {
			nil,
			nil,
		},
		"Empty": {
			map[string]*hclwrite.File{},
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
				root := os.TempDir()
				err := core.WriteTerraformFiles(root, tc.value)
				if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else {
					for fname, expected := range tc.expectedFiles {
						_, exists := tc.value[fname]
						if !exists {
							t.Errorf("Case %q: expected file %s doesn't exist !", tcname, fname)
						} else {
							actual, err2 := os.ReadFile(path.Join(root, fname))
							if err2 != nil {
								t.Errorf("Case %q: %s", tcname, err2)
							} else if string(actual) != expected {
								t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, diff.LineDiff(expected, string(actual)))
							}
						}
					}
				}
			},
		)
	}
}

func TestWriteTerraformFiles_onError(t *testing.T) {
	t.Parallel()

	file1 := hclwrite.NewEmptyFile()
	file1.Body().AppendBlock(hclwrite.NewBlock("type1", []string{"label1"}))

	notWritableDir := path.Join(os.TempDir(), "not_writable_dir")

	err := os.MkdirAll(notWritableDir, os.FileMode(0))
	if err != nil {
		t.Fatalf(err.Error())
	}

	cases := map[string]struct {
		root  string
		value map[string]*hclwrite.File
		error error
	}{
		"Unknown dir": {
			"/an_unknown_dir/somewhere",
			map[string]*hclwrite.File{
				"repo.repo1.tf": file1,
			},
			fmt.Errorf("error while writing terraform files:\n\t - workspace path doesn't exist: /an_unknown_dir/somewhere"),
		},
		"Unable to write": {
			notWritableDir,
			map[string]*hclwrite.File{
				"repo.repo1.tf": file1,
			},
			fmt.Errorf("error while writing terraform files:\n\t - open %s/repo.repo1.tf: permission denied", notWritableDir),
		},
	}

	for tcname, tc := range cases {
		tcname := tcname // Reinit var for parallel test
		tc := tc         // Reinit var for parallel test

		t.Run(
			tcname,
			func(t *testing.T) {
				t.Parallel()
				if err2 := core.WriteTerraformFiles(tc.root, tc.value); err2 == nil {
					t.Errorf("Case %q: expected an error but everything went well", tcname)
				} else if err2.Error() != tc.error.Error() {
					t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, diff.LineDiff(tc.error.Error(), err2.Error()))
				}
			},
		)
	}
}
