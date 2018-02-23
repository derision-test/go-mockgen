package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	importPaths = kingpin.Arg("path", "").Required().Strings()
	pkgName     = kingpin.Flag("package", "").Short('p').Required().String()
	interfaces  = kingpin.Flag("interfaces", "").Short('i').Strings()
)

func main() {
	kingpin.Parse()

	allSpecs := map[string]*wrappedSpec{}
	for _, path := range *importPaths {
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

	for _, name := range *interfaces {
		if _, ok := allSpecs[name]; !ok {
			abort(fmt.Errorf("interface %s not found in supplied import paths", name))
		}
	}

	if err := generate(allSpecs, *pkgName); err != nil {
		abort(err)
	}
}

func shouldInclude(name string) bool {
	for _, v := range *interfaces {
		if v == name {
			return true
		}
	}

	return len(*interfaces) == 0
}

func abort(err error) {
	fmt.Printf("error: %s", err.Error())
	os.Exit(1)
}
