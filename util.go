package main

import (
	"fmt"
	"os"
	"sort"
)

func abort(err error) {
	fmt.Printf("error: %s", err.Error())
	os.Exit(1)
}

func getNames(specs map[string]*wrappedSpec) []string {
	names := []string{}
	for name := range specs {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}
