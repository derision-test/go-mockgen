package main

import (
	"fmt"
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

	wrappedType struct {
		typ types.Type
	}

	wrappedTuple struct {
		name string
		typ  *wrappedType
	}

	wrappedInterface struct {
		methods map[string]*wrappedMethod
	}

	wrappedMethod struct {
		params   []*wrappedTuple
		results  []*wrappedType
		variadic bool
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

//
// Types

func (t *wrappedType) Text() string {
	return t.typ.String()
}

func (t *wrappedType) Zero() {
	// TOOD
}

func sanitizeType(typ types.Type, pkg *types.Package, packageName string, packageNames []string) types.Type {
	switch t := typ.(type) {
	case *types.Basic:
		return typ

	case *types.Interface:
		return sanitizeInterface(t, pkg, packageName, packageNames)

	case *types.Named:
		if name := t.String(); stringInSlice(name, packageNames) {
			fullName := fmt.Sprintf("%s.%s", packageName, name)

			return types.NewNamed(
				types.NewTypeName(0, pkg, fullName, nil),
				nil,
				nil,
			)
		}

		return typ

	case *types.Pointer:
		return types.NewPointer(sanitizeType(t.Elem(), pkg, packageName, packageNames))

	case *types.Slice:
		return types.NewSlice(sanitizeType(t.Elem(), pkg, packageName, packageNames))

	case *types.Struct:
		fields := []*types.Var{}
		for i := 0; i < t.NumFields(); i++ {
			fields = append(fields, sanitizeField(
				t.Field(i),
				pkg,
				packageName,
				packageNames,
			))
		}

		return types.NewStruct(fields, nil)

	case *types.Signature:
		params := sanitizeParams(
			t.Params(),
			pkg,
			packageName,
			packageNames,
		)

		results := sanizeResults(
			t.Results(),
			pkg,
			packageName,
			packageNames,
		)

		return types.NewSignature(
			nil,
			types.NewTuple(params...),
			types.NewTuple(results...),
			t.Variadic(),
		)

	default:
		panic(fmt.Sprintf("unsupported case: %#v\n", typ))
	}
}

func sanitizeInterface(typ *types.Interface, pkg *types.Package, packageName string, packageNames []string) *types.Interface {
	methods := []*types.Func{}
	for i := 0; i < typ.NumMethods(); i++ {
		var (
			method    = typ.Method(i)
			sanitized = sanitizeType(method.Type(), pkg, packageName, packageNames)
		)

		methods = append(methods, types.NewFunc(
			0,
			pkg,
			method.Name(),
			sanitized.(*types.Signature),
		))
	}

	t := types.NewInterface(methods, nil)
	t.Complete()
	return t
}

func sanitizeField(typ *types.Var, pkg *types.Package, packageName string, packageNames []string) *types.Var {
	inner := sanitizeType(
		typ.Type(),
		pkg,
		packageName,
		packageNames,
	)

	return types.NewVar(
		0,
		pkg,
		typ.Name(),
		inner,
	)
}

func sanitizeParams(ps *types.Tuple, pkg *types.Package, packageName string, packageNames []string) []*types.Var {
	params := []*types.Var{}

	for i := 0; i < ps.Len(); i++ {
		name, typ := ith(ps, i)

		params = append(params, types.NewVar(
			0,
			pkg,
			name,
			sanitizeType(
				typ,
				pkg,
				packageName,
				packageNames,
			),
		))
	}

	return params
}

func sanizeResults(rs *types.Tuple, pkg *types.Package, packageName string, packageNames []string) []*types.Var {
	results := []*types.Var{}

	for i := 0; i < rs.Len(); i++ {
		_, typ := ith(rs, i)

		results = append(results, types.NewVar(
			0,
			pkg,
			"",
			sanitizeType(
				typ,
				pkg,
				packageName,
				packageNames,
			),
		))
	}

	return results
}

func deconstructMethod(signature *types.Signature) *wrappedMethod {
	var (
		ps      = signature.Params()
		rs      = signature.Results()
		params  = []*wrappedTuple{}
		results = []*wrappedType{}
	)

	for i := 0; i < ps.Len(); i++ {
		param := ps.At(i)

		params = append(params, &wrappedTuple{
			param.Name(),
			&wrappedType{param.Type()},
		})
	}

	for i := 0; i < rs.Len(); i++ {
		results = append(results, &wrappedType{rs.At(i).Type()})
	}

	return &wrappedMethod{
		params:   params,
		results:  results,
		variadic: signature.Variadic(),
	}
}

//
// Helpers

func ith(typ *types.Tuple, i int) (string, types.Type) {
	return typ.At(i).Name(), typ.At(i).Type()
}

func stringInSlice(needle string, haystack []string) bool {
	for _, elem := range haystack {
		if needle == elem {
			return true
		}
	}

	return false
}
