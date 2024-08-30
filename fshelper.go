package main

import (
	"os"
)

func readDirectory(rootPath string) ([]string, error) {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	var res []string

	for _, file := range files {
		res = append(res, file.Name())
	}

	return res, nil
}

func openDirectory(path string) (*os.File, error) {
	directory, readErr := os.Open(path)
	if readErr != nil {
		//nolint:wrapcheck // Expected to return error as is
		return nil, readErr
	}

	directoryStat, statErr := directory.Stat()
	if statErr != nil {
		//nolint:wrapcheck // Expected to return error as is
		return nil, statErr
	}

	if !directoryStat.IsDir() {
		return nil, pathIsNotADirectoryError(path)
	}

	return directory, nil
}
