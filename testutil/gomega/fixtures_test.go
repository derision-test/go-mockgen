package matchers

type (
	litFunc struct {
		history []litCall
	}

	litCall struct {
		args    []interface{}
		results []interface{}
	}
)

func (f litFunc) History() []litCall     { return f.history }
func (i litCall) Args() []interface{}    { return i.args }
func (i litCall) Results() []interface{} { return i.results }
