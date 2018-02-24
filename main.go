package main

import "fmt"

func main() {
	dirname, filename, err := parseArgs()
	if err != nil {
		abort(err)
	}

	allSpecs, err := getSpecs(*ImportPaths, *Interfaces)
	if err != nil {
		abort(err)
	}

	for _, name := range *Interfaces {
		if _, ok := allSpecs[name]; !ok {
			abort(fmt.Errorf("interface %s not found in supplied import paths", name))
		}
	}

	if *ListOnly {
		for _, name := range getNames(allSpecs) {
			fmt.Printf("%s\n", name)
		}

		return
	}

	if err := generate(allSpecs, *PkgName, dirname, filename, *Force); err != nil {
		abort(err)
	}
}
