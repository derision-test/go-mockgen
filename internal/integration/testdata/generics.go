package testdata

type I1[T any] interface {
	M1(T)
}

type I2[T1, T2 any] interface {
	I1[T1]
	M2(T1) T2
}

type SignedIntegerConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type I3[T SignedIntegerConstraint] interface {
	M3(I1[T]) int
}

type I4[T interface{ ~string }] interface {
	M4(I1[T]) int
}

type fooer[T any] interface {
	Foo() T
}
