package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/derision-test/go-mockgen/internal/mockgen/generation"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
	"golang.org/x/tools/go/packages"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := mainErr(); err != nil {
		message := fmt.Sprintf("error: %s\n", err.Error())

		if solvableError, ok := err.(solvableError); ok {
			message += "\nPossible solutions:\n"

			for _, hint := range solvableError.Solutions() {
				message += fmt.Sprintf("  - %s\n", hint)
			}

			message += "\n"
		}

		log.Fatalf(message)
	}
}

type solvableError interface {
	Solutions() []string
}

func mainErr() error {
	allOptions, err := parseAndValidateOptions()
	if err != nil {
		return err
	}

	var importPaths []string
	for _, opts := range allOptions {
		importPaths = append(importPaths, opts.ImportPaths...)
	}

	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName | packages.NeedImports | packages.NeedSyntax | packages.NeedTypes}, importPaths...)
	if err != nil {
		return fmt.Errorf("could not load packages %s (%s)", strings.Join(importPaths, ","), err.Error())
	}

	for _, opts := range allOptions {
		ifaces, err := types.Extract(pkgs, opts.ImportPaths, opts.Interfaces, opts.Exclude)
		if err != nil {
			return err
		}

		nameMap := make(map[string]struct{}, len(ifaces))
		for _, t := range ifaces {
			nameMap[strings.ToLower(t.Name)] = struct{}{}
		}

		for _, name := range opts.Interfaces {
			if _, ok := nameMap[strings.ToLower(name)]; !ok {
				return fmt.Errorf("type '%s' not found in supplied import paths", name)
			}
		}

		if err := generation.Generate(ifaces, opts); err != nil {
			return err
		}
	}

	return nil
}
