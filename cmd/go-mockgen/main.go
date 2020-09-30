package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/derision-test/go-mockgen/internal/mockgen/generation"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := mainErr(); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func mainErr() error {
	opts, err := parseArgs()
	if err != nil {
		return err
	}

	ifaces, err := types.Extract(opts.ImportPaths, opts.Interfaces)
	if err != nil {
		return err
	}

	nameMap := map[string]struct{}{}
	for _, t := range ifaces {
		nameMap[strings.ToLower(t.Name)] = struct{}{}
	}

	for _, name := range opts.Interfaces {
		if _, ok := nameMap[strings.ToLower(name)]; !ok {
			return fmt.Errorf("type '%s' not found in supplied import paths", name)
		}
	}

	return generation.Generate(ifaces, opts)
}
