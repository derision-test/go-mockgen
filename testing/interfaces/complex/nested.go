package complex

import "fmt"

type SupInterface1 interface {
	Foo()
}

type SupInterface2 interface {
	Bar()
}

type SubInterface interface {
	SupInterface1
	SupInterface2
	fmt.Stringer

	Baz()
}
