package types

import (
	"fmt"
	"go/types"
	"sort"
)

type Interface struct {
	Name       string
	ImportPath string
	Type       InterfaceType
	Methods    []*Method
}

type InterfaceType int

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

type Method struct {
	Name     string
	Params   []types.Type
	Results  []types.Type
	Variadic bool
}

type Package struct {
	Name  string
	Types map[string]*Interface
}

func NewPackage(name string, types map[string]*Interface) *Package {
	return &Package{
		Name:  name,
		Types: types,
	}
}

type Packages struct {
	packages map[string]*Package
}

func NewPackages(packages map[string]*Package) *Packages {
	return &Packages{
		packages: packages,
	}
}

func (p *Packages) GetNames() []string {
	nameMap := map[string]struct{}{}
	for _, pkg := range p.packages {
		for name, _ := range pkg.Types {
			nameMap[name] = struct{}{}
		}
	}

	names := []string{}
	for name := range nameMap {
		names = append(names, name)
	}

	sort.Strings(names)
	return names
}

func (p *Packages) GetInterface(name string) (*Interface, error) {
	candidates := []*Interface{}
	for _, pkg := range p.packages {
		if t, ok := pkg.Types[name]; ok {
			if t.Type == InterfaceTypeInterface {
				candidates = append(candidates, t)
			}
		}
	}

	if len(candidates) > 1 {
		return nil, fmt.Errorf("type '%s' is multiply-defined in supplied import paths", name)
	}

	if len(candidates) == 1 {
		return candidates[0], nil
	}

	return nil, nil
}
