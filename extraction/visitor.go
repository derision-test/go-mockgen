package extraction

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/efritz/go-mockgen/specs"
)

type visitor struct {
	pkg   *types.Package
	specs specs.InterfaceSpecs
	err   error
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
				if err := v.visitTypeSpec(typeSpec); err != nil && v.err == nil {
					v.err = err
				}
			}
		}
	}

	return v
}

func (v *visitor) visitTypeSpec(typeSpec *ast.TypeSpec) error {
	typ := v.getInterfaceObject(typeSpec, v.pkg.Scope())
	if typ == nil {
		return nil
	}

	methods := specs.MethodSpecs{}
	for i := 0; i < typ.NumMethods(); i++ {
		method := typ.Method(i)

		if !method.Exported() {
			return fmt.Errorf(
				"interface %s contains unexported method %s",
				typeSpec.Name,
				typ.Method(i).Name())
		}

		methods[method.Name()] = deconstructMethod(method.Type().(*types.Signature))
	}

	var (
		name      = typeSpec.Name.Name
		titleName = title(name)
		lowerName = strings.ToLower(name)
	)

	spec := &specs.InterfaceSpec{
		Name:      name,
		TitleName: titleName,
		Methods:   methods,
	}

	v.specs[lowerName] = spec
	return nil
}

func (v *visitor) getInterfaceObject(typeSpec *ast.TypeSpec, scope *types.Scope) *types.Interface {
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

func title(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}
