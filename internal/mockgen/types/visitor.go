package types

import (
	"go/ast"
	"go/types"
)

type visitor struct {
	importPath string
	pkgType    *types.Package
	types      map[string]*Interface
}

func newVisitor(importPath string, pkgType *types.Package) *visitor {
	return &visitor{
		importPath: importPath,
		pkgType:    pkgType,
		types:      map[string]*Interface{},
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.File:
		return v

	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				name := typeSpec.Name.Name
				_, obj := v.pkgType.Scope().Innermost(typeSpec.Pos()).LookupParent(name, 0)

				switch t := obj.Type().Underlying().(type) {
				case *types.Interface:
					v.types[name] = newInterfaceFromTypeSpec(name, v.importPath, t)
				}
			}
		}
	}

	return nil
}
