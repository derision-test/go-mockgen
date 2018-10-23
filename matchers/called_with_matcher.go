package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type calledWithMatcher struct {
	args []interface{}
}

func BeCalledWith(args ...interface{}) types.GomegaMatcher {
	return &calledWithMatcher{
		args: args,
	}
}

func (m *calledWithMatcher) Match(actual interface{}) (bool, error) {
	history, ok := getCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("calledWithMatcher expects a mock function description.  Got:\n%s", format.Object(actual, 1))
	}

	n, err := getMatchingCallCounts(history, m.args)
	if n > 0 {
		return true, nil
	}

	return false, err
}

func (m *calledWithMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to contain at least one call with argument list matching", m.args)
}

func (m *calledWithMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to contain at least one call with argument list matching", m.args)
}
