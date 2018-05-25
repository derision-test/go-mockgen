package generation

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

func zeroValue(typ types.Type, importPath string) *jen.Statement {
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
			return compose(generateQualifiedName(t, importPath), jen.Block())
		}

		return zeroValue(t.Underlying(), importPath)

	case *types.Struct:
		return generateType(typ, importPath, false).Block()
	}

	return jen.Nil()
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
