package complex

type EmbeddedTypes interface {
	Param(struct {
		X string
		Y bool
	})

	Result() (struct{ Z int }, error)
}
