package main

import (
	"fmt"
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

	wrappedSpec struct {
		importPath string
		spec       *interfaceSpec
	}
)

func getSpecs(importPaths []string, interfaces []string) (map[string]*wrappedSpec, error) {
	allSpecs := map[string]*wrappedSpec{}
	for _, path := range importPaths {
		pkg, pkgType, err := parseImportPath(path)
		if err != nil {
			abort(err)
		}

		specs := getInterfaceSpecs(pkg, pkgType)
		if len(specs) == 0 {
			return nil, fmt.Errorf("no interfaces found in path %s", path)
		}

		for name, spec := range specs {
			if shouldInclude(name, interfaces) {
				if _, ok := allSpecs[name]; ok {
					return nil, fmt.Errorf("ambiguous interface %s in supplied import paths", name)
				}

				allSpecs[name] = &wrappedSpec{spec: spec, importPath: path}
			}
		}
	}

	return allSpecs, nil
}

func shouldInclude(name string, interfaces []string) bool {
	for _, v := range interfaces {
		if v == name {
			return true
		}
	}

	return len(interfaces) == 0
}

func getInterfaceSpecs(pkg *ast.Package, pkgType *types.Package) map[string]*interfaceSpec {
	visitor := newInterfaceExtractor(pkgType)
	for _, file := range pkg.Files {
		ast.Walk(visitor, file)
	}

	return visitor.specs
}

func newInterfaceExtractor(pkg *types.Package) *interfaceExtractor {
	return &interfaceExtractor{
		pkg:   pkg,
		specs: map[string]*interfaceSpec{},
	}
}

func (e *interfaceExtractor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				e.visitTypeSpec(typeSpec)
			}
		}
	}

	return e
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
