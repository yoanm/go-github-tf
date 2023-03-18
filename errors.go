package main

import (
	"errors"
	"fmt"

	"github.com/yoanm/github-tf/core"
)

var (
	errDuringWorkspaceLoading = errors.New("error during workspace loading")
	errDuringConfigsLoading   = errors.New("error during configs loading")
	errDuringTemplateLoading  = errors.New("error during templates loading")

	errInputDirectoryDoesntExist = errors.New("input directory doesn't exist")
	errRepositoryAlreadyImported = errors.New("repository already imported")

	errPathIsNotADirectory = errors.New("path is not a directory")
)

func WorkspaceLoadingError(errList []error) error {
	const separator = "\n"

	return fmt.Errorf("%w:%s%w", errDuringWorkspaceLoading, separator, core.JoinErrors(errList, separator))
}

func ConfigDirectoryLoadingError(errList []error) error {
	const separator = "\n\t - "

	return fmt.Errorf("%w:%s%w", errDuringConfigsLoading, separator, core.JoinErrors(errList, separator))
}

func TemplateLoadingError(errList []error) error {
	const separator = "\n\t - "

	return fmt.Errorf("%w:%s%w", errDuringTemplateLoading, separator, core.JoinErrors(errList, separator))
}

func InputDirectoryDoesntExistError(path string) error {
	return fmt.Errorf("%w: %s", errInputDirectoryDoesntExist, path)
}

func AlreadyImportedRepositoryError(fName string, repoName string, firstFName string) error {
	return fmt.Errorf(
		"%w: repository %s imported by %s, but already imported by %s",
		errRepositoryAlreadyImported,
		repoName,
		fName,
		firstFName,
	)
}

func PathIsNotADirectoryError(path string) error {
	return fmt.Errorf("%w: %s", errPathIsNotADirectory, path)
}
