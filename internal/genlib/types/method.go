package types

import "go/types"

type Method struct {
	Name     string
	Params   []types.Type
	Results  []types.Type
	Variadic bool
}

func DeconstructMethod(name string, signature *types.Signature) *Method {
	var (
		ps      = signature.Params()
		rs      = signature.Results()
		params  = []types.Type{}
		results = []types.Type{}
	)

	for i := 0; i < ps.Len(); i++ {
		params = append(params, ps.At(i).Type())
	}

	for i := 0; i < rs.Len(); i++ {
		results = append(results, rs.At(i).Type())
	}

	return &Method{
		Name:     name,
		Params:   params,
		Results:  results,
		Variadic: signature.Variadic(),
	}
}
