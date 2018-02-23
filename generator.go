package main

import (
	"bytes"
	"fmt"

	"github.com/dave/jennifer/jen"
)

const (
	mockFormat          = "Mock%s"
	constructorFormat   = "NewMock%s"
	innerMethodFormat   = "%sFunc"
	parameterNameFormat = "v%d"
)

func generate(specs map[string]*wrappedSpec) error {
	file := jen.NewFile("test")

	for name, spec := range specs {
		generateInterfaceDefinition(file, name, spec)
		generateTypeTest(file, name, spec)
		generateConstructor(file, name, spec)
		generateMethodImplementations(file, name, spec)
	}

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return err
	}

	fmt.Printf("%s\n", buffer.String())
	return nil
}

// generateInterfaceDefinition
//
// type Mock{{Interface}} struct {
//     {{Method}}Name func({{params...}}) {{results...}}
// }

func generateInterfaceDefinition(file *jen.File, interfaceName string, interfaceSpec *wrappedSpec) {
	fields := []jen.Code{}
	for methodName, method := range interfaceSpec.spec.methods {
		fields = append(fields, generateMethodField(
			methodName,
			method,
			interfaceSpec.importPath,
		))
	}

	file.Type().Id(fmt.Sprintf(mockFormat, interfaceName)).Struct(fields...)
}

func generateMethodField(methodName string, method *methodSpec, importPath string) *jen.Statement {
	return jen.Id(fmt.Sprintf(innerMethodFormat, methodName)).
		Func().
		Params(generateParams(method, importPath)...).
		Params(generateResults(method, importPath)...)
}

// generateTypeTest
//
// var _ {{Interface}} = NewMock{{Interface}}()

func generateTypeTest(file *jen.File, interfaceName string, interfaceSpec *wrappedSpec) {
	constructorName := fmt.Sprintf(constructorFormat, interfaceName)

	file.Var().
		Id("_").
		Qual(interfaceSpec.importPath, interfaceName).
		Op("=").
		Id(constructorName).
		Call()
}

// generateConstructor
//
// func NewMock{{Interface}} *Mock{{Interface}} {
//     return &Mock{{Interface}}{
//         {{Method}}Func func({{params...}}) {{results...}} { return {{result-zero-values...}} }
//     }
// }

func generateConstructor(file *jen.File, interfaceName string, interfaceSpec *wrappedSpec) {
	structName := fmt.Sprintf(mockFormat, interfaceName)
	constructorName := fmt.Sprintf(constructorFormat, interfaceName)

	body := jen.Return().
		Op("&").
		Id(structName).
		Values(generateDefaults(interfaceSpec.spec, interfaceSpec.importPath)...)

	file.Func().
		Id(constructorName).
		Params().
		Op("*").
		Id(structName).
		Block(body)
}

func generateDefaults(interfaceSpec *interfaceSpec, importPath string) []jen.Code {
	defaults := []jen.Code{}
	for methodName, method := range interfaceSpec.methods {
		defaults = append(defaults, generateDefault(method, methodName, importPath))
	}

	return defaults
}

func generateDefault(method *methodSpec, methodName, importPath string) *jen.Statement {
	zeroes := []jen.Code{}
	for _, typ := range method.results {
		zeroes = append(zeroes, zeroValue(
			typ,
			importPath,
		))
	}

	var body jen.Code
	if len(zeroes) != 0 {
		body = jen.Return().List(zeroes...)
	}

	defaultImpl := jen.Func().
		Params(generateParams(method, importPath)...).
		Params(generateResults(method, importPath)...).
		Block(body)

	return compose(jen.Id(fmt.Sprintf(innerMethodFormat, methodName)).Op(":"), defaultImpl)
}

// generateMethodImplementations
//
// func (m *Mock{{Interface}}) {{Method}}({{params...}}) {{results...}} {
//     return m.{{Method}}Func({{params...}})
// }

func generateMethodImplementations(file *jen.File, interfaceName string, interfaceSpec *wrappedSpec) {
	for methodName, method := range interfaceSpec.spec.methods {
		generateMethodImplementation(file, interfaceName, interfaceSpec.importPath, methodName, method)
	}
}

func generateMethodImplementation(file *jen.File, interfaceName string, importPath, methodName string, method *methodSpec) {
	names := []jen.Code{}
	for i := range method.params {
		name := jen.Id(fmt.Sprintf(parameterNameFormat, i))

		if method.variadic && i == len(method.params)-1 {
			name = name.Op("...")
		}

		names = append(names, name)
	}

	params := generateParams(method, importPath)
	for i, param := range params {
		params[i] = compose(jen.Id(fmt.Sprintf(parameterNameFormat, i)), param)
	}

	receiver := jen.Id("m").
		Op("*").
		Id(fmt.Sprintf(mockFormat, interfaceName))

	file.Func().
		Params(receiver).
		Id(methodName).
		Params(params...).
		Params(generateResults(method, importPath)...).
		Block(generateFunctionBody(methodName, method, names))
}

func generateFunctionBody(methodName string, method *methodSpec, names []jen.Code) *jen.Statement {
	body := jen.Id("m").
		Op(".").
		Id(fmt.Sprintf(innerMethodFormat, methodName)).
		Call(names...)

	if len(method.results) == 0 {
		return body
	}

	return compose(jen.Return(), body)
}

//
// Common Helpers

func generateParams(method *methodSpec, importPath string) []jen.Code {
	params := []jen.Code{}
	for i, typ := range method.params {
		params = append(params, generateType(
			typ,
			importPath,
			method.variadic && i == len(method.params)-1,
		))
	}

	return params
}

func generateResults(method *methodSpec, importPath string) []jen.Code {
	results := []jen.Code{}
	for _, typ := range method.results {
		results = append(results, generateType(
			typ,
			importPath,
			false,
		))
	}

	return results
}

func compose(stmt1 *jen.Statement, stmt2 jen.Code) *jen.Statement {
	composed := append(*stmt1, stmt2)
	return &composed
}
