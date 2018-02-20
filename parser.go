package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

var fset = token.NewFileSet()

func parseDir(name string) (*ast.Package, error) {
	pkgs, err := parser.ParseDir(fset, name, fileFilter, 0)
	if err != nil {
		return nil, err
	}

	pkg := getFirst(pkgs)
	if pkg == nil {
		return nil, fmt.Errorf("could not import package %s", name)
	}

	return pkg, nil
}

func fileFilter(info os.FileInfo) bool {
	var (
		name = info.Name()
		ext  = path.Ext(name)
		base = strings.TrimSuffix(name, ext)
	)

	return !info.IsDir() && ext == ".go" && !strings.HasSuffix(base, "_test")
}

func getFirst(pkgs map[string]*ast.Package) *ast.Package {
	for _, pkg := range pkgs {
		return pkg
	}

	return nil
}
