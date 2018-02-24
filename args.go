package main

import (
	"fmt"
	"path"
	"regexp"

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

var identPattern = regexp.MustCompile("^[A-Za-z]([A-Za-z0-9_]*[A-Za-z])?$")

func parseArgs() (string, string, error) {
	kingpin.Parse()

	dirname, filename, err := validateOutputPath(*OutputDir, *OutputFilename)
	if err != nil {
		return "", "", err
	}

	if *PkgName == "" && !*ListOnly {
		if dirname == "" {
			kingpin.Fatalf("could not infer package, try --help")
		}

		*PkgName = path.Base(dirname)
	}

	if !identPattern.Match([]byte(*PkgName)) {
		kingpin.Fatalf("illegal package name supplied, try --help")
	}

	if *Prefix != "" && !identPattern.Match([]byte(*Prefix)) {
		kingpin.Fatalf("illegal prefix supplied, try --help")
	}

	return dirname, filename, nil
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
