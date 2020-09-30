package testdata

type OptionValidator interface {
	Validate(options ...Option) error
}

type Option struct{}
