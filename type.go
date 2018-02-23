package main

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

func generateType(typ types.Type, importPath string, variadic bool) *jen.Statement {
	recur := func(typ types.Type) *jen.Statement {
		return generateType(typ, importPath, false)
	}

	switch t := typ.(type) {
	case *types.Basic:
		return jen.Id(typ.String())

	case *types.Chan:
		if t.Dir() == types.RecvOnly {
			return compose(jen.Op("<-").Chan(), recur(t.Elem()))
		}

		if t.Dir() == types.SendOnly {
			return compose(jen.Chan().Op("<-"), recur(t.Elem()))
		}

		return compose(jen.Chan(), recur(t.Elem()))

	case *types.Interface:
		methods := []jen.Code{}
		for i := 0; i < t.NumMethods(); i++ {
			methods = append(methods, compose(jen.Id(t.Method(i).Name()), recur(t.Method(i).Type())))
		}

		return jen.Interface(methods...)

	case *types.Map:
		return compose(jen.Map(recur(t.Key())), recur(t.Elem()))

	case *types.Named:
		return generateQualifiedName(t, importPath)

	case *types.Pointer:
		return compose(jen.Op("*"), recur(t.Elem()))

	case *types.Signature:
		params := []jen.Code{}
		for i := 0; i < t.Params().Len(); i++ {
			params = append(params, compose(jen.Id(t.Params().At(i).Name()), recur(t.Params().At(i).Type())))
		}

		results := []jen.Code{}
		for i := 0; i < t.Results().Len(); i++ {
			results = append(results, recur(t.Results().At(i).Type()))
		}

		return jen.Func().Params(params...).Params(results...)

	case *types.Slice:
		return compose(getSliceTypePrefix(variadic), recur(t.Elem()))

	case *types.Struct:
		fields := []jen.Code{}
		for i := 0; i < t.NumFields(); i++ {
			fields = append(fields, compose(jen.Id(t.Field(i).Name()), recur(t.Field(i).Type())))
		}

		return jen.Struct(fields...)

	default:
		panic(fmt.Sprintf("unsupported case: %#v\n", typ))
	}
}

func zeroValue(typ types.Type, importPath string) *jen.Statement {
	switch t := typ.(type) {
	case *types.Basic:
		kind := t.Kind()

		if kind == types.Bool {
			return jen.False()
		} else if kind == types.String {
			return jen.Lit("")
		} else if isIntegerType(kind) {
			return jen.Lit(0)
		}

	case *types.Named:
		if _, ok := t.Underlying().(*types.Struct); ok {
			return compose(generateQualifiedName(t, importPath), jen.Block())
		}

		return zeroValue(t.Underlying(), importPath)

	case *types.Struct:
		return generateType(typ, importPath, false).Block()
	}

	return jen.Nil()
}

//
// Helpers

func getSliceTypePrefix(variadic bool) *jen.Statement {
	if variadic {
		return jen.Op("...")
	}

	return jen.Index()
}

func isIntegerType(kind types.BasicKind) bool {
	kinds := []types.BasicKind{
		types.Int,
		types.Int8,
		types.Int16,
		types.Int32,
		types.Int64,
		types.Uint,
		types.Uint8,
		types.Uint16,
		types.Uint32,
		types.Uint64,
		types.Uintptr,
		types.Float32,
		types.Float64,
		types.Byte,
		types.Rune,
		types.Complex64,
		types.Complex128,
	}

	for _, k := range kinds {
		if k == kind {
			return true
		}
	}

	return false
}

func generateQualifiedName(t *types.Named, importPath string) *jen.Statement {
	name := t.Obj().Name()

	if t.Obj().Pkg() == nil {
		return jen.Id(name)
	}

	if path := t.Obj().Pkg().Path(); path != "" {
		return jen.Qual(path, name)
	}

	return jen.Qual(importPath, name)
}
