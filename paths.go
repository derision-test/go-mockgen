package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/alecthomas/kingpin"
)

func validateOutputPath(dirname, filename string) (string, string, error) {
	if dirname == "" && filename == "" {
		return "", "", nil
	}

	if filename != "" {
		if dirname != "" {
			kingpin.Fatalf("dirname and filename are mutually exclusive, try --help")
		}

		dirname, filename = path.Dir(filename), path.Base(filename)
	}

	exists, err := pathExists(dirname)
	if err != nil {
		return "", "", err
	}

	if !exists {
		return "", "", fmt.Errorf("directory %s does not exist", dirname)
	}

	return dirname, filename, nil
}

func getFilename(dirname, interfaceName string) string {
	return path.Join(dirname, fmt.Sprintf("%s_mock.go", strings.ToLower(interfaceName)))
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func anyPathExists(paths []string) (string, error) {
	for _, path := range paths {
		exists, err := pathExists(path)
		if err != nil {
			return "", err
		}

		if exists {
			return path, nil
		}
	}

	return "", nil
}
