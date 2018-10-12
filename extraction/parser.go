package extraction

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/efritz/go-mockgen/paths"
)

type pathParser struct {
	wd         string
	fset       *token.FileSet
	typeConfig types.Config
}

func newPathParser(wd string) *pathParser {
	return &pathParser{
		wd:         wd,
		fset:       token.NewFileSet(),
		typeConfig: types.Config{Importer: importer.For("source", nil)},
	}
}

func (p *pathParser) parse(path string) (*ast.Package, *types.Package, error) {
	pkgs, err := p.importPackage(path)
	if err != nil {
		return nil, nil, err
	}

	files := []*ast.File{}
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
	}

	pkgType, err := p.typeConfig.Check("", p.fset, files, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not import package %s (%s)", path, err.Error())
	}

	if pkg := getFirst(pkgs); pkg != nil {
		return pkg, pkgType, nil
	}

	return nil, nil, fmt.Errorf("could not import package %s (no results)", path)
}

func (p *pathParser) importPackage(path string) (map[string]*ast.Package, error) {
	for _, possiblePath := range getPossiblePaths(p.wd, path) {
		pkgs, err := parser.ParseDir(p.fset, possiblePath, fileFilter, 0)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				return nil, fmt.Errorf(
					"could not import package from %s (%s)",
					possiblePath,
					err.Error(),
				)
			}

			continue
		}

		return pkgs, nil
	}

	return nil, fmt.Errorf("could not locate package %s", path)
}

func getPossiblePaths(wd, path string) []string {
	var (
		root       = filepath.Join(paths.Gopath(), "src")
		globalPath = filepath.Join(root, path)
	)

	if !strings.HasPrefix(wd, root) {
		return []string{globalPath}
	}

	paths := []string{}
	for wd != root {
		paths = append(paths, filepath.Join(wd, "vendor", path))
		wd = filepath.Dir(wd)
	}

	return append(paths, globalPath)
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
