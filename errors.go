package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/yoanm/go-github-tf/core"
)

var (
	errDuringWorkspaceLoading = errors.New("error during workspace loading")
	errDuringConfigsLoading   = errors.New("error during configs loading")
	errDuringTemplateLoading  = errors.New("error during templates loading")

	errInputDirectoryDoesntExist = errors.New("input directory doesn't exist")
	errRepositoryAlreadyImported = errors.New("repository already imported")

	errPathIsNotADirectory = errors.New("path is not a directory")
)

func workspaceLoadingError(errList []error) error {
	const separator = "\n"

	return fmt.Errorf("%w:%s%w", errDuringWorkspaceLoading, separator, core.JoinErrors(errList, separator))
}

func configDirectoryLoadingError(errList []error) error {
	const separator = "\n\t - "

	return fmt.Errorf("%w:%s%w", errDuringConfigsLoading, separator, core.JoinErrors(errList, separator))
}

func templateLoadingError(errList []error) error {
	const separator = "\n\t - "

	return fmt.Errorf("%w:%s%w", errDuringTemplateLoading, separator, core.JoinErrors(errList, separator))
}

func inputDirectoryDoesntExistError(path string) error {
	return fmt.Errorf("%w: %s", errInputDirectoryDoesntExist, path)
}

func alreadyImportedRepositoryError(repoName string, filepathList []string) error {
	sort.Strings(filepathList)

	return fmt.Errorf(
		"%w: %q imported by %s",
		errRepositoryAlreadyImported,
		repoName,
		strings.Join(filepathList, ", "),
	)
}

func pathIsNotADirectoryError(path string) error {
	return fmt.Errorf("%w: %s", errPathIsNotADirectory, path)
}
