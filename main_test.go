package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmdtest"
)

func TestCLIFlags(t *testing.T) {
	ts := configure(t, "testdata/base")
	ts.Run(t, false)
}

func TestCLIWrite_withLoadingErrors(t *testing.T) {
	cases := []string{
		"invalid-config-files",
		"invalid-config-files-2",
		"invalid-templates-files",
		"multiple-invalid-files",
		"multiple-unknown-files",
		"permission-issue",
		"permission-issue-2",
		"invalid-workspace-dir",
		"multiple-config-for-same-repo",
		"multiple-config-for-same-repo-2",
		"invalid-yaml",
		"invalid-yaml-2",
	}
	for _, cname := range cases {
		t.Run(
			cname,
			func(t2 *testing.T) {
				ts := configure(t2, filepath.Join("testdata/write/loading-errors", cname))
				// ts.KeepRootDirs = true
				ts.Run(t2, false)
			},
		)
	}
}

func TestCLIWrite_withComputationErrors(t *testing.T) {
	cases := []string{
		"unknown-template",
		"default-branch-template-without-default-branch",
	}
	for _, cname := range cases {
		t.Run(
			cname,
			func(t2 *testing.T) {
				ts := configure(t2, filepath.Join("testdata/write/computation-errors", cname))
				// ts.KeepRootDirs = true
				ts.Run(t2, false)
			},
		)
	}
}

func TestCLIWrite_withTerraformErrors(t *testing.T) {
	cases := []string{
		"missing-terraform-directory",
		"terraform-directory-as-file",
		"permission-issue",
	}
	for _, cname := range cases {
		t.Run(
			cname,
			func(t2 *testing.T) {
				ts := configure(t2, filepath.Join("testdata/write/terraform-errors", cname))
				// ts.KeepRootDirs = true
				ts.Run(t2, false)
			},
		)
	}
}

func TestCLIWrite_working(t *testing.T) {
	cases := []string{
		"base",
		"full",
		"yml-vs-yaml",
		"with-templates",
		"with-templates-and-anchors",
		"multiple-branch-protection-for-same-pattern",
		"default-branch-branch-protection-template-with-existing-config",
	}
	for _, cname := range cases {
		t.Run(
			cname,
			func(t2 *testing.T) {
				ts := configure(t2, filepath.Join("testdata/write/working", cname))
				// ts.KeepRootDirs = true
				ts.Run(t2, false)
			},
		)
	}
}

func configure(t *testing.T, testdataPath string) *cmdtest.TestSuite {
	ts, err := cmdtest.Read(testdataPath)
	if err != nil {
		t.Fatal(err)
	}

	ts.Commands["github-tf"] = cmdtest.InProcessProgram("github-tf", run)
	ts.Commands["chmod"] = chmodCmd
	ts.Setup = func(rootDir string) error {
		_, testFileName, _, ok := runtime.Caller(0)
		if !ok {
			return fmt.Errorf("failed get real working directory from caller")
		}

		projectRootDir := filepath.Dir(testFileName)
		// fmt.Printf("Project dir %s\n", projectRootDir)
		// fmt.Printf("ROOTDIR %s\n", rootDir)

		// copy {testdataPath}/testdata to ROOTDIR/testdata if it exists
		testdataSourcePath := filepath.Join(projectRootDir, testdataPath, "testdata")
		if _, err = os.Stat(testdataSourcePath); !os.IsNotExist(err) {
			fmt.Printf("Copy testdata %s\n", testdataSourcePath)
			testdataTargetPath := filepath.Join(rootDir, "testdata")
			cmd := exec.Command("cp", "-r", testdataSourcePath, testdataTargetPath)
			if _, err = cmd.Output(); err != nil {
				return fmt.Errorf("Error during testdata copy (%s -> %s): %s", testdataSourcePath, testdataTargetPath, err)
			}
		}

		return nil
	}

	return ts
}

func chmodCmd(args []string, inputFile string) ([]byte, error) {
	if inputFile != "" {
		return nil, fmt.Errorf("input redirection not supported")
	}
	if err := checkPath(args[0]); err != nil {
		return nil, err
	}
	perm, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}
	if err = os.Chmod(args[0], os.FileMode(perm)); err != nil {
		return nil, err
	}
	return nil, nil
}

func checkPath(path string) error {
	if strings.ContainsRune(path, '/') || strings.ContainsRune(path, '\\') {
		return fmt.Errorf("argument must be in the current directory (%q contains '/')", path)
	}
	return nil
}