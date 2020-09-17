package extraction

import (
	"go/ast"
	"go/token"
	gotypes "go/types"

	"github.com/efritz/go-mockgen/internal/genlib/types"
)

type visitor struct {
	importPath string
	pkgType    *gotypes.Package
	types      map[string]*types.Interface
}

func newVisitor(importPath string, pkgType *gotypes.Package) *visitor {
	return &visitor{
		importPath: importPath,
		pkgType:    pkgType,
		types:      map[string]*types.Interface{},
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				v.deconstructTypeSpec(typeSpec)
			}
		}
	}

	return v
}

func (v *visitor) deconstructTypeSpec(typeSpec *ast.TypeSpec) {
	name := typeSpec.Name.Name

	switch t := getUnderlyingType(v.pkgType, name, typeSpec.Pos()).(type) {
	case *gotypes.Struct:
		v.types[name] = types.DeconstructStruct(name, v.importPath, t)
	case *gotypes.Interface:
		v.types[name] = types.DeconstructInterface(name, v.importPath, t)
	}
}

func getUnderlyingType(pkgType *gotypes.Package, name string, pos token.Pos) gotypes.Type {
	_, obj := pkgType.Scope().Innermost(pos).LookupParent(name, 0)
	return obj.Type().Underlying()
}
