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
	fset       *token.FileSet
	typeConfig types.Config
}

func newPathParser() *pathParser {
	return &pathParser{
		fset:       token.NewFileSet(),
		typeConfig: types.Config{Importer: importer.For("source", nil)},
	}
}

func (p *pathParser) parse(path string) (*ast.Package, *types.Package, error) {
	pkgs, err := parser.ParseDir(p.fset, filepath.Join(paths.Gopath(), "src", path), fileFilter, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("could not import package %s (%s)", path, err.Error())
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
