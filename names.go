package main

import (
	"go/ast"
)

type nameExtractor struct {
	names []string
}

func getNames(pkg *ast.Package) []string {
	e := newNameExtractor()
	walk(pkg, e)
	return e.names
}

func newNameExtractor() *nameExtractor {
	return &nameExtractor{
		names: []string{},
	}
}

func (e *nameExtractor) visitTypeSpec(spec *ast.TypeSpec) {
	name := spec.Name.Name

	if ast.IsExported(name) {
		e.names = append(e.names, name)
	}
}
