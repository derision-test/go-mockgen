package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/efritz/go-mockgen/extraction"
	"github.com/efritz/go-mockgen/generation"
)

const Version = "0.1.0"

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() error {
	wd, dirname, filename, err := parseArgs()
	if err != nil {
		return err
	}

	allSpecs, err := extraction.Extract(wd, *importPaths, *interfaces)
	if err != nil {
		return err
	}

	if *listOnly {
		for _, name := range allSpecs.Names() {
			fmt.Printf("%s\n", name)
		}

		return nil
	}

	for _, name := range *interfaces {
		if _, ok := allSpecs[strings.ToLower(name)]; !ok {
			return fmt.Errorf("interface %s not found in supplied import paths", name)
		}
	}

	return generation.Generate(allSpecs, *pkgName, *prefix, dirname, filename, *force)
}
