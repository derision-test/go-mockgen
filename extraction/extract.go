package extraction

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
	"unicode"

	"github.com/efritz/go-mockgen/specs"
)

func Extract(wd string, importPaths []string, interfaces []string) (specs.Specs, error) {
	var (
		parser   = newPathParser(wd)
		allSpecs = specs.Specs{}
	)

	for _, path := range importPaths {
		pkg, pkgType, err := parser.parse(path)
		if err != nil {
			return nil, err
		}

		interfaceSpecs, err := getInterfaceSpecs(pkg, pkgType)
		if err != nil {
			return nil, err
		}

		if len(interfaceSpecs) == 0 {
			return nil, fmt.Errorf("no interfaces found in path %s", path)
		}

		for name, spec := range interfaceSpecs {
			if !shouldInclude(name, interfaces) {
				continue
			}

			for method := range spec.Methods {
				if !unicode.IsUpper([]rune(method)[0]) {
					return nil, fmt.Errorf(
						"interface %s contains unexported method %s",
						name,
						method,
					)
				}
			}

			if _, ok := allSpecs[name]; ok {
				return nil, fmt.Errorf("ambiguous interface %s in supplied import paths", name)
			}

			allSpecs[name] = &specs.WrappedSpec{
				Spec:       spec,
				ImportPath: path,
			}
		}
	}

	return allSpecs, nil
}

func getInterfaceSpecs(pkg *ast.Package, pkgType *types.Package) (specs.InterfaceSpecs, error) {
	visitor := newVisitor(pkgType)
	for _, file := range pkg.Files {
		ast.Walk(visitor, file)
	}

	return visitor.specs, visitor.err
}

func shouldInclude(name string, interfaces []string) bool {
	for _, v := range interfaces {
		if strings.ToLower(v) == strings.ToLower(name) {
			return true
		}
	}

	return len(interfaces) == 0
}
