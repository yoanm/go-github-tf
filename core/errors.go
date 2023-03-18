package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var (
	ErrRepositoryNameIsMandatory = errors.New("repository name is mandatory")

	ErrWorkspacePathDoesntExist              = errors.New("workspace path doesn't exist")
	ErrWorkspacePathIsExpectedToBeADirectory = errors.New("workspace path is expected to be a directory")

	ErrNoTemplateAvailable                 = errors.New("not found as none available")
	ErrNoRepositoryTemplateAvailable       = fmt.Errorf("%s template %w", RepositoryTemplateType, ErrNoTemplateAvailable)
	ErrNoBranchTemplateAvailable           = fmt.Errorf("%s template %w", BranchTemplateType, ErrNoTemplateAvailable)
	ErrNoBranchProtectionTemplateAvailable = fmt.Errorf(
		"%s template %w",
		BranchProtectionTemplateType,
		ErrNoTemplateAvailable,
	)

	ErrTemplateNotFound                 = errors.New("not found")
	ErrRepositoryTemplateNotFound       = fmt.Errorf("%s template %w", RepositoryTemplateType, ErrTemplateNotFound)
	ErrBranchTemplateNotFound           = fmt.Errorf("%s template %w", BranchTemplateType, ErrTemplateNotFound)
	ErrBranchProtectionTemplateNotFound = fmt.Errorf("%s template %w", BranchProtectionTemplateType, ErrTemplateNotFound)

	ErrMaxTemplateCount = errors.New("maximum template count reached")
	ErrMaxTemplateDepth = errors.New("maximum template depth reached")

	ErrDuringWriteTerraformFiles = errors.New("error while writing terraform files")
	ErrDuringFileGeneration      = errors.New("error while generating files")
	ErrDuringComputation         = errors.New("error during computation")

	// Json Schema and validation.
	ErrSchemaValidation        = errors.New("schema validation error")
	ErrEmptySchema             = errors.New("empty schema")
	ErrSchemaNotFound          = errors.New("schema not found")
	ErrSchemaIsNil             = errors.New("schema is nil")
	ErrDuringSchemaCompilation = errors.New("error during schema compilation")

	ErrFileError             = errors.New("file")
	ErrBranchError           = errors.New("branch")
	ErrDefaultBranchError    = errors.New("default branch")
	ErrBranchProtectionError = errors.New("branch protection")
)

func BranchError(branch string, err error) error {
	return fmt.Errorf("%w %s: %w", ErrBranchError, branch, err)
}

func DefaultBranchError(err error) error {
	return fmt.Errorf("%w: %w", ErrDefaultBranchError, err)
}

func BranchProtectionError(index int, err error) error {
	return fmt.Errorf("%w #%d: %w", ErrBranchProtectionError, index, err)
}

func FileError(filepath string, err error) error {
	return fmt.Errorf("%w %s: %w", ErrFileError, filepath, err)
}

func FileGenerationError(msgList []string) error {
	return fmt.Errorf("%w:\n\t - %s", ErrDuringFileGeneration, strings.Join(msgList, "\n\t - "))
}

func WorkspacePathDoesntExistError(path string) error {
	return fmt.Errorf("%w: %s", ErrWorkspacePathDoesntExist, path)
}

func WorkspacePathIsExpectedToBeADirectoryError(path string) error {
	return fmt.Errorf("%w: %s", ErrWorkspacePathIsExpectedToBeADirectory, path)
}

func ComputationError(errList []error) error {
	return fmt.Errorf("%w:\n\t - %w", ErrDuringComputation, JoinErrors(errList, "\n\t - "))
}

func RepositoryNameIsMandatoryForConfigIndexError(index int) error {
	return fmt.Errorf("config #%d: %w", index, ErrRepositoryNameIsMandatory)
}

func RepositoryNameIsMandatoryForRepoError(index int) error {
	return fmt.Errorf("repo #%d: %w", index, ErrRepositoryNameIsMandatory)
}

func UnknownTemplateError(tplType string, tplName string) error {
	var baseError error

	switch tplType {
	case RepositoryTemplateType:
		baseError = ErrRepositoryTemplateNotFound
	case BranchTemplateType:
		baseError = ErrBranchTemplateNotFound
	case BranchProtectionTemplateType:
		baseError = ErrBranchProtectionTemplateNotFound
	default:
		return fmt.Errorf("\"%s\" %s template %w", tplName, tplType, ErrTemplateNotFound)
	}

	return fmt.Errorf("\"%s\" %w", tplName, baseError)
}

func NoTemplateAvailableError(tplType string) error {
	switch tplType {
	case RepositoryTemplateType:
		return ErrNoRepositoryTemplateAvailable
	case BranchTemplateType:
		return ErrNoBranchTemplateAvailable
	case BranchProtectionTemplateType:
		return ErrNoBranchProtectionTemplateAvailable
	default:
		return fmt.Errorf("%s template %w", tplType, ErrNoTemplateAvailable)
	}
}

func MaxTemplateCountReachedError(tplType string, path []string) error {
	pathString := "ROOT"

	if len(path) > 0 {
		pathString = strings.Join(path, "->")
	}

	return fmt.Errorf(
		"%w: more than %d %s template detected for %s",
		ErrMaxTemplateCount,
		TemplateMaxCount,
		tplType,
		pathString,
	)
}

func MaxTemplateDepthReachedError(tplType string, path []string) error {
	pathString := "ROOT"

	if len(path) > 0 {
		pathString = strings.Join(path, "->")
	}

	return fmt.Errorf(
		"%w: more than %d levels of %s template detected for %s",
		ErrMaxTemplateDepth,
		TemplateMaxDepth,
		tplType,
		pathString,
	)
}

func SchemaValidationError(path string, location string, msg string) error {
	return fmt.Errorf(
		"%w: file %s: %s %s",
		ErrSchemaValidation,
		path,
		location,
		msg,
	)
}

func EmptySchemaError(url string) error {
	return fmt.Errorf("%w: url %q", ErrEmptySchema, url)
}

func SchemaNotFoundError(url string) error {
	return fmt.Errorf("%w: url %q", ErrSchemaNotFound, url)
}

func SchemaIsNilError(url string) error {
	return fmt.Errorf("%w: url %q", ErrSchemaIsNil, url)
}

func SchemaCompilationError(url string, msg string) error {
	return fmt.Errorf("%w: url %s / error %s", ErrDuringSchemaCompilation, url, msg)
}

func TerraformFileWritingErrors(errList []error) error {
	return fmt.Errorf("%w:\n\t - %w", ErrDuringWriteTerraformFiles, JoinErrors(errList, "\n\t - "))
}

func SortErrorsByKey(errList map[string]error) []error {
	// sort file to always get a predictable output (for tests mostly)
	newErrorList := []error{}
	keys := make([]string, 0, len(errList))

	for k := range errList {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, file := range keys {
		newErrorList = append(newErrorList, errList[file])
	}

	return newErrorList
}

func JoinErrors(errList []error, separator string) error {
	if separator == "\n" {
		return errors.Join(errList...)
	}

	err := errList[0]
	for _, subErr := range errList[1:] {
		err = fmt.Errorf("%w%s%w", err, separator, subErr)
	}

	return err
}
