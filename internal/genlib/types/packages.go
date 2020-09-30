package types

import (
	"fmt"
	"sort"
)

type (
	Packages struct {
		packages map[string]*Package
	}

	TypeGetter func(pkgs *Packages, name string) (*Interface, error)
)

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

func (p *Packages) GetType(name string) (*Interface, error)      { return p.getType(name, aType) }
func (p *Packages) GetStruct(name string) (*Interface, error)    { return p.getType(name, sType) }
func (p *Packages) GetInterface(name string) (*Interface, error) { return p.getType(name, iType) }

func (p *Packages) getType(name string, matcher func(InterfaceType) bool) (*Interface, error) {
	candidates := []*Interface{}
	for _, pkg := range p.packages {
		if t, ok := pkg.Types[name]; ok {
			if matcher(t.Type) {
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

//
// Helpers

func aType(ifaceType InterfaceType) bool                           { return true }
func sType(ifaceType InterfaceType) bool                           { return ifaceType == InterfaceTypeStruct }
func iType(ifaceType InterfaceType) bool                           { return ifaceType == InterfaceTypeInterface }
func GetType(pkgs *Packages, name string) (*Interface, error)      { return pkgs.GetType(name) }
func GetStruct(pkgs *Packages, name string) (*Interface, error)    { return pkgs.GetStruct(name) }
func GetInterface(pkgs *Packages, name string) (*Interface, error) { return pkgs.GetInterface(name) }
