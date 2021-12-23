package generation

import "github.com/dave/jennifer/jen"

func compose(stmt *jen.Statement, tail ...jen.Code) *jen.Statement {
	head := *stmt
	for _, value := range tail {
		head = append(head, value)
	}

	return &head
}

func addComment(code *jen.Statement, level int, commentText string) *jen.Statement {
	if commentText == "" {
		return code
	}

	comment := generateComment(level, commentText)
	return compose(comment, code)
}

func selfAppend(sliceRef *jen.Statement, value jen.Code) jen.Code {
	return compose(sliceRef, jen.Op("=").Id("append").Call(sliceRef, value))
}
