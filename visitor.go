package main

import "go/ast"

type (
	visitor struct {
		typeSpecVisitor typeSpecVisitor
	}

	typeSpecVisitor interface {
		visitTypeSpec(*ast.TypeSpec)
	}
)

func walk(pkg *ast.Package, typeSpecVisitor typeSpecVisitor) {
	visitor := newVisitor(typeSpecVisitor)

	for _, file := range pkg.Files {
		visitor.walk(file)
	}
}

func newVisitor(typeSpecVisitor typeSpecVisitor) *visitor {
	return &visitor{typeSpecVisitor: typeSpecVisitor}
}

func (v *visitor) walk(expr *ast.File) {
	ast.Walk(v, expr)
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.GenDecl:
		for _, spec := range n.Specs {
			if typeSpec, ok := spec.(*ast.TypeSpec); ok {
				v.typeSpecVisitor.visitTypeSpec(typeSpec)
			}
		}
	}

	return v
}
