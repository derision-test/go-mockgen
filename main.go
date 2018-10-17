package main

import (
	"fmt"
	gotypes "go/types"
	"log"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/efritz/go-genlib/command"
	"github.com/efritz/go-genlib/generation"
	"github.com/efritz/go-genlib/types"
)

const (
	Name        = "go-mockgen"
	PackageName = "github.com/efritz/go-mockgen"
	Description = "go-mockgen generates mock implementations from interface definitions."
	Version     = "0.1.0"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("go-mockgen: ")
}

func main() {
	if err := command.Run(Name, Description, Version, types.GetInterface, generate); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func generate(ifaces []*types.Interface, opts *command.Options) error {
	return generation.Generate(PackageName, ifaces, opts, generateFilename, generateInterface)
}

func generateFilename(name string) string {
	return fmt.Sprintf("%s_mock.go", name)
}

func generateInterface(file *jen.File, iface *types.Interface, prefix string) {
	var (
		titleName      = title(iface.Name)
		mockStructName = fmt.Sprintf("Mock%s%s", prefix, titleName)
	)

	file.Add(generateStruct(iface, prefix, titleName, mockStructName))
	file.Add(generateConstructor(iface, mockStructName))

	methodGenerators := []func(*types.Interface, *types.Method, string, string, string) jen.Code{
		generateParamSetStruct,
		generateOverrideMethod,
		generateCallCountMethod,
		generateCallParamsMethod,
	}

	for _, method := range iface.Methods {
		for _, generator := range methodGenerators {
			file.Add(generator(
				iface,
				method,
				prefix,
				titleName,
				mockStructName,
			))

			file.Line()
		}
	}
}

func generateStruct(iface *types.Interface, prefix, titleName, mockStructName string) jen.Code {
	structFields := []jen.Code{}
	for _, method := range iface.Methods {
		hookFuncField := jen.
			Id(fmt.Sprintf("%sFunc", method.Name)).
			Func().
			Params(generation.GenerateParamTypes(method, iface.ImportPath, false)...).
			Params(generation.GenerateResultTypes(method, iface.ImportPath)...)

		callHistoryField := jen.
			Id(fmt.Sprintf("%sFuncCallHistory", method.Name)).
			Index().
			Id(fmt.Sprintf("%s%s%sParamSet", prefix, titleName, method.Name))

		structFields = append(structFields, hookFuncField)
		structFields = append(structFields, callHistoryField)
	}

	structFields = append(structFields, jen.Id("mutex").Qual("sync", "RWMutex"))

	return jen.
		Type().
		Id(mockStructName).
		Struct(structFields...)
}

func generateConstructor(iface *types.Interface, mockStructName string) jen.Code {
	constructorFields := []jen.Code{}
	for _, method := range iface.Methods {
		fieldTypes := generation.GenerateParamTypes(
			method,
			iface.ImportPath,
			false,
		)

		zeroFunction := generateZeroFunction(
			iface.ImportPath,
			method.Results,
			fieldTypes,
			generation.GenerateResultTypes(method, iface.ImportPath),
		)

		constructorFields = append(constructorFields, generation.Compose(
			jen.Id(fmt.Sprintf("%sFunc", method.Name)).Op(":"),
			zeroFunction,
		))
	}

	return generation.GenerateFunction(
		fmt.Sprintf("New%s", mockStructName),
		nil,
		[]jen.Code{jen.Op("*").Id(mockStructName)},
		jen.Return().Op("&").Id(mockStructName).Values(constructorFields...),
	)
}

func generateParamSetStruct(
	iface *types.Interface,
	method *types.Method,
	prefix string,
	titleName string,
	mockStructName string,
) jen.Code {
	fieldTypes := generation.GenerateParamTypes(
		method,
		iface.ImportPath,
		true,
	)

	return jen.
		Type().
		Id(fmt.Sprintf("%s%s%sParamSet", prefix, titleName, method.Name)).
		Struct(generateParamSetStructFields(fieldTypes)...)
}

func generateOverrideMethod(
	iface *types.Interface,
	method *types.Method,
	prefix string,
	titleName string,
	mockStructName string,
) jen.Code {
	callTarget := jen.
		Id("m").
		Dot(fmt.Sprintf("%sFunc", method.Name))

	historyInstance := generateParamSetInstance(
		fmt.Sprintf("%s%s%sParamSet", prefix, titleName, method.Name),
		len(method.Params),
	)

	appendHistory := selfAppend(
		jen.Id("m").Dot(fmt.Sprintf("%sFuncCallHistory", method.Name)),
		historyInstance,
	)

	return generation.GenerateOverride(
		"m",
		mockStructName,
		iface.ImportPath,
		method,
		jen.Id("m").Dot("mutex").Dot("RLock").Call(),
		appendHistory,
		jen.Id("m").Dot("mutex").Dot("RUnlock").Call(),
		generation.GenerateDecoratedCall(method, callTarget),
		generation.GenerateDecoratedReturn(method),
	)
}

func generateCallCountMethod(
	iface *types.Interface,
	method *types.Method,
	prefix string,
	titleName string,
	mockStructName string,
) jen.Code {
	return generation.GenerateMethod(
		"m",
		mockStructName,
		fmt.Sprintf("%sFuncCallCount", method.Name),
		nil,
		[]jen.Code{jen.Int()},
		jen.Id("m").Dot("mutex").Dot("RLock").Call(),
		generation.Compose(jen.Defer(), jen.Id("m").Dot("mutex").Dot("RUnlock").Call()),
		jen.Return(jen.Len(jen.Id("m").Dot(fmt.Sprintf("%sFuncCallHistory", method.Name)))),
	)
}

func generateCallParamsMethod(
	iface *types.Interface,
	method *types.Method,
	prefix string,
	titleName string,
	mockStructName string,
) jen.Code {
	return generation.GenerateMethod(
		"m",
		mockStructName,
		fmt.Sprintf("%sFuncCallParams", method.Name),
		nil,
		[]jen.Code{index(fmt.Sprintf("%s%s%sParamSet", prefix, titleName, method.Name))},
		jen.Id("m").Dot("mutex").Dot("RLock").Call(),
		generation.Compose(jen.Defer(), jen.Id("m").Dot("mutex").Dot("RUnlock").Call()),
		jen.Return(jen.Id("m").Dot(fmt.Sprintf("%sFuncCallHistory", method.Name))),
	)
}

func generateZeroFunction(
	importPath string,
	results []gotypes.Type,
	paramTypes []jen.Code,
	resultTypes []jen.Code,
) jen.Code {
	zeroes := []jen.Code{}
	for _, typ := range results {
		zeroes = append(zeroes, generation.GenerateZeroValue(
			typ,
			importPath,
		))
	}

	return generation.GenerateFunction(
		"",
		paramTypes,
		resultTypes,
		jen.Return().List(zeroes...),
	)
}

func generateParamSetStructFields(paramTypesNoDots []jen.Code) []jen.Code {
	fields := []jen.Code{}
	for i, param := range paramTypesNoDots {
		fields = append(fields, jen.Id(fmt.Sprintf("Arg%d", i)).Add(param))
	}

	return fields
}

func generateParamSetInstance(paramSetStructName string, paramCount int) jen.Code {
	names := []jen.Code{}
	for i := 0; i < paramCount; i++ {
		names = append(names, jen.Id(fmt.Sprintf("v%d", i)))
	}

	return jen.Id(paramSetStructName).Values(names...)
}

func index(name string) jen.Code {
	return jen.Index().Id(name)
}

func selfAppend(sliceRef *jen.Statement, value jen.Code) jen.Code {
	return generation.Compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}

func title(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}
