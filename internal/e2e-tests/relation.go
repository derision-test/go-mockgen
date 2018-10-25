package e2etests

type Parent interface {
	AddChild(c Child)
	GetChildren() []Child
	GetChild(i int) (Child, error)
}

type Child interface {
	Parent() Parent
}
