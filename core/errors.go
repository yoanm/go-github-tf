package core

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errDuringWriteTerraformFiles     = errors.New("error while writing terraform files")
	errRepositoryNameIsMandatory     = errors.New("repository name is mandatory")
	errDuringFileGeneration          = errors.New("error while generating files")
	errFileOpenNoSuchFileOrDirectory = errors.New("no such file or directory")
	errPathIsNotADirectory           = errors.New("path is not a directory")
	errTemplateUnavailable           = errors.New("template unavailable")
	errDuringComputation             = errors.New("error during computation")
	errRepositoryNameIsMissing       = errors.New("repository name is missing")
	errUnknownTemplate               = errors.New("unknown template")
	errMaxTemplateCount              = errors.New("maximum template count reached")
	errMaxTemplateDepth              = errors.New("maximum template depth reached")

	// Json Schema and validation.
	errSchemaValidation        = errors.New("schema validation error")
	errEmptySchema             = errors.New("empty schema")
	errSchemaNotFound          = errors.New("schema not found")
	errSchemaIsNil             = errors.New("schema is nil")
	errDuringSchemaCompilation = errors.New("error during schema compilation")
)

func WriteTerraformFileError(errList []error) error {
	msgList := []string{}
	for _, err := range errList {
		msgList = append(msgList, err.Error())
	}

	return fmt.Errorf("%w\n\t - %s", errDuringWriteTerraformFiles, strings.Join(msgList, "\n\t - "))
}

func FileGenerationError(msgList []string) error {
	return fmt.Errorf("%w:\n\t - %s", errDuringFileGeneration, strings.Join(msgList, "\n\t - "))
}

func FileOpenNoSuchFileOrDirectoryError(path string) error {
	return fmt.Errorf("%w: %s", errFileOpenNoSuchFileOrDirectory, path)
}

func PathIsNotADirectoryError(path string) error {
	return fmt.Errorf("%w: %s", errPathIsNotADirectory, path)
}

func TemplateUnavailableError(tplType string) error {
	return fmt.Errorf("%w: unable to load %s template", errTemplateUnavailable, tplType)
}

func ComputationError(msgList []string) error {
	return fmt.Errorf("%w:\n\t - %s", errDuringComputation, strings.Join(msgList, "\n\t - "))
}

func RepositoryNameIsMandatoryForConfigIndexError(index int) error {
	return fmt.Errorf("%w: config #%d", errRepositoryNameIsMandatory, index)
}

func RepositoryNameIsMissingForRepoError(index int) error {
	return fmt.Errorf("%w: repo #%d", errRepositoryNameIsMissing, index)
}

func UnknownTemplateError(tplType string, tplName string) error {
	return fmt.Errorf("%w: %s template %s", errUnknownTemplate, tplType, tplName)
}

func MaxTemplateCountReachedError(tplType string, path []string) error {
	pathString := "ROOT"

	if len(path) > 0 {
		pathString = strings.Join(path, "->")
	}

	return fmt.Errorf(
		"%w: more than %d %s template detected for %s",
		errMaxTemplateCount,
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
		errMaxTemplateDepth,
		TemplateMaxDepth,
		tplType,
		pathString,
	)
}

func SchemaValidationError(path string, location string, msg string) error {
	return fmt.Errorf(
		"%w: file %s: %s %s",
		errSchemaValidation,
		path,
		location,
		msg,
	)
}

func EmptySchemaError(url string) error {
	return fmt.Errorf("%w: url %q", errEmptySchema, url)
}

func SchemaNotFoundError(url string) error {
	return fmt.Errorf("%w: url %q", errSchemaNotFound, url)
}

func SchemaIsNilError(url string) error {
	return fmt.Errorf("%w: url %q", errSchemaIsNil, url)
}

func SchemaCompilationError(url string, msg string) error {
	return fmt.Errorf("%w: url %s / error %s", errDuringSchemaCompilation, url, msg)
}
