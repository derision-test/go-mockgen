package main

import (
	"go/ast"
	"go/types"
)

type (
	interfaceExtractor struct {
		pkg   *types.Package
		specs map[string]*interfaceSpec
	}

	interfaceSpec struct {
		methods map[string]*methodSpec
	}

	methodSpec struct {
		params   []types.Type
		results  []types.Type
		variadic bool
	}
)

func getInterfaceSpecs(pkg *ast.Package, pkgType *types.Package) map[string]*interfaceSpec {
	e := newInterfaceExtractor(pkgType)
	walk(pkg, e)
	return e.specs
}

func newInterfaceExtractor(pkg *types.Package) *interfaceExtractor {
	return &interfaceExtractor{
		pkg:   pkg,
		specs: map[string]*interfaceSpec{},
	}
}

func (e *interfaceExtractor) visitTypeSpec(typeSpec *ast.TypeSpec) {
	typ := e.getInterfaceObject(typeSpec, e.pkg.Scope())
	if typ == nil {
		return
	}

	methods := map[string]*methodSpec{}
	for i := 0; i < typ.NumMethods(); i++ {
		methods[typ.Method(i).Name()] = deconstructMethod(typ.Method(i).Type().(*types.Signature))
	}

	e.specs[typeSpec.Name.Name] = &interfaceSpec{
		methods: methods,
	}
}

func (e *interfaceExtractor) getInterfaceObject(typeSpec *ast.TypeSpec, scope *types.Scope) *types.Interface {
	_, obj := scope.Innermost(typeSpec.Pos()).LookupParent(typeSpec.Name.Name, 0)

	switch t := obj.Type().Underlying().(type) {
	case *types.Interface:
		return t
	default:
		return nil
	}
}

func deconstructMethod(signature *types.Signature) *methodSpec {
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

	return &methodSpec{
		params:   params,
		results:  results,
		variadic: signature.Variadic(),
	}
}
