package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type calledMatcher struct{}

func BeCalled() types.GomegaMatcher {
	return &calledMatcher{}
}

func (m *calledMatcher) Match(actual interface{}) (bool, error) {
	history, ok := getCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("calledMatcher expects a mock function description.  Got:\n%s", format.Object(actual, 1))
	}

	return len(history) > 0, nil
}

func (m *calledMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nto be called at least once", format.Object(actual, 1))
}

func (m *calledMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nnot to be called at least once", format.Object(actual, 1))
}
