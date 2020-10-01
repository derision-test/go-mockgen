package testutil

type mockFunc struct {
	history []mockCall
}

type mockCall struct {
	args    []interface{}
	results []interface{}
}

func newHistory(calls ...mockCall) *mockFunc {
	return &mockFunc{
		history: calls,
	}
}

func (m mockFunc) History() []mockCall    { return m.history }
func (m mockCall) Args() []interface{}    { return m.args }
func (m mockCall) Results() []interface{} { return m.results }
