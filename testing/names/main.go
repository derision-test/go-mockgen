package something

type OuterType func() error

type (
	StructType     struct{}
	IntefaceType   interface{}
	SimpleType     int
	unexportedType int
)
