package generation

import (
	gotypes "go/types"
	"strings"

	"github.com/derision-test/go-mockgen/internal/mockgen/types"
)

func getType(kind gotypes.BasicKind) gotypes.Type {
	return gotypes.Typ[kind].Underlying()
}

func makeBareInterface(methods ...*types.Method) *types.Interface {
	return &types.Interface{
		Name:       TestTitleName,
		ImportPath: TestImportPath,
		Type:       types.InterfaceTypeInterface,
		Methods:    methods,
	}
}

func makeInterface(methods ...*types.Method) (*wrappedInterface, string) {
	return wrapInterface(makeBareInterface(methods...), TestPrefix, TestTitleName, TestMockStructName, ""), ""
}

func makeMethod(methods ...*types.Method) (*wrappedInterface, *wrappedMethod, string) {
	wrapped, _ := makeInterface(methods...)
	return wrapped, wrapped.wrappedMethods[0], ""
}

func strip(block string) string {
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "\t\t") {
			lines[i] = line[2:]
		}
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
