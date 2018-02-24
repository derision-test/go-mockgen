package main

import "fmt"

func main() {
	dirname, filename, err := parseArgs()
	if err != nil {
		abort(err)
	}

	allSpecs, err := getSpecs(*importPaths, *interfaces)
	if err != nil {
		abort(err)
	}

	for _, name := range *interfaces {
		if _, ok := allSpecs[name]; !ok {
			abort(fmt.Errorf("interface %s not found in supplied import paths", name))
		}
	}

	if *listOnly {
		for _, name := range getNames(allSpecs) {
			fmt.Printf("%s\n", name)
		}

		return
	}

	if err := generate(allSpecs, *pkgName, *prefix, dirname, filename, *force); err != nil {
		abort(err)
	}
}
