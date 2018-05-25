package generation

import (
	"fmt"

	"github.com/dave/jennifer/jen"

	"github.com/efritz/go-mockgen/specs"
)

const (
	mockFormat          = "Mock%s%s"
	constructorFormat   = "NewMock%s%s"
	innerMethodFormat   = "%sFunc"
	callCountFormat     = "%sFuncCallCount"
	defaultMethodFormat = "default%sFunc"
	parameterNameFormat = "v%d"
)

type interfaceGenerator struct {
	file   *jen.File
	name   string
	prefix string
	spec   *specs.WrappedSpec
}

func newInterfaceGenerator(
	file *jen.File,
	name string,
	prefix string,
	spec *specs.WrappedSpec,
) *interfaceGenerator {
	return &interfaceGenerator{
		file:   file,
		name:   name,
		prefix: prefix,
		spec:   spec,
	}
}

//
// Code Generation

func (g *interfaceGenerator) generate() {
	fns := []func(){
		g.generateInterfaceDefinition,
		g.generateTypeTest,
		g.generateConstructor,
		g.generateMethodImplementations,
		g.generateDefaultMethodImplementations,
	}

	for _, fn := range fns {
		fn()
	}
}

//
// Definition

func (g *interfaceGenerator) generateInterfaceDefinition() {
	var (
		name   = fmt.Sprintf(mockFormat, g.prefix, g.name)
		fields = []jen.Code{}
	)

	for name, method := range g.spec.Spec.Methods {
		var (
			methodName    = fmt.Sprintf(innerMethodFormat, name)
			callCountName = fmt.Sprintf(callCountFormat, name)
			params        = generateParams(method, g.spec.ImportPath)
			results       = generateResults(method, g.spec.ImportPath)
		)

		// {{Method}}Name func({{params...}}) {{results...}}
		fields = append(fields, jen.Id(methodName).Func().Params(params...).Params(results...))

		// {{Method}}NameCallCOunt int
		fields = append(fields, jen.Id(callCountName).Id("int"))
	}

	// type Mock{{Interface}} struct { [fields] }
	g.file.Type().Id(name).Struct(fields...)
}

//
// Type Test

func (g *interfaceGenerator) generateTypeTest() {
	var (
		pkgName  = stripVendor(g.spec.ImportPath)
		ctorName = fmt.Sprintf(constructorFormat, g.prefix, g.name)
	)

	// var _ {{Interface}} = NewMock{{Interface}}()
	g.file.Var().Id("_").Qual(pkgName, g.name).Op("=").Id(ctorName).Call()
}

//
// Constructor

func (g *interfaceGenerator) generateConstructor() {
	var (
		structName      = fmt.Sprintf(mockFormat, g.prefix, g.name)
		constructorName = fmt.Sprintf(constructorFormat, g.prefix, g.name)
		initStatement   = jen.Id("m").Op(":=").Op("&").Id(structName).Values()
		returnStatement = jen.Return().Id("m")
		defaults        = []jen.Code{}
	)

	for name := range g.spec.Spec.Methods {
		var (
			methodName        = fmt.Sprintf(innerMethodFormat, name)
			defaultMethodName = fmt.Sprintf(defaultMethodFormat, name)
		)

		// m.{{Method}}Func = m.default{{Method}}Func
		defaults = append(defaults, jen.Id("m").Dot(methodName).Op("=").Id("m").Dot(defaultMethodName))
	}

	// m := &Mock{{Interface}}{}; [defaults] ; return m
	body := append([]jen.Code{initStatement}, append(defaults, returnStatement)...)

	// func NewMock{{Interface}} *Mock{{Interface}} { [body] }
	g.file.Func().Id(constructorName).Params().Op("*").Id(structName).Block(body...)
}

//
// Method Implementations

func (g *interfaceGenerator) generateMethodImplementations() {
	for methodName, method := range g.spec.Spec.Methods {
		g.generateMethodImplementation(methodName, method)
	}
}

func (g *interfaceGenerator) generateMethodImplementation(name string, method *specs.MethodSpec) {
	var (
		structName = fmt.Sprintf(mockFormat, g.prefix, g.name)
		params     = generateParams(method, g.name)
		results    = generateResults(method, g.spec.ImportPath)
		paramNames = g.generateParameterNames(method)
		body       = g.generateMethodBody(name, method, paramNames)
	)

	for i, param := range params {
		params[i] = compose(jen.Id(fmt.Sprintf(parameterNameFormat, i)), param)
	}

	// func (m *Mock{{Interface}}) {{Method}}({{params...}}) {{results...}} { [body] }
	g.file.Func().Params(jen.Id("m").Op("*").Id(structName)).Id(name).Params(params...).Params(results...).Block(body...)
}

func (g *interfaceGenerator) generateParameterNames(method *specs.MethodSpec) []jen.Code {
	names := []jen.Code{}
	for i := range method.Params {
		name := jen.Id(fmt.Sprintf(parameterNameFormat, i))

		if method.Variadic && i == len(method.Params)-1 {
			name = name.Op("...")
		}

		names = append(names, name)
	}

	return names
}

func (g *interfaceGenerator) generateMethodBody(name string, method *specs.MethodSpec, names []jen.Code) []jen.Code {
	var (
		methodName    = fmt.Sprintf(innerMethodFormat, name)
		callCountName = fmt.Sprintf(callCountFormat, name)
	)

	// m.{{Method}}FuncCallCount++
	incr := jen.Id("m").Dot(callCountName).Op("++")

	// m.{{Method}}Func({{params...}})
	dispatch := jen.Id("m").Dot(methodName).Call(names...)

	if len(method.Results) != 0 {
		// return [dispatch]
		dispatch = compose(jen.Return(), dispatch)
	}

	return []jen.Code{incr, dispatch}
}

//
// Default Method Implementations

func (g *interfaceGenerator) generateDefaultMethodImplementations() {
	for name, method := range g.spec.Spec.Methods {
		g.generateDefaultMethodImplementation(fmt.Sprintf(defaultMethodFormat, name), method)
	}
}

func (g *interfaceGenerator) generateDefaultMethodImplementation(name string, method *specs.MethodSpec) {
	params := generateParams(method, g.spec.ImportPath)
	for i, param := range params {
		params[i] = compose(jen.Id(fmt.Sprintf(parameterNameFormat, i)), param)
	}

	var (
		structName = fmt.Sprintf(mockFormat, g.prefix, g.name)
		receiver   = jen.Id("m").Op("*").Id(structName)
		results    = generateResults(method, g.spec.ImportPath)
		body       = g.generateDefaultMethodBody(method)
	)

	// func (m *Mock{{Interface}}) default{{Method}}({{params...}}) {{results...}} { [body] }
	g.file.Func().Params(receiver).Id(name).Params(params...).Params(results...).Block(body)
}

func (g *interfaceGenerator) generateDefaultMethodBody(method *specs.MethodSpec) jen.Code {
	zeroes := []jen.Code{}
	for _, typ := range method.Results {
		zeroes = append(zeroes, zeroValue(
			typ,
			g.spec.ImportPath,
		))
	}

	// return {{result-zero-values...}}
	return jen.Return().List(zeroes...)
}

//
// Common Helpers

func generateParams(method *specs.MethodSpec, importPath string) []jen.Code {
	params := []jen.Code{}
	for i, typ := range method.Params {
		params = append(params, generateType(
			typ,
			importPath,
			method.Variadic && i == len(method.Params)-1,
		))
	}

	return params
}

func generateResults(method *specs.MethodSpec, importPath string) []jen.Code {
	results := []jen.Code{}
	for _, typ := range method.Results {
		results = append(results, generateType(
			typ,
			importPath,
			false,
		))
	}

	return results
}
