package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	// TODO - add descriptions
	ImportPaths    = kingpin.Arg("path", "").Required().Strings()
	PkgName        = kingpin.Flag("package", "").Short('p').Required().String()
	Interfaces     = kingpin.Flag("interfaces", "").Short('i').Strings()
	OutputDir      = kingpin.Flag("dirname", "").Short('d').String()
	OutputFilename = kingpin.Flag("filename", "").Short('o').String()
	Force          = kingpin.Flag("force", "").Short('f').Bool()
)

func main() {
	kingpin.Parse()

	dirname, filename, err := validateOutputPath(
		*OutputDir,
		*OutputFilename,
	)

	if err != nil {
		abort(err)
	}

	allSpecs := map[string]*wrappedSpec{}
	for _, path := range *ImportPaths {
		pkg, pkgType, err := parseImportPath(path)
		if err != nil {
			abort(err)
		}

		specs := getInterfaceSpecs(pkg, pkgType)
		if len(specs) == 0 {
			abort(fmt.Errorf("no interfaces found in path %s", path))
		}

		for name, spec := range specs {
			if shouldInclude(name) {
				if _, ok := allSpecs[name]; ok {
					abort(fmt.Errorf("ambiguous interface %s in supplied import paths", name))
				}

				allSpecs[name] = &wrappedSpec{spec: spec, importPath: path}
			}
		}
	}

	for _, name := range *Interfaces {
		if _, ok := allSpecs[name]; !ok {
			abort(fmt.Errorf("interface %s not found in supplied import paths", name))
		}
	}

	if err := generate(allSpecs, *PkgName, dirname, filename, *Force); err != nil {
		abort(err)
	}
}

func shouldInclude(name string) bool {
	for _, v := range *Interfaces {
		if v == name {
			return true
		}
	}

	return len(*Interfaces) == 0
}

func abort(err error) {
	fmt.Printf("error: %s", err.Error())
	os.Exit(1)
}
