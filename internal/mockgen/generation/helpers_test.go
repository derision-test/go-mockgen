package generation

import (
	gotypes "go/types"
	"strings"

	"github.com/derision-test/go-mockgen/internal/mockgen/types"
)

const (
	TestPrefix         = "Test"
	TestTitleName      = "Client"
	TestMockStructName = "MockTestClient"
	TestImportPath     = "github.com/derision-test/go-mockgen/test"
)

var (
	boolType        = getType(gotypes.Bool)
	stringType      = getType(gotypes.String)
	stringSliceType = gotypes.NewSlice(getType(gotypes.String))

	TestMethodStatus = &types.Method{
		Name:    "Status",
		Params:  []gotypes.Type{},
		Results: []gotypes.Type{stringType, boolType},
	}

	TestMethodDo = &types.Method{
		Name:    "Do",
		Params:  []gotypes.Type{stringType},
		Results: []gotypes.Type{boolType},
	}

	TestMethodDof = &types.Method{
		Name:     "Dof",
		Params:   []gotypes.Type{stringType, stringSliceType},
		Results:  []gotypes.Type{boolType},
		Variadic: true,
	}
)

func getType(kind gotypes.BasicKind) gotypes.Type {
	return gotypes.Typ[kind].Underlying()
}

func makeBareInterface(methods ...*types.Method) *types.Interface {
	return &types.Interface{
		Name:       TestTitleName,
		ImportPath: TestImportPath,
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
