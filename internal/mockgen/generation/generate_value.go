package generation

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

func generateZeroValue(typ types.Type, importPath, outputImportPath string) *jen.Statement {
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
			return compose(generateQualifiedName(t, importPath, outputImportPath), jen.Block())
		}

		return generateZeroValue(t.Underlying(), importPath, outputImportPath)

	case *types.Struct:
		return generateType(typ, importPath, outputImportPath, false).Block()
	}

	return jen.Nil()
}

func isIntegerType(kind types.BasicKind) bool {
	switch kind {
	case
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
		types.Complex64,
		types.Complex128:
		// Note:
		// ~types.Byte -> Uint8
		// ~types.Rune -> types.Int32
		return true

	default:
		return false
	}
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
