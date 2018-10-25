package e2etests

type OptionValidator interface {
	Validate(options ...Option) error
}

type Option struct{}
