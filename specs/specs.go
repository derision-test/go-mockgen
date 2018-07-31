package specs

import (
	"go/types"
	"sort"
)

type (
	WrappedSpec struct {
		Spec       *InterfaceSpec
		ImportPath string
	}

	InterfaceSpec struct {
		Name      string
		TitleName string
		Methods   MethodSpecs
	}

	MethodSpec struct {
		Params   []types.Type
		Results  []types.Type
		Variadic bool
	}

	Specs          map[string]*WrappedSpec
	InterfaceSpecs map[string]*InterfaceSpec
	MethodSpecs    map[string]*MethodSpec
)

func (s Specs) Names() []string {
	names := []string{}
	for _, spec := range s {
		names = append(names, spec.Spec.Name)
	}

	sort.Strings(names)
	return names
}

func (s WrappedSpec) MethodNames() []string {
	names := []string{}
	for name := range s.Spec.Methods {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

func (s WrappedSpec) Method(name string) *MethodSpec {
	return s.Spec.Methods[name]
}
