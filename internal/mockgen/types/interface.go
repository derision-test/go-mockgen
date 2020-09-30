package types

import (
	"go/types"
	"sort"
)

type (
	Interface struct {
		Name       string
		ImportPath string
		Type       InterfaceType
		Methods    []*Method
	}

	InterfaceType int
)

const (
	InterfaceTypeStruct InterfaceType = iota
	InterfaceTypeInterface
)

func (i *Interface) MethodNames() []string {
	names := []string{}
	for _, method := range i.Methods {
		names = append(names, method.Name)
	}

	return names
}

func (i *Interface) Method(name string) *Method {
	for _, method := range i.Methods {
		if method.Name == name {
			return method
		}
	}

	return nil
}

func DeconstructStruct(name, importPath string, typeSpec *types.Struct) *Interface {
	methodMap := map[string]*Method{}
	methodNames := []string{}

	for i := 0; i < typeSpec.NumFields(); i++ {
		field := typeSpec.Field(i)
		name := field.Name()

		if signature, ok := field.Type().(*types.Signature); ok {
			methodMap[name] = DeconstructMethod(name, signature)
			methodNames = append(methodNames, name)
		}
	}

	sort.Strings(methodNames)

	methods := []*Method{}
	for _, name := range methodNames {
		methods = append(methods, methodMap[name])
	}

	return &Interface{
		Name:       name,
		ImportPath: importPath,
		Type:       InterfaceTypeStruct,
		Methods:    methods,
	}
}

func DeconstructInterface(name, importPath string, typeSpec *types.Interface) *Interface {
	methodMap := map[string]*Method{}
	methodNames := []string{}

	for i := 0; i < typeSpec.NumMethods(); i++ {
		method := typeSpec.Method(i)
		name := method.Name()
		signature := method.Type().(*types.Signature)

		methodMap[name] = DeconstructMethod(name, signature)
		methodNames = append(methodNames, name)
	}

	sort.Strings(methodNames)

	methods := []*Method{}
	for _, name := range methodNames {
		methods = append(methods, methodMap[name])
	}

	return &Interface{
		Name:       name,
		ImportPath: importPath,
		Type:       InterfaceTypeInterface,
		Methods:    methods,
	}
}
