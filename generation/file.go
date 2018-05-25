package generation

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/efritz/go-mockgen/specs"
)

func generateContent(allSpecs specs.Specs, pkgName, prefix string) (string, error) {
	file := jen.NewFile(pkgName)
	file.HeaderComment("DO NOT EDIT")
	file.HeaderComment("Code generated automatically by github.com/efritz/go-mockgen")
	file.HeaderComment(fmt.Sprintf("$ %s", strings.Join(os.Args, " ")))

	for _, name := range allSpecs.Names() {
		generator := newInterfaceGenerator(
			file,
			name,
			prefix,
			allSpecs[name],
		)

		generator.generate()
	}

	buffer := &bytes.Buffer{}
	if err := file.Render(buffer); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
