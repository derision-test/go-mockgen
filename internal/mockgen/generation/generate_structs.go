package generation

import (
	"fmt"
	"strconv"

	"github.com/dave/jennifer/jen"
	"github.com/dustin/go-humanize"
)

func generateMockStruct(iface *wrappedInterface, outputImportPath string) jen.Code {
	mockStructName := iface.mockStructName
	commentText := fmt.Sprintf(
		`%s is a mock implementation of the %s interface (from the package %s) used for unit testing.`,
		mockStructName,
		iface.Name,
		iface.ImportPath,
	)

	structFields := make([]jen.Code, 0, len(iface.Methods))
	for _, method := range iface.wrappedMethods {
		mockFuncFieldName := fmt.Sprintf("%sFunc", method.Name)
		mockFuncStructName := fmt.Sprintf("%s%s%sFunc", iface.prefix, iface.titleName, method.Name)
		commentText := fmt.Sprintf(
			`%s is an instance of a mock function object controlling the behavior of the method %s.`,
			mockFuncFieldName,
			method.Name,
		)

		hook := jen.Id(mockFuncFieldName).Op("*").Id(mockFuncStructName)
		structFields = append(structFields, addComment(hook, 2, commentText))
	}

	// <Name>Func *<Prefix><InterfaceName><Name>Func, ...
	return generateStruct(mockStructName, commentText, structFields)
}

func generateMockFuncStruct(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	mockStructName := iface.mockStructName
	mockFuncStructName := fmt.Sprintf("%s%s%sFunc", iface.prefix, iface.titleName, method.Name)
	mockFuncCallStructName := fmt.Sprintf("%s%s%sFuncCall", iface.prefix, iface.titleName, method.Name)
	commentText := fmt.Sprintf(
		`%s describes the behavior when the %s method of the parent %s instance is invoked.`,
		mockFuncStructName,
		method.Name,
		mockStructName,
	)

	return generateStruct(mockFuncStructName, commentText, []jen.Code{
		compose(jen.Id("defaultHook"), method.signature),     // defaultHook <signature>
		compose(jen.Id("hooks").Index(), method.signature),   // hooks []<signature>
		jen.Id("history").Index().Id(mockFuncCallStructName), // history []<prefix>FuncCall
		jen.Id("mutex").Qual("sync", "Mutex"),                // mutex sync.Mutex
	})
}

func generateMockFuncCallStruct(iface *wrappedInterface, method *wrappedMethod, outputImportPath string) jen.Code {
	mockStructName := iface.mockStructName
	mockFuncCallStructName := fmt.Sprintf("%s%s%sFuncCall", iface.prefix, iface.titleName, method.Name)
	commentText := fmt.Sprintf(
		`%s is an object that describes an invocation of method %s on an instance of %s.`,
		mockFuncCallStructName,
		method.Name,
		mockStructName,
	)

	makeFields := func(prefix string, params []jen.Code, makeComment commentFactory) []jen.Code {
		fields := make([]jen.Code, 0, len(params))
		for i, param := range params {
			name := prefix + strconv.Itoa(i)
			field := jen.Id(name).Add(param)
			fields = append(fields, addComment(field, 2, makeComment(method, name, i)))
		}

		return fields
	}

	argFields := makeFields("Arg", method.dotlessParamTypes, argFieldComment)    // Arg<n> <ParamType #n>, ...
	resultFields := makeFields("Result", method.resultTypes, resultFieldComment) // Result<n> <ResultType #n>, ...
	return generateStruct(mockFuncCallStructName, commentText, append(argFields, resultFields...))
}

func generateStruct(name string, commentText string, structFields []jen.Code) jen.Code {
	typeDeclaration := jen.Type().Id(name).Struct(structFields...)
	return addComment(typeDeclaration, 1, commentText)
}

type commentFactory func(method *wrappedMethod, name string, i int) string

var (
	_ commentFactory = argFieldComment
	_ commentFactory = resultFieldComment
)

func argFieldComment(method *wrappedMethod, name string, i int) string {
	if i == len(method.dotlessParamTypes)-1 && method.Variadic {
		return fmt.Sprintf(
			`%s is a slice containing the values of the variadic arguments passed to this method invocation.`,
			name,
		)
	}

	return fmt.Sprintf(
		`%s is the value of the %s argument passed to this method invocation.`,
		name,
		humanize.Ordinal(i+1),
	)
}

func resultFieldComment(method *wrappedMethod, name string, i int) string {
	return fmt.Sprintf(
		`%s is the value of the %s result returned from this method invocation.`,
		name,
		humanize.Ordinal(i+1),
	)
}
