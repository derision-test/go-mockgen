package localtypes

type EmbeddedStruct interface {
	Foo(struct{ z X })
}

type InterfaceStruct interface {
	Bar() interface {
		Baz(y Y) struct{ z Z }
	}
}
