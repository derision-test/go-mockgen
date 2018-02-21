package simple

type NoParams interface {
	NoParamsNoResults()
	NoParamsOneResult() error
	NoParamsMultipleResults() ([]string, error)
}

type Params interface {
	OneParam(foo string)
	MultipleParams(foo, bar string, baz bool)
	Unnamed(string, string, bool)
	Variadic(format string, params ...interface{})
}
