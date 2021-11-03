package generation

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/dave/jennifer/jen"
	"github.com/derision-test/go-mockgen/internal/mockgen/types"
	"github.com/dustin/go-humanize"
)

type wrappedInterface struct {
	*types.Interface
	prefix         string
	titleName      string
	mockStructName string
	wrappedMethods []*wrappedMethod
}

type wrappedMethod struct {
	*types.Method
	iface             *types.Interface
	dotlessParamTypes []jen.Code
	paramTypes        []jen.Code
	resultTypes       []jen.Code
	signature         jen.Code
}

type (
	topLevelGenerator func(*wrappedInterface, string) jen.Code
	methodGenerator   func(*wrappedInterface, *wrappedMethod, string) jen.Code
)

const (
	mockStructFormat  = "Mock%s%s"
	funcStructFormat  = "%s%s%sFunc"
	callStructFormat  = "%s%s%sFuncCall"
	funcFieldFormat   = "%sFunc"
	argFieldFormat    = "Arg%d"
	resultFieldFormat = "Result%d"
	argVarFormat      = "v%d"
	resultVarFormat   = "r%d"
)

func generateInterface(file *jen.File, iface *types.Interface, prefix, outputImportPath string) {
	titleName := title(iface.Name)
	mockStructName := fmt.Sprintf(mockStructFormat, prefix, titleName)
	wrappedInterface := wrapInterface(iface, prefix, titleName, mockStructName, outputImportPath)

	topLevelGenerators := []topLevelGenerator{
		generateMockStruct,
		generateMockStructConstructor,
		generateMockStructStrictConstructor,
		generateMockStructFromConstructor,
	}

	methodGenerators := []methodGenerator{
		generateFuncStruct,
		generateFunc,
		generateFuncSetHookMethod,
		generateFuncPushHookMethod,
		generateFuncSetReturnMethod,
		generateFuncPushReturnMethod,
		generateFuncNextHookMethod,
		generateFuncAppendCallMethod,
		generateFuncHistoryMethod,
		generateCallStruct,
		generateCallArgsMethod,
		generateCallResultsMethod,
	}

	for _, generator := range topLevelGenerators {
		file.Add(generator(wrappedInterface, outputImportPath))
		file.Line()
	}

	for _, method := range wrappedInterface.wrappedMethods {
		for _, generator := range methodGenerators {
			file.Add(generator(wrappedInterface, method, outputImportPath))
			file.Line()
		}
	}
}

func wrapInterface(iface *types.Interface, prefix, titleName, mockStructName, outputImportPath string) *wrappedInterface {
	wrapped := &wrappedInterface{
		Interface:      iface,
		prefix:         prefix,
		titleName:      titleName,
		mockStructName: mockStructName,
	}

	for _, method := range iface.Methods {
		wrapped.wrappedMethods = append(wrapped.wrappedMethods, wrapMethod(iface, method, outputImportPath))
	}

	return wrapped
}

func wrapMethod(iface *types.Interface, method *types.Method, outputImportPath string) *wrappedMethod {
	m := &wrappedMethod{
		Method:            method,
		iface:             iface,
		dotlessParamTypes: GenerateParamTypes(method, iface.ImportPath, outputImportPath, true),
		paramTypes:        GenerateParamTypes(method, iface.ImportPath, outputImportPath, false),
		resultTypes:       GenerateResultTypes(method, iface.ImportPath, outputImportPath),
	}

	m.signature = jen.Func().Params(m.paramTypes...).Params(m.resultTypes...)
	return m
}

//
// Mock Struct Generation

func generateMockStruct(iface *wrappedInterface, outputImportPath string) jen.Code {
	commentFmt := "%s is a mock implementation of the %s interface (from the package %s) used for unit testing."
	comment := GenerateComment(1, commentFmt, iface.mockStructName, iface.Name, iface.ImportPath)

	structFields := []jen.Code{}
	for _, method := range iface.Methods {
		name := fmt.Sprintf(funcFieldFormat, method.Name)

		commentFmt := "%s is an instance of a mock function object controlling the behavior of the method %s."
		comment := GenerateComment(2, commentFmt, name, method.Name)

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

func generateMockStructConstructor(iface *wrappedInterface, outputImportPath string) jen.Code {
	name := fmt.Sprintf("New%s", iface.mockStructName)
	commentFmt := "%s creates a new mock of the %s interface. All methods return zero values for all results, unless overwritten."
	comment := GenerateComment(1, commentFmt, name, iface.Name)

	constructorFields := []jen.Code{}
	for _, method := range iface.wrappedMethods {
		zeroes := []jen.Code{}
		for _, typ := range method.Results {
			zeroes = append(zeroes, GenerateZeroValue(
				typ,
				iface.ImportPath,
				outputImportPath,
			))
		}

		zeroFunction := GenerateFunction(
			"",
			method.paramTypes,
			GenerateResultTypes(method.Method, iface.ImportPath, outputImportPath),
			jen.Return().List(zeroes...),
		)

		innerStructField := Compose(
			jen.Line(),
			Compose(jen.Id("defaultHook").Op(":"), zeroFunction),
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

	functionDecl := GenerateFunction(
		name,
		nil,
		[]jen.Code{jen.Op("*").Id(iface.mockStructName)},
		jen.Return().Op("&").Id(iface.mockStructName).Values(constructorFields...),
	)

	return Compose(comment, functionDecl)
}

func generateMockStructStrictConstructor(iface *wrappedInterface, outputImportPath string) jen.Code {
	name := fmt.Sprintf("NewStrict%s", iface.mockStructName)
	commentFmt := "%s creates a new mock of the %s interface. All methods panic on invocation, unless overwritten."
	comment := GenerateComment(1, commentFmt, name, iface.Name)

	constructorFields := []jen.Code{}
	for _, method := range iface.wrappedMethods {
		zeroes := []jen.Code{}
		for _, typ := range method.Results {
			zeroes = append(zeroes, GenerateZeroValue(
				typ,
				iface.ImportPath,
				outputImportPath,
			))
		}

		panickingFunction := GenerateFunction(
			"",
			method.paramTypes,
			GenerateResultTypes(method.Method, iface.ImportPath, outputImportPath),
			jen.Panic(jen.Lit(fmt.Sprintf("unexpected invocation of %s.%s", iface.mockStructName, method.Method.Name))),
		)

		innerStructField := Compose(
			jen.Line(),
			Compose(jen.Id("defaultHook").Op(":"), panickingFunction),
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

	functionDecl := GenerateFunction(
		name,
		nil,
		[]jen.Code{jen.Op("*").Id(iface.mockStructName)},
		jen.Return().Op("&").Id(iface.mockStructName).Values(constructorFields...),
	)

	return Compose(comment, functionDecl)
}

func generateMockStructFromConstructor(iface *wrappedInterface, outputImportPath string) jen.Code {
	ifaceName := jen.Qual(SanitizeImportPath(iface.ImportPath, outputImportPath), iface.Name)

	var surrogate *jen.Statement
	if !unicode.IsUpper([]rune(iface.Name)[0]) {
		name := fmt.Sprintf("surrogateMock%s", iface.titleName)
		commentFmt := "%s is a copy of the %s interface (from the package %s). It is redefined here as it is unexported in the source package."
		comment := GenerateComment(1, commentFmt, name, iface.Name, iface.ImportPath)

		signatures := []jen.Code{}
		for _, method := range iface.wrappedMethods {
			signatures = append(signatures, jen.Id(method.Name).Params(method.paramTypes...).Params(method.resultTypes...))
		}

		ifaceName = jen.Id(name)
		surrogate = comment.Type().Id(name).Interface(signatures...).Line()
	}

	name := fmt.Sprintf("New%sFrom", iface.mockStructName)
	commentFmt := "%s creates a new mock of the %s interface. All methods delegate to the given implementation, unless overwritten."
	comment := GenerateComment(1, commentFmt, name, iface.mockStructName)

	constructorFields := []jen.Code{}
	for _, method := range iface.Methods {
		innerStructField := Compose(
			jen.Line(),
			Compose(jen.Id("defaultHook").Op(":"), jen.Id("i").Dot(method.Name)),
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

	functionDecl := GenerateFunction(
		name,
		[]jen.Code{Compose(jen.Id("i"), ifaceName)},
		[]jen.Code{jen.Op("*").Id(iface.mockStructName)},
		jen.Return().Op("&").Id(iface.mockStructName).Values(constructorFields...),
	)

	if surrogate != nil {
		comment = Compose(surrogate, comment)
	}

	return Compose(comment, functionDecl)
}

//
// Func Struct Generation

func generateFuncStruct(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	name := fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)
	commentFmt := "%s describes the behavior when the %s method of the parent %s instance is invoked."
	comment := GenerateComment(1, commentFmt, name, method.Name, iface.mockStructName)

	defaultHookField := Compose(jen.Id("defaultHook"), method.signature)
	hooksField := Compose(jen.Id("hooks").Index(), method.signature)
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

func generateFunc(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	commentFmt := "%s delegates to the next hook function in the queue and stores the parameter and result values of this invocation."
	comment := GenerateComment(1, commentFmt, method.Name)

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

	methodDecl := GenerateOverride(
		jen.Id("m").Op("*").Id(iface.mockStructName),
		method.iface.ImportPath,
		outputImportPath,
		method.Method,
		GenerateDecoratedCall(method.Method, hook),
		appendFuncCall,
		GenerateDecoratedReturn(method.Method),
	)

	return Compose(comment, methodDecl)
}

func generateFuncSetHookMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	commentFmt := "SetDefaultHook sets function that is called when the %s method of the parent %s instance is invoked and the hook queue is empty."
	comment := GenerateComment(1, commentFmt, method.Name, iface.mockStructName)

	methodDecl := GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"SetDefaultHook",
		[]jen.Code{Compose(jen.Id("hook"), method.signature)},
		nil,
		jen.Id("f").Dot("defaultHook").Op("=").Id("hook"),
	)

	return Compose(comment, methodDecl)
}

func generateFuncPushHookMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	commentFmt := "PushHook adds a function to the end of hook queue. Each invocation of the %s method of the parent %s instance invokes the hook at the front of the queue and discards it. After the queue is empty, the default hook function is invoked for any future action."
	comment := GenerateComment(1, commentFmt, method.Name, iface.mockStructName)

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

	methodDecl := GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"PushHook",
		[]jen.Code{Compose(jen.Id("hook"), method.signature)},
		nil,
		lock,
		selfAppend(jen.Id("f").Dot("hooks"), jen.Id("hook")),
		unlock,
	)

	return Compose(comment, methodDecl)
}

func generateFuncSetReturnMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	return generateReturnMethod(iface, method, "SetDefault")
}

func generateFuncPushReturnMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	return generateReturnMethod(iface, method, "Push")
}

func generateReturnMethod(iface *wrappedInterface, method *wrappedMethod, methodPrefix string) jen.Code {
	commentFmt := "%sReturn calls %sDefaultHook with a function that returns the given values."
	comment := GenerateComment(1, commentFmt, methodPrefix, methodPrefix)

	names := []jen.Code{}
	namedResults := []jen.Code{}
	for i, t := range method.resultTypes {
		name := jen.Id(fmt.Sprintf(resultVarFormat, i))
		names = append(names, name)
		namedResults = append(namedResults, Compose(name, t))
	}

	function := GenerateFunction(
		"",
		method.paramTypes,
		method.resultTypes,
		jen.Return().List(names...),
	)

	methodDecl := GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		fmt.Sprintf("%sReturn", methodPrefix),
		namedResults,
		nil,
		jen.Id("f").Dot(fmt.Sprintf("%sHook", methodPrefix)).Call(function),
	)

	return Compose(comment, methodDecl)
}

func generateFuncNextHookMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
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

	return GenerateMethod(
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

func generateFuncAppendCallMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
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

	return GenerateMethod(
		jen.Id("f").Op("*").Id(fmt.Sprintf(funcStructFormat, iface.prefix, iface.titleName, method.Name)),
		"appendCall",
		[]jen.Code{Compose(jen.Id("r0"), jen.Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)))},
		nil,
		lock,
		selfAppend(jen.Id("f").Dot("history"), jen.Id("r0")),
		unlock,
	)
}

func generateFuncHistoryMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	commentFmt := "History returns a sequence of %s objects describing the invocations of this function."
	comment := GenerateComment(1, commentFmt, fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name))

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

	methodDecl := GenerateMethod(
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

	return Compose(comment, methodDecl)
}

//
// Call Struct Generation

func generateCallStruct(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	name := fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)
	commentFmt := "%s is an object that describes an invocation of method %s on an instance of %s."
	comment := GenerateComment(1, commentFmt, name, method.Name, iface.mockStructName)

	fields := []jen.Code{}
	for i, param := range method.dotlessParamTypes {
		name := fmt.Sprintf(argFieldFormat, i)

		var commentText string
		if i == len(method.dotlessParamTypes)-1 && method.Variadic {
			commentText = fmt.Sprintf("%s is a slice containing the values of the variadic arguments passed to this method invocation.", name)
		} else {
			commentText = fmt.Sprintf("%s is the value of the %s argument passed to this method invocation.", name, humanize.Ordinal(i+1))
		}

		fields = append(fields, GenerateComment(2, commentText).Id(name).Add(param))
	}

	for i, param := range method.resultTypes {
		name := fmt.Sprintf(resultFieldFormat, i)
		commentFmt := "%s is the value of the %s result returned from this method invocation."
		comment := GenerateComment(2, commentFmt, name, humanize.Ordinal(i+1))
		fields = append(fields, comment.Id(name).Add(param))
	}

	return comment.
		Type().
		Id(name).
		Struct(fields...)
}

func generateCallArgsMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	var commentText string
	if method.Variadic {
		commentText = "Args returns an interface slice containing the arguments of this invocation. The variadic slice argument is flattened in this array such that one positional argument and three variadic arguments would result in a slice of four, not two."
	} else {
		commentText = "Args returns an interface slice containing the arguments of this invocation."
	}

	comment := GenerateComment(1, commentText)

	values := []jen.Code{}
	for i := range method.Params {
		values = append(values, jen.Id("c").Dot(fmt.Sprintf(argFieldFormat, i)))
	}

	var body jen.Code
	if method.Variadic {
		lastIndex := len(values) - 1
		nonVariadicValues := values[:lastIndex]

		body = jen.
			Id("trailing").
			Op(":=").
			Index().
			Interface().
			Values().
			Line().
			For(Compose(jen.Id("_").Op(",").Id("val").Op(":=").Range(), values[lastIndex])).
			Block(selfAppend(jen.Id("trailing"), jen.Id("val"))).
			Line().
			Line().
			Return().
			Append(jen.Index().Interface().Values(nonVariadicValues...), jen.Id("trailing").Op("..."))
	} else {
		body = jen.Return().Index().Interface().Values(values...)
	}

	methodDecl := GenerateMethod(
		jen.Id("c").Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)),
		"Args",
		nil,
		[]jen.Code{jen.Index().Interface()},
		body,
	)

	return Compose(comment, methodDecl)
}

func generateCallResultsMethod(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	comment := GenerateComment(1, "Results returns an interface slice containing the results of this invocation.")

	values := []jen.Code{}
	for i := range method.Results {
		values = append(values, jen.Id("c").Dot(fmt.Sprintf(resultFieldFormat, i)))
	}

	methodDecl := GenerateMethod(
		jen.Id("c").Id(fmt.Sprintf(callStructFormat, iface.prefix, iface.titleName, method.Name)),
		"Results",
		nil,
		[]jen.Code{jen.Index().Interface()},
		jen.Return().Index().Interface().Values(values...),
	)

	return Compose(comment, methodDecl)
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
	return Compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}
