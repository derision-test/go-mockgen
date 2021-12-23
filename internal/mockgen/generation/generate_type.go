package generation

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

type typeGenerator func(typ types.Type) *jen.Statement

func generateType(typ types.Type, importPath, outputImportPath string, variadic bool) *jen.Statement {
	recur := func(typ types.Type) *jen.Statement {
		return generateType(typ, importPath, outputImportPath, false)
	}

	switch t := typ.(type) {
	case *types.Basic:
		return generateBasicType(t, recur)
	case *types.Chan:
		return generateChanType(t, recur)
	case *types.Interface:
		return generateInterfaceType(t, recur)
	case *types.Map:
		return generateMapType(t, recur)
	case *types.Named:
		return generateNamedType(t, importPath, outputImportPath, recur)
	case *types.Pointer:
		return generatePointerType(t, recur)
	case *types.Signature:
		return generateSignatureType(t, recur)
	case *types.Slice:
		return generateSliceType(t, variadic, recur)
	case *types.Struct:
		return generateStructType(t, recur)

	default:
		panic(fmt.Sprintf("unsupported case: %#v\n", typ))
	}
}

func generateBasicType(t *types.Basic, _ typeGenerator) *jen.Statement {
	return jen.Id(t.String())
}

func generateChanType(t *types.Chan, generate typeGenerator) *jen.Statement {
	c := jen.Chan()

	if t.Dir() == types.RecvOnly {
		c = compose(jen.Op("<-"), c)
	} else if t.Dir() == types.SendOnly {
		c = compose(c, jen.Op("<-"))
	}

	return compose(c, generate(t.Elem()))
}

func generateInterfaceType(t *types.Interface, generate typeGenerator) *jen.Statement {
	methods := make([]jen.Code, 0, t.NumMethods())
	for i := 0; i < t.NumMethods(); i++ {
		methods = append(methods, compose(jen.Id(t.Method(i).Name()), generate(t.Method(i).Type())))
	}

	return jen.Interface(methods...)
}

func generateMapType(t *types.Map, generate typeGenerator) *jen.Statement {
	return compose(jen.Map(generate(t.Key())), generate(t.Elem()))
}

func generateNamedType(t *types.Named, importPath, outputImportPath string, _ typeGenerator) *jen.Statement {
	return generateQualifiedName(t, importPath, outputImportPath)
}

func generatePointerType(t *types.Pointer, generate typeGenerator) *jen.Statement {
	return compose(jen.Op("*"), generate(t.Elem()))
}

func generateSignatureType(t *types.Signature, generate typeGenerator) *jen.Statement {
	params := make([]jen.Code, 0, t.Params().Len())
	for i := 0; i < t.Params().Len(); i++ {
		params = append(params, compose(jen.Id(t.Params().At(i).Name()), generate(t.Params().At(i).Type())))
	}

	results := make([]jen.Code, 0, t.Results().Len())
	for i := 0; i < t.Results().Len(); i++ {
		results = append(results, generate(t.Results().At(i).Type()))
	}

	return jen.Func().Params(params...).Params(results...)
}

func generateSliceType(t *types.Slice, variadic bool, generate typeGenerator) *jen.Statement {
	if variadic {
		return compose(jen.Op("..."), generate(t.Elem()))
	}

	return compose(jen.Index(), generate(t.Elem()))
}

func generateStructType(t *types.Struct, generate typeGenerator) *jen.Statement {
	fields := make([]jen.Code, 0, t.NumFields())
	for i := 0; i < t.NumFields(); i++ {
		fields = append(fields, compose(jen.Id(t.Field(i).Name()), generate(t.Field(i).Type())))
	}

	return jen.Struct(fields...)
}
