package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
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

func stripVendor(path string) string {
	parts := strings.Split(path, "/vendor/")
	return parts[len(parts)-1]
}
