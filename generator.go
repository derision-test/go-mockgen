package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
)

const (
	mockFormat          = "Mock%s%s"
	constructorFormat   = "NewMock%s%s"
	innerMethodFormat   = "%sFunc"
	parameterNameFormat = "v%d"
)

func generate(specs map[string]*wrappedSpec, pkgName, prefix, dirname, filename string, force bool) error {
	if dirname != "" && filename == "" {
		return generateMultipleFiles(specs, pkgName, prefix, dirname, force)
	}

	return generateOneFile(specs, pkgName, prefix, path.Join(dirname, filename), force)
}

func generateMultipleFiles(specs map[string]*wrappedSpec, pkgName, prefix, dirname string, force bool) error {
	if !force {
		paths := []string{}
		for name := range specs {
			paths = append(paths, getFilename(dirname, name, prefix))
		}

		conflict, err := anyPathExists(paths)
		if err != nil {
			return err
		}

		if conflict != "" {
			return fmt.Errorf("filename %s already exists", conflict)
		}
	}

	for name, spec := range specs {
		content, err := generateContent(map[string]*wrappedSpec{name: spec}, pkgName, prefix)
		if err != nil {
			return err
		}

		if err := writeFile(getFilename(dirname, name, prefix), content); err != nil {
			return err
		}
	}

	return nil
}

func generateOneFile(specs map[string]*wrappedSpec, pkgName, prefix, filename string, force bool) error {
	content, err := generateContent(specs, pkgName, prefix)
	if err != nil {
		return err
	}

	if filename != "" && !force {
		exists, err := pathExists(filename)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("filename %s already exists", filename)
		}

		return writeFile(filename, content)
	}

	fmt.Printf("%s\n", content)
	return nil
}

func generateContent(specs map[string]*wrappedSpec, pkgName, prefix string) (string, error) {
	file := jen.NewFile(pkgName)
	for _, name := range getNames(specs) {
		generateFile(file, name, prefix, specs[name])
	}

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func getFilename(dirname, interfaceName, prefix string) string {
	filename := fmt.Sprintf("%s_mock.go", interfaceName)
	if prefix != "" {
		filename = fmt.Sprintf("%s_%s", prefix, filename)
	}

	return path.Join(dirname, strings.Replace(strings.ToLower(filename), "-", "_", -1))
}

func pathExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func anyPathExists(paths []string) (string, error) {
	for _, path := range paths {
		exists, err := pathExists(path)
		if err != nil {
			return "", err
		}

		if exists {
			return path, nil
		}
	}

	return "", nil
}

//
// Code Generation

func generateFile(file *jen.File, name, prefix string, spec *wrappedSpec) {
	generateInterfaceDefinition(file, name, prefix, spec)
	generateTypeTest(file, name, prefix, spec)
	generateConstructor(file, name, prefix, spec)
	generateMethodImplementations(file, name, prefix, spec)
}

// generateInterfaceDefinition
//
// type Mock{{Interface}} struct {
//     {{Method}}Name func({{params...}}) {{results...}}
// }

func generateInterfaceDefinition(file *jen.File, interfaceName, prefix string, interfaceSpec *wrappedSpec) {
	fields := []jen.Code{}
	for methodName, method := range interfaceSpec.spec.methods {
		fields = append(fields, generateMethodField(
			methodName,
			method,
			interfaceSpec.importPath,
		))
	}

	file.Type().Id(fmt.Sprintf(mockFormat, prefix, interfaceName)).Struct(fields...)
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

func generateTypeTest(file *jen.File, interfaceName, prefix string, interfaceSpec *wrappedSpec) {
	constructorName := fmt.Sprintf(constructorFormat, prefix, interfaceName)

	file.Var().
		Id("_").
		Qual(stripVendor(interfaceSpec.importPath), interfaceName).
		Op("=").
		Id(constructorName).
		Call()
}

func stripVendor(path string) string {
	parts := strings.Split(path, "/vendor/")
	return parts[len(parts)-1]
}

// generateConstructor
//
// func NewMock{{Interface}} *Mock{{Interface}} {
//     return &Mock{{Interface}}{
//         {{Method}}Func func({{params...}}) {{results...}} { return {{result-zero-values...}} }
//     }
// }

func generateConstructor(file *jen.File, interfaceName, prefix string, interfaceSpec *wrappedSpec) {
	structName := fmt.Sprintf(mockFormat, prefix, interfaceName)
	constructorName := fmt.Sprintf(constructorFormat, prefix, interfaceName)

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

func generateMethodImplementations(file *jen.File, interfaceName, prefix string, interfaceSpec *wrappedSpec) {
	for methodName, method := range interfaceSpec.spec.methods {
		generateMethodImplementation(file, interfaceName, prefix, interfaceSpec.importPath, methodName, method)
	}
}

func generateMethodImplementation(file *jen.File, interfaceName, prefix string, importPath, methodName string, method *methodSpec) {
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
		Id(fmt.Sprintf(mockFormat, prefix, interfaceName))

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
