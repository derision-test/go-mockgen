package generation

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

func stripVendor(path string) string {
	parts := strings.Split(path, "/vendor/")
	return parts[len(parts)-1]
}

func compose(stmt1 *jen.Statement, stmt2 jen.Code) *jen.Statement {
	composed := append(*stmt1, stmt2)
	return &composed
}
