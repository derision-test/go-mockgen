package extraction

import (
	"go/ast"
	"go/types"

	"github.com/efritz/go-mockgen/specs"
)

type visitor struct {
	pkg   *types.Package
	specs specs.InterfaceSpecs
}

func newVisitor(pkg *types.Package) *visitor {
	return &visitor{
		pkg:   pkg,
		specs: specs.InterfaceSpecs{},
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				v.visitTypeSpec(typeSpec)
			}
		}
	}

	return v
}

func (v *visitor) visitTypeSpec(typeSpec *ast.TypeSpec) {
	typ := v.getInterfaceObject(typeSpec, v.pkg.Scope())
	if typ == nil {
		return
	}

	methods := specs.MethodSpecs{}
	for i := 0; i < typ.NumMethods(); i++ {
		methods[typ.Method(i).Name()] = deconstructMethod(typ.Method(i).Type().(*types.Signature))
	}

	v.specs[typeSpec.Name.Name] = &specs.InterfaceSpec{
		Methods: methods,
	}
}

func (v *visitor) getInterfaceObject(typeSpec *ast.TypeSpec, scope *types.Scope) *types.Interface {
	if !typeSpec.Name.IsExported() {
		return nil
	}

	_, obj := scope.Innermost(typeSpec.Pos()).LookupParent(typeSpec.Name.Name, 0)

	switch t := obj.Type().Underlying().(type) {
	case *types.Interface:
		return t
	default:
		return nil
	}
}

func deconstructMethod(signature *types.Signature) *specs.MethodSpec {
	var (
		ps      = signature.Params()
		rs      = signature.Results()
		params  = []types.Type{}
		results = []types.Type{}
	)

	for i := 0; i < ps.Len(); i++ {
		params = append(params, ps.At(i).Type())
	}

	for i := 0; i < rs.Len(); i++ {
		results = append(results, rs.At(i).Type())
	}

	return &specs.MethodSpec{
		Params:   params,
		Results:  results,
		Variadic: signature.Variadic(),
	}
}
