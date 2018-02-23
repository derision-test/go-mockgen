package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
)

var (
	importPath = kingpin.Arg("path", "").Required().String()
	interfaces = kingpin.Arg("interfaces", "").Strings()
)

func main() {
	kingpin.Parse()

	pkg, pkgType, err := parseImportPath(*importPath)
	if err != nil {
		abort(err)
	}

	names := getNames(pkg)
	specs := filter(getInterfaceSpecs(pkg, pkgType))

	if len(specs) == 0 {
		abort(fmt.Errorf("no interfaces found"))
	}

	if err := generate(specs, *importPath, names); err != nil {
		abort(err)
	}
}

func filter(specs map[string]*interfaceSpec) map[string]*interfaceSpec {
	if len(*interfaces) == 0 {
		return specs
	}

	filtered := map[string]*interfaceSpec{}
	for _, name := range *interfaces {
		spec, ok := specs[name]
		if !ok {
			abort(fmt.Errorf("interface %s not found in import path", name))
		}

		filtered[name] = spec
	}

	return filtered
}

func abort(err error) {
	fmt.Printf("error: %s", err.Error())
	os.Exit(1)
}
