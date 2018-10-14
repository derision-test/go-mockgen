package main

import (
	"fmt"
	gotypes "go/types"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/efritz/go-genlib/command"
	"github.com/efritz/go-genlib/generation"
	"github.com/efritz/go-genlib/types"
)

const (
	Name        = "go-mockgen"
	Description = "go-mockgen generates mock implementations from interface definitions."
	Version     = "0.1.0"
)

func main() {
	if err := command.Run(Name, Description, Version, types.GetInterface, generate); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}

func generate(ifaces []*types.Interface, opts *command.Options) error {
	return generation.Generate("github.com/efritz/go-mockgen", ifaces, opts, generateFilename, generateInterface)
}

func generateFilename(name string) string {
	return fmt.Sprintf("%s_mock.go", name)
}

func generateInterface(file *jen.File, iface *types.Interface, prefix string) {
	var (
		titleName         = title(iface.Name)
		mockStructName    = fmt.Sprintf("Mock%s%s", prefix, titleName)
		constructorName   = fmt.Sprintf("New%s", mockStructName)
		structFields      = []jen.Code{}
		constructorFields = []jen.Code{}
		methodParts       = map[string]methodParts{}
	)

	funcs := []func(string) jen.Code{
		func(name string) jen.Code { return methodParts[name].paramSetDef },
		func(name string) jen.Code { return methodParts[name].overrideMethodDef },
		func(name string) jen.Code { return methodParts[name].callCountMethodDef },
		func(name string) jen.Code { return methodParts[name].callParamMethodDef },
	}

	for _, method := range iface.Methods {
		parts := generateMethodParts(
			iface,
			method,
			prefix,
			titleName,
			mockStructName,
		)

		constructorFields = append(constructorFields, parts.constructorField)
		structFields = append(structFields, parts.hookFunctionField, parts.callParamsMethodField)
		methodParts[method.Name] = parts
	}

	structDef := jen.Type().Id(mockStructName).Struct(append(
		structFields,
		jen.Id("mutex").Qual("sync", "RWMutex"),
	)...)

	constructorDef := generation.GenerateFunction(
		constructorName,
		nil,
		[]jen.Code{jen.Op("*").Id(mockStructName)},
		jen.Return().Op("&").Id(mockStructName).Values(constructorFields...),
	)

	file.Add(structDef)
	file.Add(constructorDef)

	for _, name := range iface.MethodNames() {
		for _, f := range funcs {
			file.Add(f(name))
			file.Line()
		}
	}
}

type methodParts struct {
	constructorField      jen.Code
	hookFunctionField     jen.Code
	callParamsMethodField jen.Code
	paramSetDef           jen.Code
	overrideMethodDef     jen.Code
	callCountMethodDef    jen.Code
	callParamMethodDef    jen.Code
}

func generateMethodParts(
	iface *types.Interface,
	method *types.Method,
	prefix string,
	titleName string,
	mockStructName string,
) (parts methodParts) {
	var (
		hookFunctionName     = fmt.Sprintf("%sFunc", method.Name)
		paramSetStructName   = fmt.Sprintf("%s%s%sParamSet", prefix, titleName, method.Name)
		callCountMethodName  = fmt.Sprintf("%sFuncCallCount", method.Name)
		callParamsMethodName = fmt.Sprintf("%sFuncCallParams", method.Name)
		callHistoryFieldName = fmt.Sprintf("_%sFuncCallHistory", method.Name)
		paramTypes           = generation.GenerateParamTypes(method, iface.ImportPath, false)
		resultTypes          = generation.GenerateResultTypes(method, iface.ImportPath)
		paramSetStructFields = makeParamSetFields(generation.GenerateParamTypes(method, iface.ImportPath, true))
		lock                 = jen.Id("m").Dot("mutex").Dot("RLock").Call()
		unlock               = jen.Id("m").Dot("mutex").Dot("RUnlock").Call()
		deferUnlock          = generation.Compose(jen.Defer(), unlock)
		callHistoryFieldRef  = jen.Id("m").Dot(callHistoryFieldName)
		zeroFunction         = makeZeroFunction(iface.ImportPath, method.Results, paramTypes, resultTypes)
		appendParamSet       = selfAppend(callHistoryFieldRef, makeParamSet(paramSetStructName, len(method.Params)))
	)

	parts.constructorField = generation.Compose(
		jen.Id(hookFunctionName).Op(":"), // <Method>Func:
		zeroFunction,                     // func(t1, ..., tn) { return z1, ..., zn }
	)

	parts.hookFunctionField = jen.
		Id(hookFunctionName).  // <Method>Func
		Func().                // func
		Params(paramTypes...). // (...)
		Params(resultTypes...) // (...)

	parts.callParamsMethodField = jen.
		Id(callHistoryFieldName). // _<Method>FuncCallHistory
		Index().                  // []
		Id(paramSetStructName)    // <Method>ParamSet

	parts.paramSetDef = jen.
		Type().                         // type
		Id(paramSetStructName).         // <Prefix><Struct><Method>ParamSet
		Struct(paramSetStructFields...) // struct { ... }

	parts.overrideMethodDef = generation.GenerateOverride(
		"m",                                    // func (m
		mockStructName,                         // m*Mock<Struct>)
		iface.ImportPath,                       //
		method,                                 // <Method>(v0, v1, ...) (...) {
		lock,                                   // m.mutex.Lock()
		appendParamSet,                         // m.<CallHistory> = append(m.<CallHistory>, ParamSet{...})
		unlock,                                 // m.mutex.Unlock()
		generation.GenerateSuperCall(method),   // r0, ..., rm := m.<Method>(v0, ..., vn)
		generation.GenerateSuperReturn(method), // return r0, ..., rm }
	)

	parts.callCountMethodDef = generation.GenerateMethod(
		"m",                                      // func (m
		mockStructName,                           // *Mock<Struct>)
		callCountMethodName,                      // <Method>FuncCallCount
		nil,                                      // ()
		[]jen.Code{jen.Int()},                    // int
		lock,                                     // m.mutex.Lock()
		deferUnlock,                              // defer m.mutex.Unlock()
		jen.Return(jen.Len(callHistoryFieldRef)), // return len(m._<Method>FuncCallHistory) }
	)

	parts.callParamMethodDef = generation.GenerateMethod(
		"m",                                   // func (m
		mockStructName,                        // *Mock<Struct>)
		callParamsMethodName,                  // <Method>FuncCallParams
		nil,                                   // ()
		[]jen.Code{index(paramSetStructName)}, // []<Struct>ParamSet {
		lock,                                  // m.mutex.Lock()
		deferUnlock,                           // defer m.mutex.Unlock()
		jen.Return(callHistoryFieldRef),       // return m._<Method>FuncCallHistory }
	)

	return
}

//
// Helpers

func title(s string) string {
	if s == "" {
		return s
	}

	return strings.ToUpper(string(s[0])) + s[1:]
}

func index(name string) jen.Code {
	return jen.Index().Id(name)
}

func selfAppend(sliceRef *jen.Statement, value jen.Code) jen.Code {
	return generation.Compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}

func makeZeroFunction(importPath string, results []gotypes.Type, paramTypes, resultTypes []jen.Code) jen.Code {
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

func makeParamSetFields(paramTypesNoDots []jen.Code) []jen.Code {
	paramSetStructFields := []jen.Code{}
	for i, param := range paramTypesNoDots {
		paramSetStructFields = append(paramSetStructFields, jen.Id(fmt.Sprintf("Arg%d", i)).Add(param))
	}

	return paramSetStructFields
}

func makeParamSet(paramSetStructName string, paramCount int) jen.Code {
	names := []jen.Code{}
	for i := 0; i < paramCount; i++ {
		names = append(names, jen.Id(fmt.Sprintf("v%d", i)))
	}

	return jen.Id(paramSetStructName).Values(names...)
}
