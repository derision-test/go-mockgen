package main

import (
	"fmt"
	"path"

	"github.com/alecthomas/kingpin"
)

var (
	ImportPaths    = kingpin.Arg("path", "").Required().Strings()
	PkgName        = kingpin.Flag("package", "").Short('p').String()
	Prefix         = kingpin.Flag("prefix", "").String()
	Interfaces     = kingpin.Flag("interfaces", "").Short('i').Strings()
	OutputDir      = kingpin.Flag("dirname", "").Short('d').String()
	OutputFilename = kingpin.Flag("filename", "").Short('o').String()
	Force          = kingpin.Flag("force", "").Short('f').Bool()
	ListOnly       = kingpin.Flag("list", "").Bool()
)

func parseArgs() (string, string, error) {
	kingpin.Parse()

	if *PkgName == "" && !*ListOnly {
		kingpin.Fatalf("required flag --package not provided, try --help")
	}

	return validateOutputPath(
		*OutputDir,
		*OutputFilename,
	)
}

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
