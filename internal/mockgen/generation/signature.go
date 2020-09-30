package generation

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
)

func GenerateFunction(methodName string, params, results []jen.Code, body ...jen.Code) jen.Code {
	return jen.Func().
		Id(methodName).
		Params(params...).
		Params(results...).
		Block(body...)
}

func GenerateMethod(receiver jen.Code, methodName string, params, results []jen.Code, body ...jen.Code) jen.Code {
	return jen.Func().
		Params(receiver).
		Id(methodName).
		Params(params...).
		Params(results...).
		Block(body...)
}

func GenerateOverride(receiver jen.Code, importPath, outputImportPath string, method *types.Method, body ...jen.Code) jen.Code {
	params := GenerateParamTypes(method, importPath, outputImportPath, false)
	for i, param := range params {
		params[i] = Compose(jen.Id(fmt.Sprintf("v%d", i)), param)
	}

	return GenerateMethod(
		receiver,
		method.Name,
		params,
		GenerateResultTypes(method, importPath, outputImportPath),
		body...,
	)
}

func GenerateParamTypes(method *types.Method, importPath, outputImportPath string, omitDots bool) []jen.Code {
	params := []jen.Code{}
	for i, typ := range method.Params {
		params = append(params, GenerateType(
			typ,
			importPath,
			outputImportPath,
			method.Variadic && i == len(method.Params)-1 && !omitDots,
		))
	}

	return params
}

func GenerateResultTypes(method *types.Method, importPath, outputImportPath string) []jen.Code {
	results := []jen.Code{}
	for _, typ := range method.Results {
		results = append(results, GenerateType(
			typ,
			importPath,
			outputImportPath,
			false,
		))
	}

	return results
}

func GenerateDecoratedCall(method *types.Method, target *jen.Statement) jen.Code {
	names := []jen.Code{}
	for i := range method.Params {
		name := jen.Id(fmt.Sprintf("v%d", i))
		if method.Variadic && i == len(method.Params)-1 {
			name = Compose(name, jen.Op("..."))
		}

		names = append(names, name)
	}

	dispatch := target.Call(names...)
	if len(method.Results) == 0 {
		return dispatch
	}

	assign := jen.Id("r0")
	for i := 1; i < len(method.Results); i++ {
		assign = assign.Op(",").Id(fmt.Sprintf("r%d", i))
	}

	return Compose(assign.Op(":="), dispatch)
}

func GenerateDecoratedReturn(method *types.Method) jen.Code {
	ret := jen.Return()

	if len(method.Results) > 0 {
		ret = ret.Id("r0")

		for i := 1; i < len(method.Results); i++ {
			ret = ret.Op(",").Id(fmt.Sprintf("r%d", i))
		}
	}

	return ret
}
