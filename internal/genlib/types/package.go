package types

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
