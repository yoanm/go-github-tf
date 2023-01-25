package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/hashicorp/hcl/v2/hclwrite"

	"github.com/yoanm/go-tfsig/testutils"
)

func TestGenerateHclRepoFiles(t *testing.T) {
	cases := map[string]struct {
		value         []*GhRepoConfig
		expectedFiles map[string]string
	}{
		"Full": {
			[]*GhRepoConfig{GetFullConfig(1), GetFullConfig(2)},
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
			[]*GhRepoConfig{},
			nil,
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				files, err := GenerateHclRepoFiles(tc.value)
				if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else {
					for fname, goldenfile := range tc.expectedFiles {
						tffile, exists := files[fname]
						if !exists {
							t.Errorf("Case %q: expected file %s doesn't exist !", tcname, fname)
						} else {
							if err := testutils.EnsureFileEqualsGoldenFile(tffile, goldenfile); err != nil {
								t.Errorf("Case %q file %s: %v", tcname, fname, err)
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
		t.Run(
			tcname,
			func(t *testing.T) {
				root := os.TempDir()
				err := WriteTerraformFiles(root, tc.value)
				if err != nil {
					t.Errorf("Case %q: %s", tcname, err)
				} else {
					for fname, expected := range tc.expectedFiles {
						_, exists := tc.value[fname]
						if !exists {
							t.Errorf("Case %q: expected file %s doesn't exist !", tcname, fname)
						} else {
							actual, err := ioutil.ReadFile(path.Join(root, fname))
							if err != nil {
								t.Errorf("Case %q: %s", tcname, err)
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
			fmt.Errorf("error while writing files:\n\t - open /an_unknown_dir/somewhere: no such file or directory"),
		},
		"Unable to write": {
			notWritableDir,
			map[string]*hclwrite.File{
				"repo.repo1.tf": file1,
			},
			fmt.Errorf("error while writing files:\n\t - open %s/repo.repo1.tf: permission denied", notWritableDir),
		},
	}

	for tcname, tc := range cases {
		t.Run(
			tcname,
			func(t *testing.T) {
				err := WriteTerraformFiles(tc.root, tc.value)
				if err == nil {
					t.Errorf("Case %q: expected an error but everything went well", tcname)
				} else if err.Error() != tc.error.Error() {
					t.Errorf("Case %q:\n- expected\n+ actual\n\n%v", tcname, diff.LineDiff(tc.error.Error(), err.Error()))
				}
			},
		)
	}
}
