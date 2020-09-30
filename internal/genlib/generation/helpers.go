package generation

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/mitchellh/go-wordwrap"
)

var (
	maxAllowance = 80
	minAllowance = maxAllowance - indent*maxLevels
	indent       = 4
	maxLevels    = 3
)

func Compose(stmt1 *jen.Statement, stmt2 jen.Code) *jen.Statement {
	composed := append(*stmt1, stmt2)
	return &composed
}

func GenerateComment(level int, format string, args ...interface{}) *jen.Statement {
	allowance := maxAllowance - indent*level - 3
	if allowance < minAllowance {
		allowance = minAllowance
	}

	var (
		commentText  = fmt.Sprintf(format, args...)
		wrapped      = wordwrap.WrapString(commentText, uint(allowance))
		lines        = strings.Split(wrapped, "\n")
		commentBlock = jen.Comment(lines[0]).Line()
	)

	for i := 1; i < len(lines); i++ {
		commentBlock = commentBlock.Comment(lines[i]).Line()
	}

	return commentBlock
}
