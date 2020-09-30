package extraction

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/token"
	gotypes "go/types"
	"log"
	"os"

	"github.com/derision-test/go-mockgen/internal/genlib/paths"
	"github.com/derision-test/go-mockgen/internal/genlib/types"
	gopackages "golang.org/x/tools/go/packages"
)

type Extractor struct {
	workingDirectory string
	fset             *token.FileSet
	typeConfig       gotypes.Config
}

func NewExtractor() (*Extractor, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory (%s)", err.Error())
	}

	return &Extractor{
		workingDirectory: workingDirectory,
		fset:             token.NewFileSet(),
		typeConfig:       gotypes.Config{Importer: importer.For("source", nil)},
	}, nil
}

func (e *Extractor) Extract(importPaths []string) (*types.Packages, error) {
	packageConfig := &gopackages.Config{
		Mode: gopackages.LoadSyntax,
	}

	packages := map[string]*types.Package{}
	for _, importPath := range importPaths {
		path, dir := paths.ResolveImportPath(e.workingDirectory, importPath)

		log.Printf(
			"parsing package '%s'\n",
			paths.GetRelativePath(dir),
		)

		pkgs, err := gopackages.Load(packageConfig, importPath)
		if err != nil {
			return nil, fmt.Errorf(
				"could not load package %s (%s)",
				importPath,
				err.Error(),
			)
		}

		for _, err := range pkgs[0].Errors {
			return nil, fmt.Errorf(
				"malformed package %s (%s)",
				importPath,
				err.Msg,
			)
		}

		visitor := newVisitor(path, pkgs[0].Types)
		for _, file := range pkgs[0].Syntax {
			ast.Walk(visitor, file)
		}

		packages[path] = types.NewPackage(path, visitor.types)
	}

	return types.NewPackages(packages), nil
}
