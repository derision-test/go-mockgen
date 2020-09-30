package generation

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
)

func SanitizeImportPath(path, outputImportPath string) string {
	path = stripVendor(path)
	if path == outputImportPath {
		return ""
	}

	return path
}

func GenerateType(typ types.Type, importPath, outputImportPath string, variadic bool) *jen.Statement {
	recur := func(typ types.Type) *jen.Statement {
		return GenerateType(typ, importPath, outputImportPath, false)
	}

	switch t := typ.(type) {
	case *types.Basic:
		return jen.Id(typ.String())

	case *types.Chan:
		if t.Dir() == types.RecvOnly {
			return Compose(jen.Op("<-").Chan(), recur(t.Elem()))
		}

		if t.Dir() == types.SendOnly {
			return Compose(jen.Chan().Op("<-"), recur(t.Elem()))
		}

		return Compose(jen.Chan(), recur(t.Elem()))

	case *types.Interface:
		methods := []jen.Code{}
		for i := 0; i < t.NumMethods(); i++ {
			methods = append(methods, Compose(jen.Id(t.Method(i).Name()), recur(t.Method(i).Type())))
		}

		return jen.Interface(methods...)

	case *types.Map:
		return Compose(jen.Map(recur(t.Key())), recur(t.Elem()))

	case *types.Named:
		return generateQualifiedName(t, importPath, outputImportPath)

	case *types.Pointer:
		return Compose(jen.Op("*"), recur(t.Elem()))

	case *types.Signature:
		params := []jen.Code{}
		for i := 0; i < t.Params().Len(); i++ {
			params = append(params, Compose(jen.Id(t.Params().At(i).Name()), recur(t.Params().At(i).Type())))
		}

		results := []jen.Code{}
		for i := 0; i < t.Results().Len(); i++ {
			results = append(results, recur(t.Results().At(i).Type()))
		}

		return jen.Func().Params(params...).Params(results...)

	case *types.Slice:
		return Compose(getSliceTypePrefix(variadic), recur(t.Elem()))

	case *types.Struct:
		fields := []jen.Code{}
		for i := 0; i < t.NumFields(); i++ {
			fields = append(fields, Compose(jen.Id(t.Field(i).Name()), recur(t.Field(i).Type())))
		}

		return jen.Struct(fields...)

	default:
		panic(fmt.Sprintf("unsupported case: %#v\n", typ))
	}
}

func stripVendor(path string) string {
	parts := strings.Split(path, "/vendor/")
	return parts[len(parts)-1]
}

func generateQualifiedName(t *types.Named, importPath, outputImportPath string) *jen.Statement {
	name := t.Obj().Name()

	if t.Obj().Pkg() == nil {
		return jen.Id(name)
	}

	if path := t.Obj().Pkg().Path(); path != "" {
		return jen.Qual(SanitizeImportPath(path, outputImportPath), name)
	}

	return jen.Qual(SanitizeImportPath(importPath, outputImportPath), name)
}

func getSliceTypePrefix(variadic bool) *jen.Statement {
	if variadic {
		return jen.Op("...")
	}

	return jen.Index()
}
