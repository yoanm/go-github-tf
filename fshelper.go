package main

import (
	"fmt"
	"os"
)

func readDirectory(rootPath string) ([]string, error) {
	directory, dirErr := openDirectory(rootPath)
	if dirErr != nil {
		return nil, dirErr
	}

	return directory.Readdirnames(0)
}

func openDirectory(path string) (*os.File, error) {
	directory, readErr := os.Open(path)
	if readErr != nil {
		return nil, readErr
	}
	directoryStat, statErr := directory.Stat()
	if statErr != nil {
		return nil, statErr
	}
	if !directoryStat.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	return directory, nil
}
