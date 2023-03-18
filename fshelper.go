package main

import (
	"errors"
	"fmt"
	"os"
)

var errPathIsNotADirectory = errors.New("path is not a directory")

func PathIsNotADirectoryError(path string) error {
	return fmt.Errorf("%w: %s", errPathIsNotADirectory, path)
}

func readDirectory(rootPath string) ([]string, error) {
	directory, dirErr := openDirectory(rootPath)
	if dirErr != nil {
		return nil, dirErr
	}

	//nolint:wrapcheck // Expecred to return error as is
	return directory.Readdirnames(0)
}

func openDirectory(path string) (*os.File, error) {
	directory, readErr := os.Open(path)
	if readErr != nil {
		//nolint:wrapcheck // Expecred to return error as is
		return nil, readErr
	}

	directoryStat, statErr := directory.Stat()
	if statErr != nil {
		//nolint:wrapcheck // Expecred to return error as is
		return nil, statErr
	}

	if !directoryStat.IsDir() {
		return nil, PathIsNotADirectoryError(path)
	}

	return directory, nil
}
