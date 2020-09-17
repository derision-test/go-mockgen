package generation

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

func GenerateZeroValue(typ types.Type, importPath, outputImportPath string) *jen.Statement {
	switch t := typ.(type) {
	case *types.Basic:
		kind := t.Kind()

		if kind == types.Bool {
			return jen.False()
		} else if kind == types.String {
			return jen.Lit("")
		} else if isIntegerType(kind) {
			return jen.Lit(0)
		}

	case *types.Named:
		if shouldEmitNamedType(t) {
			return Compose(generateQualifiedName(t, importPath, outputImportPath), jen.Block())
		}

		return GenerateZeroValue(t.Underlying(), importPath, outputImportPath)

	case *types.Struct:
		return GenerateType(typ, importPath, outputImportPath, false).Block()
	}

	return jen.Nil()
}

func isIntegerType(kind types.BasicKind) bool {
	kinds := []types.BasicKind{
		types.Int,
		types.Int8,
		types.Int16,
		types.Int32,
		types.Int64,
		types.Uint,
		types.Uint8,
		types.Uint16,
		types.Uint32,
		types.Uint64,
		types.Uintptr,
		types.Float32,
		types.Float64,
		types.Byte,
		types.Rune,
		types.Complex64,
		types.Complex128,
	}

	for _, k := range kinds {
		if k == kind {
			return true
		}
	}

	return false
}

func shouldEmitNamedType(t *types.Named) bool {
	if _, ok := t.Underlying().(*types.Struct); ok {
		return true
	}

	if _, ok := t.Underlying().(*types.Array); ok {
		return true
	}

	return false
}
