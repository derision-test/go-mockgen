package testdata

type InterfaceType interface {
	ComplexParam(interface{ M(int) bool })
	CopmlexResult() (interface{ M(int) bool }, error)
}
