package generation

import (
	"fmt"

	"github.com/dave/jennifer/jen"

	"github.com/efritz/go-mockgen/specs"
)

const (
	// Structs
	mockFormat     = "Mock%s%s"
	paramSetFormat = "%s%s%sParamSet"

	// Fields
	innerMethodFormat           = "%sFunc"
	statCallParamsFormat        = "hist%s"
	statsMutexName              = "mutex"
	exportedParameterNameFormat = "Arg%d"

	// Methods
	callCountFormat  = "%sFuncCallCount"
	callParamsFormat = "%sFuncCallParams"

	// Variables
	parameterNameFormat = "v%d"

	// Functions
	constructorFormat   = "NewMock%s%s"
	defaultMethodFormat = "default%sFunc"
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
		g.generateParamSetDefinitions,
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

	for _, name := range g.spec.MethodNames() {
		var (
			method             = g.spec.Method(name)
			methodName         = fmt.Sprintf(innerMethodFormat, name)
			params             = generateParams(method, g.spec.ImportPath, false)
			results            = generateResults(method, g.spec.ImportPath)
			statCallParamsName = fmt.Sprintf(statCallParamsFormat, name)
			callParamSetName   = fmt.Sprintf(paramSetFormat, g.prefix, g.name, name)
		)

		// {{Method}}Func func({{params...}}) {{results...}}
		fields = append(fields, jen.Id(methodName).Func().Params(params...).Params(results...))

		// stat{{Method}}FuncParams []{{Method}}FuncParamSet
		fields = append(fields, jen.Id(statCallParamsName).Index().Id(callParamSetName))
	}

	// Lock sync.RWMutex
	fields = append(fields, jen.Id(statsMutexName).Qual("sync", "RWMutex"))

	// type Mock{{Interface}} struct { [fields] }
	g.file.Type().Id(name).Struct(fields...)
}

//
// ParamSet Definitions

func (g *interfaceGenerator) generateParamSetDefinitions() {
	for _, name := range g.spec.MethodNames() {
		g.generateParamSetDefinition(name, g.spec.Method(name))
	}
}

func (g *interfaceGenerator) generateParamSetDefinition(name string, method *specs.MethodSpec) {
	var (
		structName = fmt.Sprintf(paramSetFormat, g.prefix, g.name, name)
		params     = generateParams(method, g.spec.ImportPath, true)
		paramNames = generateExportedParamNames(method)
		fields     = []jen.Code{}
	)

	for i, param := range params {
		fields = append(fields, jen.Add(paramNames[i], param))
	}

	// type {{Interface}}{{Method}}FuncParamSet struct { [fields] }
	g.file.Type().Id(structName).Struct(fields...)
}

//
// Constructor

func (g *interfaceGenerator) generateConstructor() {
	var (
		structName      = fmt.Sprintf(mockFormat, g.prefix, g.name)
		constructorName = fmt.Sprintf(constructorFormat, g.prefix, g.name)
		initStatement   = jen.Id("m").Op(":=").Op("&").Id(structName).Values()
		returnStatement = jen.Return().Id("m")
		fields          = []jen.Code{}
	)

	for _, name := range g.spec.MethodNames() {
		var (
			methodName        = fmt.Sprintf(innerMethodFormat, name)
			defaultMethodName = fmt.Sprintf(defaultMethodFormat, name)
		)

		// m.{{Method}}Func = m.default{{Method}}Func
		fields = append(fields, jen.Id("m").Dot(methodName).Op("=").Id("m").Dot(defaultMethodName))
	}

	// m := &Mock{{Interface}}{}; [fields] ; return m
	body := append([]jen.Code{initStatement}, append(fields, returnStatement)...)

	// func NewMock{{Interface}} *Mock{{Interface}} { [body] }
	g.file.Func().Id(constructorName).Params().Op("*").Id(structName).Block(body...)
}

//
// Method Implementations

func (g *interfaceGenerator) generateMethodImplementations() {
	for _, name := range g.spec.MethodNames() {
		g.generateMethodImplementation(name, g.spec.Method(name))
		g.generateStatMethodImplementations(name)
		g.file.Line()
	}
}

func (g *interfaceGenerator) generateStatMethodImplementations(name string) {
	var (
		structName         = fmt.Sprintf(mockFormat, g.prefix, g.name)
		statCallParamsName = fmt.Sprintf(statCallParamsFormat, name)
		callCountName      = fmt.Sprintf(callCountFormat, name)
		callParamsName     = fmt.Sprintf(callParamsFormat, name)
		callParamSetName   = fmt.Sprintf(paramSetFormat, g.prefix, g.name, name)
	)

	g.file.Func().Params(jen.Id("m").Op("*").Id(structName)).Id(callCountName).Params().Params(jen.Int()).Block(
		jen.Id("m").Dot(statsMutexName).Dot("RLock").Call(),
		jen.Defer().Id("m").Dot(statsMutexName).Dot("RUnlock").Call(),
		jen.Return(jen.Len(jen.Id("m").Dot(statCallParamsName))),
	)
	g.file.Func().Params(jen.Id("m").Op("*").Id(structName)).Id(callParamsName).
		Params().
		Params(jen.Index().Id(callParamSetName)).
		Block(
			jen.Id("m").Dot(statsMutexName).Dot("RLock").Call(),
			jen.Defer().Id("m").Dot(statsMutexName).Dot("RUnlock").Call(),
			jen.Return(jen.Id("m").Dot(statCallParamsName)),
		)
}

func (g *interfaceGenerator) generateMethodImplementation(name string, method *specs.MethodSpec) {
	var (
		structName = fmt.Sprintf(mockFormat, g.prefix, g.name)
		params     = generateParams(method, g.spec.ImportPath, false)
		results    = generateResults(method, g.spec.ImportPath)
		paramNames = generateParamNames(method, false)
		body       = g.generateMethodBody(name, method, paramNames)
	)

	for i, param := range params {
		params[i] = compose(jen.Id(fmt.Sprintf(parameterNameFormat, i)), param)
	}

	// func (m *Mock{{Interface}}) {{Method}}({{params...}}) {{results...}} { [body] }
	g.file.Func().Params(jen.Id("m").Op("*").Id(structName)).Id(name).Params(params...).Params(results...).Block(body...)
}

func (g *interfaceGenerator) generateMethodBody(name string, method *specs.MethodSpec, names []jen.Code) []jen.Code {
	var (
		methodName         = fmt.Sprintf(innerMethodFormat, name)
		statCallParamsName = fmt.Sprintf(statCallParamsFormat, name)
		paramNames         = generateParamNames(method, true)
		callParamSetName   = fmt.Sprintf(paramSetFormat, g.prefix, g.name, name)
	)

	// m.stats{{Method}}.Lock()
	lock := jen.Id("m").Dot(statsMutexName).Dot("Lock").Call()

	// m.{{Method}}FuncCallParams = append(m.{{Method}}FuncCallParams, {params...})
	params := jen.Id("m").Dot(statCallParamsName).Op("=").Id("append").Call(
		jen.Id("m").Dot(statCallParamsName),
		jen.Id(callParamSetName).Values(paramNames...),
	)

	// m.{{Method}}Func({{params...}})
	dispatch := jen.Id("m").Dot(methodName).Call(names...)

	// m.stats{{Method}}.Unlock()
	unlock := jen.Id("m").Dot(statsMutexName).Dot("Unlock").Call()

	if len(method.Results) != 0 {
		// return [dispatch]
		dispatch = compose(jen.Return(), dispatch)
	}

	return []jen.Code{lock, params, unlock, dispatch}
}

//
// Default Method Implementations

func (g *interfaceGenerator) generateDefaultMethodImplementations() {
	for _, name := range g.spec.MethodNames() {
		g.generateDefaultMethodImplementation(
			fmt.Sprintf(defaultMethodFormat, name),
			g.spec.Method(name),
		)
	}
}

func (g *interfaceGenerator) generateDefaultMethodImplementation(name string, method *specs.MethodSpec) {
	params := generateParams(method, g.spec.ImportPath, false)
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

func generateParams(method *specs.MethodSpec, importPath string, omitDots bool) []jen.Code {
	params := []jen.Code{}
	for i, typ := range method.Params {
		params = append(params, generateType(
			typ,
			importPath,
			method.Variadic && i == len(method.Params)-1 && !omitDots,
		))
	}

	return params
}

func generateParamNames(method *specs.MethodSpec, omitDots bool) []jen.Code {
	return generateParamNamesFormat(method, omitDots, parameterNameFormat)
}

func generateExportedParamNames(method *specs.MethodSpec) []jen.Code {
	return generateParamNamesFormat(method, true, exportedParameterNameFormat)
}

func generateParamNamesFormat(method *specs.MethodSpec, omitDots bool, format string) []jen.Code {
	names := []jen.Code{}
	for i := range method.Params {
		name := jen.Id(fmt.Sprintf(format, i))

		if method.Variadic && i == len(method.Params)-1 && !omitDots {
			name = name.Op("...")
		}

		names = append(names, name)
	}

	return names
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
