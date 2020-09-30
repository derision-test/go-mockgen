package types

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/derision-test/go-mockgen/internal/mockgen/paths"
	"golang.org/x/tools/go/packages"
)

func Extract(importPaths []string, targetNames []string) ([]*Interface, error) {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory (%s)", err.Error())
	}

	allPkgs := map[string]*Package{}
	for _, importPath := range importPaths {
		path, dir := paths.ResolveImportPath(workingDirectory, importPath)
		log.Printf("parsing package '%s'\n", paths.GetRelativePath(dir))

		pkg, err := getPackage(workingDirectory, importPath, path)
		if err != nil {
			return nil, err
		}

		allPkgs[path] = pkg
	}

	pkgs := NewPackages(allPkgs)

	ifaces := []*Interface{}
	for _, name := range pkgs.GetNames() {
		iface, err := getInterface(pkgs, name, targetNames)
		if err != nil {
			return nil, err
		}

		if iface != nil {
			ifaces = append(ifaces, iface)
		}
	}

	return ifaces, nil
}

func getPackage(workingDirectory string, importPath, path string) (*Package, error) {
	pkgs, err := packages.Load(&packages.Config{Mode: packages.LoadSyntax}, importPath)
	if err != nil {
		return nil, fmt.Errorf("could not load package %s (%s)", importPath, err.Error())
	}

	for _, err := range pkgs[0].Errors {
		return nil, fmt.Errorf("malformed package %s (%s)", importPath, err.Msg)
	}

	visitor := newVisitor(path, pkgs[0].Types)
	for _, file := range pkgs[0].Syntax {
		ast.Walk(visitor, file)
	}

	return NewPackage(path, visitor.types), nil
}

func getInterface(pkgs *Packages, name string, targetNames []string) (*Interface, error) {
	if !shouldInclude(name, targetNames) {
		return nil, nil
	}

	iface, err := pkgs.GetInterface(name)
	if err != nil || iface == nil {
		return nil, err
	}
	if iface == nil {
		return nil, nil
	}

	for _, method := range iface.Methods {
		if !unicode.IsUpper([]rune(method.Name)[0]) {
			return nil, fmt.Errorf(
				"type '%s' has unexported an method '%s'",
				name,
				method.Name,
			)
		}
	}

	return iface, nil
}

func shouldInclude(name string, targetNames []string) bool {
	for _, v := range targetNames {
		if strings.ToLower(v) == strings.ToLower(name) {
			return true
		}
	}

	return len(targetNames) == 0
}

type visitor struct {
	importPath string
	pkgType    *types.Package
	types      map[string]*Interface
}

func newVisitor(importPath string, pkgType *types.Package) *visitor {
	return &visitor{
		importPath: importPath,
		pkgType:    pkgType,
		types:      map[string]*Interface{},
	}
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				name := typeSpec.Name.Name
				_, obj := v.pkgType.Scope().Innermost(typeSpec.Pos()).LookupParent(name, 0)

				switch t := obj.Type().Underlying().(type) {
				case *types.Interface:
					v.types[name] = deconstructInterface(name, v.importPath, t)
				}
			}
		}
	}

	return v
}

func deconstructInterface(name, importPath string, typeSpec *types.Interface) *Interface {
	methodMap := map[string]*Method{}
	for i := 0; i < typeSpec.NumMethods(); i++ {
		method := typeSpec.Method(i)
		name := method.Name()
		methodMap[name] = deconstructMethod(name, method.Type().(*types.Signature))
	}

	methodNames := []string{}
	for k := range methodMap {
		methodNames = append(methodNames, k)
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

func deconstructMethod(name string, signature *types.Signature) *Method {
	ps := signature.Params()
	params := []types.Type{}
	for i := 0; i < ps.Len(); i++ {
		params = append(params, ps.At(i).Type())
	}

	rs := signature.Results()
	results := []types.Type{}
	for i := 0; i < rs.Len(); i++ {
		results = append(results, rs.At(i).Type())
	}

	return &Method{
		Name:     name,
		Params:   params,
		Results:  results,
		Variadic: signature.Variadic(),
	}
}
