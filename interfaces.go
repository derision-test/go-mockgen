package main

import (
	"go/ast"
	"go/types"
)

type (
	interfaceExtractor struct {
		pkg          *types.Package
		packageName  string
		packageNames []string
		specs        map[string]*wrappedInterface
	}
)

func newInterfaceExtractor(pkg *types.Package, packageName string, packageNames []string) *interfaceExtractor {
	return &interfaceExtractor{
		pkg:          pkg,
		packageName:  packageName,
		packageNames: packageNames,
		specs:        map[string]*wrappedInterface{},
	}
}

func (e *interfaceExtractor) visitTypeSpec(typeSpec *ast.TypeSpec) {
	obj := e.getInterfaceObject(typeSpec, e.pkg.Scope())
	if obj == nil {
		return
	}

	typ := sanitizeInterface(
		obj,
		e.pkg,
		e.packageName,
		e.packageNames,
	)

	methods := map[string]*wrappedMethod{}
	for i := 0; i < typ.NumMethods(); i++ {
		method := typ.Method(i)
		signature := method.Type().(*types.Signature)

		methods[method.Name()] = deconstructMethod(signature)
	}

	e.specs[typeSpec.Name.Name] = &wrappedInterface{
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
