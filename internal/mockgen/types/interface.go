package types

import (
	"go/types"
	"sort"
)

type Interface struct {
	Name       string
	ImportPath string
	Methods    []*Method
}

func newInterfaceFromTypeSpec(name, importPath string, typeSpec *types.Interface) *Interface {
	methodMap := make(map[string]*Method, typeSpec.NumMethods())
	for i := 0; i < typeSpec.NumMethods(); i++ {
		method := typeSpec.Method(i)
		name := method.Name()
		methodMap[name] = newMethodFromSignature(name, method.Type().(*types.Signature))
	}

	methodNames := make([]string, 0, len(methodMap))
	for k := range methodMap {
		methodNames = append(methodNames, k)
	}
	sort.Strings(methodNames)

	methods := make([]*Method, 0, len(methodNames))
	for _, name := range methodNames {
		methods = append(methods, methodMap[name])
	}

	return &Interface{
		Name:       name,
		ImportPath: importPath,
		Methods:    methods,
	}
}
