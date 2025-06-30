package testdata

type AnyReference interface {
	Get(key string) ([]byte, error)
	Set(key string, value any) error // any alias
}

type Foo struct {
	ID   string
	Name string
}

type Foos = []Foo

type Bar struct {
	ID   string
	Name string
}

type Bars Bars

type AliasReference interface {
	GetFoos() (Foos, error)
	GetBars() (Bars, error)
}
