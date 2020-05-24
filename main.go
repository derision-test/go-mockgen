package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/dustin/go-humanize"
	"github.com/efritz/go-genlib/command"
	"github.com/efritz/go-genlib/generation"
	"github.com/efritz/go-genlib/types"
)

type (
	generator struct {
		outputImportPath string
	}

	wrappedInterface struct {
		*types.Interface
		prefix         string
		titleName      string
		mockStructName string
		wrappedMethods []*wrappedMethod
	}

	wrappedMethod struct {
		*types.Method
		iface             *types.Interface
		dotlessParamTypes []jen.Code
		paramTypes        []jen.Code
		resultTypes       []jen.Code
		signature         jen.Code
	}

	topLevelGenerator func(*wrappedInterface) jen.Code
	methodGenerator   func(*wrappedInterface, *wrappedMethod) jen.Code
)

const (
	name        = "go-mockgen"
	packageName = "github.com/efritz/go-mockgen"
	description = "go-mockgen generates mock implementations from interface definitions."
	version     = "0.1.0"

	mockStructFormat  = "Mock%s%s"
	funcStructFormat  = "%s%s%sFunc"
	callStructFormat  = "%s%s%sFuncCall"
	funcFieldFormat   = "%sFunc"
	argFieldFormat    = "Arg%d"
	resultFieldFormat = "Result%d"
	argVarFormat      = "v%d"
	resultVarFormat   = "r%d"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := command.Run(name, description, version, types.GetInterface, generate); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func generate(ifaces []*types.Interface, opts *command.Options) error {
	g := &generator{
		outputImportPath: opts.OutputImportPath,
	}

	return generation.Generate(
		packageName,
		version,
		ifaces,
		opts,
		generateFilename,
		g.generateInterface,
	)
}

func generateFilename(name string) string {
	return fmt.Sprintf("%s_mock.go", name)
}

func (g *generator) generateInterface(file *jen.File, iface *types.Interface, prefix string) {
	var (
		titleName        = title(iface.Name)
		mockStructName   = fmt.Sprintf(mockStructFormat, prefix, titleName)
		wrappedInterface = g.wrapInterface(iface, prefix, titleName, mockStructName)
	)

	topLevelGenerators := []topLevelGenerator{
		g.generateMockStruct,
		g.generateMockStructConstructor,
		g.generateMockStructFromConstructor,
	}

	methodGenerators := []methodGenerator{
		g.generateFuncStruct,
		g.generateFunc,
		g.generateFuncSetHookMethod,
		g.generateFuncPushHookMethod,
		g.generateFuncSetReturnMethod,
		g.generateFuncPushReturnMethod,
		g.generateFuncNextHookMethod,
		g.generateFuncAppendCallMethod,
		g.generateFuncHistoryMethod,
		g.generateCallStruct,
		g.generateCallArgsMethod,
		g.generateCallResultsMethod,
	}

	for _, generator := range topLevelGenerators {
		file.Add(generator(wrappedInterface))
		file.Line()
	}

	for _, method := range wrappedInterface.wrappedMethods {
		for _, generator := range methodGenerators {
			file.Add(generator(wrappedInterface, method))
			file.Line()
		}
	}
}

func (g *generator) wrapInterface(iface *types.Interface, prefix, titleName, mockStructName string) *wrappedInterface {
	wrapped := &wrappedInterface{
		Interface:      iface,
		prefix:         prefix,
		titleName:      titleName,
		mockStructName: mockStructName,
	}

	for _, method := range iface.Methods {
		wrapped.wrappedMethods = append(wrapped.wrappedMethods, g.wrapMethod(iface, method))
	}

	return wrapped
}

func (g *generator) wrapMethod(iface *types.Interface, method *types.Method) *wrappedMethod {
	m := &wrappedMethod{
		Method:            method,
		iface:             iface,
		dotlessParamTypes: generation.GenerateParamTypes(method, iface.ImportPath, g.outputImportPath, true),
		paramTypes:        generation.GenerateParamTypes(method, iface.ImportPath, g.outputImportPath, false),
		resultTypes:       generation.GenerateResultTypes(method, iface.ImportPath, g.outputImportPath),
	}

	m.signature = jen.Func().Params(m.paramTypes...).Params(m.resultTypes...)
	return m
}

//
// Mock Struct Generation

func (g *generator) generateMockStruct(iface *wrappedInterface) jen.Code {
	comment := generation.GenerateComment(
		1,
		"%s is a mock implementation of the %s interface (from the package %s) used for unit testing.",
		iface.mockStructName,
		iface.Name,
		iface.ImportPath,
	)

	structFields := []jen.Code{}
	for _, method := range iface.Methods {
		name := fmt.Sprintf(funcFieldFormat, method.Name)

		comment := generation.GenerateComment(
			2,
			"%s is an instance of a mock function object controlling the behavior of the method %s.",
			name,
			method.Name,
		)

		hookFuncField := comment.
			Id(name).
			Op("*").
			Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name))

		structFields = append(structFields, hookFuncField)
	}

	return comment.
		Type().
		Id(iface.mockStructName).
		Struct(structFields...)
}

//
// Constructor Generation

func (g *generator) generateMockStructConstructor(iface *wrappedInterface) jen.Code {
	name := fmt.Sprintf("New%s", iface.mockStructName)

	comment := generation.GenerateComment(
		1,
		"%s creates a new mock of the %s interface. All methods return zero values for all results, unless overwritten.",
		name,
		iface.Name,
	)

	constructorFields := []jen.Code{}
	for _, method := range iface.wrappedMethods {
		zeroes := []jen.Code{}
		for _, typ := range method.Results {
			zeroes = append(zeroes, generation.GenerateZeroValue(
				typ,
				iface.ImportPath,
				g.outputImportPath,
			))
		}

		zeroFunction := generation.GenerateFunction(
			"",
			method.paramTypes,
			generation.GenerateResultTypes(method.Method, iface.ImportPath, g.outputImportPath),
			jen.Return().List(zeroes...),
		)

		innerStructField := generation.Compose(
			jen.Line(),
			generation.Compose(jen.Id("defaultHook").Op(":"), zeroFunction),
		)

		field := jen.
			Line().
			Id(fmt.Sprintf(funcFieldFormat, method.Name)).
			Op(":").
			Op("&").
			Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)).
			Values(innerStructField, jen.Line())

		constructorFields = append(constructorFields, field)
	}

	constructorFields = append(constructorFields, jen.Line())

	functionDecl := generation.GenerateFunction(
		name,
		nil,
		[]jen.Code{jen.Op("*").Id(iface.mockStructName)},
		jen.Return().Op("&").Id(iface.mockStructName).Values(constructorFields...),
	)

	return generation.Compose(comment, functionDecl)
}

func (g *generator) generateMockStructFromConstructor(iface *wrappedInterface) jen.Code {
	ifaceName := jen.Qual(generation.SanitizeImportPath(iface.ImportPath, g.outputImportPath), iface.Name)

	var surrogate *jen.Statement
	if !unicode.IsUpper([]rune(iface.Name)[0]) {
		name := fmt.Sprintf("surrogateMock%s", iface.titleName)

		comment := generation.GenerateComment(
			1,
			"%s is a copy of the %s interface (from the package %s). It is redefined here as it is unexported in the source packge.",
			name,
			iface.Name,
			iface.ImportPath,
		)

		signatures := []jen.Code{}
		for _, method := range iface.wrappedMethods {
			signatures = append(signatures, jen.Id(method.Name).Params(method.paramTypes...).Params(method.resultTypes...))
		}

		ifaceName = jen.Id(name)
		surrogate = comment.Type().Id(name).Interface(signatures...).Line()
	}

	name := fmt.Sprintf("New%sFrom", iface.mockStructName)

	comment := generation.GenerateComment(
		1,
		"%s creates a new mock of the %s interface. All methods delegate to the given implementation, unless overwritten.",
		name,
		iface.mockStructName,
	)

	constructorFields := []jen.Code{}
	for _, method := range iface.Methods {
		innerStructField := generation.Compose(
			jen.Line(),
			generation.Compose(jen.Id("defaultHook").Op(":"), jen.Id("i").Dot(method.Name)),
		)

		field := jen.
			Line().
			Id(fmt.Sprintf(funcFieldFormat, method.Name)).
			Op(":").
			Op("&").
			Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)).
			Values(innerStructField, jen.Line())

		constructorFields = append(constructorFields, field)
	}

	constructorFields = append(constructorFields, jen.Line())

	functionDecl := generation.GenerateFunction(
		name,
		[]jen.Code{generation.Compose(jen.Id("i"), ifaceName)},
		[]jen.Code{jen.Op("*").Id(iface.mockStructName)},
		jen.Return().Op("&").Id(iface.mockStructName).Values(constructorFields...),
	)

	if surrogate != nil {
		comment = generation.Compose(surrogate, comment)
	}

	return generation.Compose(comment, functionDecl)
}

//
// Func Struct Generation

func (g *generator) generateFuncStruct(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	name := fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)

	comment := generation.GenerateComment(
		1,
		"%s describes the behavior when the %s method of the parent %s instance is invoked.",
		name,
		method.Name,
		iface.mockStructName,
	)

	defaultHookField := generation.Compose(jen.Id("defaultHook"), method.signature)
	hooksField := generation.Compose(jen.Id("hooks").Index(), method.signature)
	mutexField := jen.Id("mutex").Qual("sync", "Mutex")

	historyField := jen.
		Id("history").
		Index().
		Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name))

	return comment.
		Type().
		Id(name).
		Struct(defaultHookField, hooksField, historyField, mutexField)
}

func (g *generator) generateFunc(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	comment := generation.GenerateComment(
		1,
		"%s delegates to the next hook function in the queue and stores the parameter and result values of this invocation.",
		method.Name,
	)

	hook := jen.
		Id("m").
		Dot(fmt.Sprintf(funcFieldFormat, method.Name)).
		Dot("nextHook").
		Call()

	callInstanceStructName := fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)

	names := []jen.Code{}
	for i := 0; i < len(method.Params); i++ {
		names = append(names, jen.Id(fmt.Sprintf(argVarFormat, i)))
	}

	for i := 0; i < len(method.Results); i++ {
		names = append(names, jen.Id(fmt.Sprintf(resultVarFormat, i)))
	}

	appendFuncCall := jen.
		Id("m").
		Dot(fmt.Sprintf(funcFieldFormat, method.Name)).
		Dot("appendCall").
		Call(jen.Id(callInstanceStructName).Values(names...))

	methodDecl := generation.GenerateOverride(
		jen.Id("m").Op("*").Id(iface.mockStructName),
		method.iface.ImportPath,
		g.outputImportPath,
		method.Method,
		generation.GenerateDecoratedCall(method.Method, hook),
		appendFuncCall,
		generation.GenerateDecoratedReturn(method.Method),
	)

	return generation.Compose(comment, methodDecl)
}

func (g *generator) generateFuncSetHookMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	comment := generation.GenerateComment(
		1,
		"SetDefaultHook sets function that is called when the %s method of the parent %s instance is invoked and the hook queue is empty.",
		method.Name,
		iface.mockStructName,
	)

	methodDecl := generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"SetDefaultHook",
		[]jen.Code{generation.Compose(jen.Id("hook"), method.signature)},
		nil,
		jen.Id("f").Dot("defaultHook").Op("=").Id("hook"),
	)

	return generation.Compose(comment, methodDecl)
}

func (g *generator) generateFuncPushHookMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	comment := generation.GenerateComment(
		1,
		"PushHook adds a function to the end of hook queue. Each invocation of the %s method of the parent %s instance inovkes the hook at the front of the queue and discards it. After the queue is empty, the default hook function is invoked for any future action.",
		method.Name,
		iface.mockStructName,
	)

	lock := jen.
		Id("f").
		Dot("mutex").
		Dot("Lock").
		Call()

	unlock := jen.
		Id("f").
		Dot("mutex").
		Dot("Unlock").
		Call()

	methodDecl := generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"PushHook",
		[]jen.Code{generation.Compose(jen.Id("hook"), method.signature)},
		nil,
		lock,
		selfAppend(jen.Id("f").Dot("hooks"), jen.Id("hook")),
		unlock,
	)

	return generation.Compose(comment, methodDecl)
}

func (g *generator) generateFuncSetReturnMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	return g.generateReturnMethod(iface, method, "SetDefault")
}

func (g *generator) generateFuncPushReturnMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	return g.generateReturnMethod(iface, method, "Push")
}

func (g *generator) generateReturnMethod(iface *wrappedInterface, method *wrappedMethod, methodPrefix string) jen.Code {
	comment := generation.GenerateComment(
		1,
		"%sReturn calls %sDefaultHook with a function that returns the given values.",
		methodPrefix,
		methodPrefix,
	)

	names := []jen.Code{}
	namedResults := []jen.Code{}
	for i, t := range method.resultTypes {
		name := jen.Id(fmt.Sprintf(resultVarFormat, i))
		names = append(names, name)
		namedResults = append(namedResults, generation.Compose(name, t))
	}

	function := generation.GenerateFunction(
		"",
		method.paramTypes,
		method.resultTypes,
		jen.Return().List(names...),
	)

	methodDecl := generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		fmt.Sprintf("%sReturn", methodPrefix),
		namedResults,
		nil,
		jen.Id("f").Dot(fmt.Sprintf("%sHook", methodPrefix)).Call(function),
	)

	return generation.Compose(comment, methodDecl)
}

func (g *generator) generateFuncNextHookMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	returnDefaultIfEmpty := jen.
		If(jen.Len(jen.Id("f").Dot("hooks")).Op("==").Lit(0)).
		Block(jen.Return(jen.Id("f").Dot("defaultHook")))

	lock := jen.
		Id("f").
		Dot("mutex").
		Dot("Lock").
		Call()

	deferUnlock := jen.
		Defer().
		Id("f").
		Dot("mutex").
		Dot("Unlock").
		Call()

	getFirstHook := jen.
		Id("hook").
		Op(":=").
		Id("f").
		Dot("hooks").
		Index(jen.Lit(0))

	popHook := jen.
		Id("f").
		Dot("hooks").
		Op("=").
		Id("f").
		Dot("hooks").
		Index(jen.Lit(1).Op(":"))

	return generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"nextHook",
		nil,
		[]jen.Code{method.signature},
		lock,
		deferUnlock,
		jen.Line(),
		returnDefaultIfEmpty,
		jen.Line(),
		getFirstHook,
		popHook,
		jen.Return(jen.Id("hook")),
	)
}

func (g *generator) generateFuncAppendCallMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	lock := jen.
		Id("f").
		Dot("mutex").
		Dot("Lock").
		Call()

	unlock := jen.
		Id("f").
		Dot("mutex").
		Dot("Unlock").
		Call()

	return generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"appendCall",
		[]jen.Code{generation.Compose(jen.Id("r0"), jen.Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)))},
		nil,
		lock,
		selfAppend(jen.Id("f").Dot("history"), jen.Id("r0")),
		unlock,
	)
}

func (g *generator) generateFuncHistoryMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	comment := generation.GenerateComment(
		1,
		"History returns a sequence of %s objects describing the invocations of this function.",
		fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name),
	)

	lock := jen.
		Id("f").
		Dot("mutex").
		Dot("Lock").
		Call()

	unlock := jen.
		Id("f").
		Dot("mutex").
		Dot("Unlock").
		Call()

	declareCopy := jen.
		Id("history").
		Op(":=").
		Make(
			jen.Index().Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)),
			jen.Len(jen.Id("f").Dot("history")),
		)

	doCopy := jen.Copy(
		jen.Id("history"),
		jen.Id("f").Dot("history"),
	)

	methodDecl := generation.GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"History",
		nil,
		[]jen.Code{jen.Index().Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name))},
		lock,
		declareCopy,
		doCopy,
		unlock,
		jen.Line(),
		jen.Return().Id("history"),
	)

	return generation.Compose(comment, methodDecl)
}

//
// Call Struct Generation

func (g *generator) generateCallStruct(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	name := fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)

	comment := generation.GenerateComment(
		1,
		"%s is an object that describes an invocation of method %s on an instance of %s.",
		name,
		method.Name,
		iface.mockStructName,
	)

	fields := []jen.Code{}
	for i, param := range method.dotlessParamTypes {
		name := fmt.Sprintf(argFieldFormat, i)

		var commentText string
		if i == len(method.dotlessParamTypes)-1 && method.Variadic {
			commentText = fmt.Sprintf(
				"%s is a slice containing the values of the variadic arguments passed to this method invocation.",
				name,
			)
		} else {
			commentText = fmt.Sprintf(
				"%s is the value of the %s argument passed to this method invocation.",
				name,
				humanize.Ordinal(i+1),
			)
		}

		fields = append(fields, generation.GenerateComment(2, commentText).Id(name).Add(param))
	}

	for i, param := range method.resultTypes {
		name := fmt.Sprintf(resultFieldFormat, i)

		comment := generation.GenerateComment(
			2,
			"%s is the value of the %s result returned from this method invocation.",
			name,
			humanize.Ordinal(i+1),
		)

		fields = append(fields, comment.Id(name).Add(param))
	}

	return comment.
		Type().
		Id(name).
		Struct(fields...)
}

func (g *generator) generateCallArgsMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	var commentText string
	if method.Variadic {
		commentText = "Args returns an interface slice containing the arguments of this invocation. The variadic slice argument is flattened in this array such that one positional argument and three variadic arguments would result in a slice of four, not two."
	} else {
		commentText = "Args returns an interface slice containing the arguments of this invocation."
	}

	comment := generation.GenerateComment(
		1,
		commentText,
	)

	values := []jen.Code{}
	for i := range method.Params {
		values = append(values, jen.Id("c").Dot(fmt.Sprintf(argFieldFormat, i)))
	}

	var body jen.Code

	if method.Variadic {
		var (
			lastIndex         = len(values) - 1
			nonVariadicValues = values[:lastIndex]
		)

		body = jen.
			Id("trailing").
			Op(":=").
			Index().
			Interface().
			Values().
			Line().
			For(generation.Compose(jen.Id("_").Op(",").Id("val").Op(":=").Range(), values[lastIndex])).
			Block(selfAppend(jen.Id("trailing"), jen.Id("val"))).
			Line().
			Line().
			Return().
			Append(jen.Index().Interface().Values(nonVariadicValues...), jen.Id("trailing").Op("..."))
	} else {
		body = jen.Return().Index().Interface().Values(values...)
	}

	methodDecl := generation.GenerateMethod(
		jen.Id("c").Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)),
		"Args",
		nil,
		[]jen.Code{jen.Index().Interface()},
		body,
	)

	return generation.Compose(comment, methodDecl)
}

func (g *generator) generateCallResultsMethod(iface *wrappedInterface, method *wrappedMethod) jen.Code {
	comment := generation.GenerateComment(
		1,
		"Results returns an interface slice containing the results of this invocation.",
	)

	values := []jen.Code{}
	for i := range method.Results {
		values = append(values, jen.Id("c").Dot(fmt.Sprintf(resultFieldFormat, i)))
	}

	methodDecl := generation.GenerateMethod(
		jen.Id("c").Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)),
		"Results",
		nil,
		[]jen.Code{jen.Index().Interface()},
		jen.Return().Index().Interface().Values(values...),
	)

	return generation.Compose(comment, methodDecl)
}

//
// Helpers

func title(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

func selfAppend(sliceRef *jen.Statement, value jen.Code) jen.Code {
	return generation.Compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}
