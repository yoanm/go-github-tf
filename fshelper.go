package main

import (
	"os"
)

func readDirectory(rootPath string) ([]string, error) {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		//nolint:wrapcheck // Expected to return error as is
		return nil, err
	}

	res := make([]string, 0, len(files))

	for _, file := range files {
		res = append(res, file.Name())
	}

	return res, nil
}
