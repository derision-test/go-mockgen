package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type calledNMatcher struct {
	n int
}

func BeCalledN(n int) types.GomegaMatcher {
	return &calledNMatcher{
		n: n,
	}
}

func (m *calledNMatcher) Match(actual interface{}) (bool, error) {
	history, ok := getCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("calledNMatcher expects a mock function description.  Got:\n%s", format.Object(actual, 1))
	}

	return len(history) == m.n, nil
}

func (m *calledNMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nto be called %d times", format.Object(actual, 1), m.n)
}

func (m *calledNMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nnot to be called %d times", format.Object(actual, 1), m.n)
}
