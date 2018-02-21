package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"strings"
)

var (
	fset       = token.NewFileSet()
	typeConfig = types.Config{
		Importer: importer.Default(),
	}
)

func parseDir(name string) (*ast.Package, *types.Package, error) {
	pkgs, err := parser.ParseDir(fset, name, fileFilter, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("could not import package %s", name)
	}

	files := []*ast.File{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
	}

	pkgType, err := typeConfig.Check("", fset, files, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not import package %s", name)
	}

	if pkg := getFirst(pkgs); pkg != nil {
		return pkg, pkgType, nil
	}

	return nil, nil, fmt.Errorf("could not import package %s", name)
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
	if len(pkgs) == 1 {
		for _, pkg := range pkgs {
			return pkg
		}
	}

	return nil
}
