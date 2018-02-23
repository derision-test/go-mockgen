package main

import (
	"bytes"
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

func generate(
	specs map[string]*interfaceSpec,
	packageName string,
	packageNames []string,
) error {
	file := jen.NewFile("test")

	for name, spec := range specs {
		generateMock(
			file,
			name,
			spec,
			packageName,
			packageNames,
		)
	}

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return err
	}

	fmt.Printf("%s\n", buffer.String())
	return nil
}

func generateMock(
	file *jen.File,
	interfaceName string,
	interfaceSpec *interfaceSpec,
	packageName string,
	packageNames []string,
) {
	var (
		fields   = []jen.Code{}
		defaults = []jen.Code{}
	)

	for funcName, method := range interfaceSpec.methods {
		var (
			params  = []jen.Code{}
			results = []jen.Code{}
			zeroes  = []jen.Code{}
		)

		for i, typ := range method.params {
			params = append(params, generateType(
				typ,
				packageName, packageNames,
				method.variadic && i == len(method.params)-1,
			))
		}

		for _, typ := range method.results {
			results = append(results, generateType(
				typ,
				packageName,
				packageNames,
				false,
			))

			zeroes = append(zeroes, zeroValue(
				typ,
				packageName,
				packageNames,
			))
		}

		field := jen.Id(fmt.Sprintf("%sFunc", funcName)).
			Func().
			Params(params...).
			Params(results...)

		fields = append(fields, field)

		var body jen.Code
		if len(zeroes) != 0 {
			body = jen.Return().List(zeroes...)
		}

		zero := jen.Func().Params(params...).Params(results...).Block(body)

		defaultVal := compose(jen.Id(fmt.Sprintf("%sFunc", funcName)).Op(":"), zero)

		defaults = append(defaults, defaultVal)
	}

	file.
		Type().
		Id(fmt.Sprintf("Mock%s", interfaceName)).
		Struct(fields...)

	file.Var().
		Id("_").
		Qual(packageName, interfaceName).
		Op("=").
		Id(fmt.Sprintf("NewMock%s", interfaceName)).
		Call()

	file.
		Func().
		Id(fmt.Sprintf("NewMock%s", interfaceName)).
		Params().
		Op("*").Id(fmt.Sprintf("Mock%s", interfaceName)).
		Block(jen.Return().
			Op("&").
			Id(fmt.Sprintf("Mock%s", interfaceName)).
			Values(defaults...))

	for funcName, method := range interfaceSpec.methods {
		var (
			names   = []jen.Code{}
			params  = []jen.Code{}
			results = []jen.Code{}

			body *jen.Statement
		)

		for i, typ := range method.params {
			variadic := method.variadic && i == len(method.params)-1
			name := jen.Id(fmt.Sprintf("v%d", i))
			if variadic {
				name = name.Op("...")
			}

			names = append(names, name)
			params = append(params, compose(jen.Id(fmt.Sprintf("v%d", i)), generateType(
				typ,
				packageName,
				packageNames,
				variadic,
			)))
		}

		for _, typ := range method.results {
			results = append(results, generateType(typ, packageName, packageNames, false))
		}

		if len(method.results) == 0 {
			body = jen.Id("m").
				Op(".").
				Id(fmt.Sprintf("%sFunc", funcName)).
				Call(names...)
		} else {
			body = jen.Return().
				Id("m").
				Op(".").
				Id(fmt.Sprintf("%sFunc", funcName)).
				Call(names...)
		}

		file.Func().
			Params(jen.Id("m").
				Op("*").
				Id(fmt.Sprintf("Mock%s", interfaceName))).
			Id(funcName).
			Params(params...).
			Params(results...).
			Block(body)
	}
}

func generateType(
	typ types.Type,
	packageName string,
	packageNames []string,
	variadic bool,
) *jen.Statement {
	recur := func(typ types.Type) *jen.Statement {
		return generateType(typ, packageName, packageNames, false)
	}

	switch t := typ.(type) {
	case *types.Basic:
		return jen.Id(typ.String())

	case *types.Chan:
		if t.Dir() == types.RecvOnly {
			return compose(jen.Op("<-").Chan(), recur(t.Elem()))
		}

		if t.Dir() == types.SendOnly {
			return compose(jen.Chan().Op("<-"), recur(t.Elem()))
		}

		return compose(jen.Chan(), recur(t.Elem()))

	case *types.Interface:
		methods := []jen.Code{}
		for i := 0; i < t.NumMethods(); i++ {
			method := t.Method(i)
			methods = append(methods, compose(jen.Id(method.Name()), recur(method.Type())))
		}

		return jen.Interface(methods...)

	case *types.Map:
		return compose(jen.Map(recur(t.Key())), recur(t.Elem()))

	case *types.Named:
		name := t.String()
		if stringInSlice(name, packageNames) {
			name = fmt.Sprintf("%s.%s", packageName, name)
		}

		if importPath, local := decomposePackage(name); importPath != "" {
			return jen.Qual(importPath, local)
		}

		return jen.Id(name)

	case *types.Pointer:
		return compose(jen.Op("*"), recur(t.Elem()))

	case *types.Slice:
		var prefix = jen.Index()
		if variadic {
			prefix = jen.Op("...")
		}

		return compose(prefix, recur(t.Elem()))

	case *types.Struct:
		fields := []jen.Code{}
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			fields = append(fields, compose(jen.Id(field.Name()), recur(field.Type())))
		}

		return jen.Struct(fields...)

	case *types.Signature:
		params := []jen.Code{}
		for i := 0; i < t.Params().Len(); i++ {
			typ := t.Params().At(i)
			params = append(params, compose(jen.Id(typ.Name()), recur(typ.Type())))
		}

		results := []jen.Code{}
		for i := 0; i < t.Results().Len(); i++ {
			results = append(results, recur(t.Results().At(i).Type()))
		}

		return jen.Func().Params(params...).Params(results...)

	default:
		panic(fmt.Sprintf("unsupported case: %#v\n", typ))
	}
}

func zeroValue(
	typ types.Type,
	packageName string,
	packageNames []string,
) *jen.Statement {
	switch t := typ.(type) {
	case *types.Basic:
		kind := t.Kind()

		if kind == types.Bool {
			return jen.False()
		} else if kind == types.String {
			return jen.Lit("")
		} else if isIntegerType(kind) {
			return jen.Lit(0)
		}

	case *types.Named:
		switch t.Underlying().(type) {
		case *types.Struct:
			path := t.Obj().Pkg().Path()
			if path == "" {
				path = packageName
			}

			return jen.Qual(path, t.Obj().Name()).Block()
		default:
			return zeroValue(t.Underlying(), packageName, packageNames)
		}

	case *types.Struct:
		return generateType(typ, packageName, packageNames, false).Block()
	}

	return jen.Nil()
}

func isIntegerType(kind types.BasicKind) bool {
	kinds := []types.BasicKind{
		types.Int,
		types.Int8,
		types.Int16,
		types.Int32,
		types.Int64,
		types.Uint,
		types.Uint8,
		types.Uint16,
		types.Uint32,
		types.Uint64,
		types.Uintptr,
		types.Float32,
		types.Float64,
		types.Byte,
		types.Rune,
		types.Complex64,
		types.Complex128,
	}

	for _, k := range kinds {
		if k == kind {
			return true
		}
	}

	return false
}

func compose(stmt1, stmt2 *jen.Statement) *jen.Statement {
	composed := append(*stmt1, stmt2)
	return &composed
}
