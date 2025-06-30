package generation

import (
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

func generateComment(level int, commentText string) *jen.Statement {
	allowance := maxAllowance - indent*level - 3
	if allowance < minAllowance {
		allowance = minAllowance
	}

	wrapped := wordwrap.WrapString(commentText, uint(allowance))
	lines := strings.Split(wrapped, "\n")
	commentBlock := jen.Comment(lines[0]).Line()

	for i := 1; i < len(lines); i++ {
		commentBlock = commentBlock.Comment(lines[i]).Line()
	}

	return commentBlock
}
