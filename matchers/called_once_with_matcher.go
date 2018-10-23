package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type calledOnceWithMatcher struct {
	args []interface{}
}

func BeCalledOnceWith(args ...interface{}) types.GomegaMatcher {
	return &calledOnceWithMatcher{
		args: args,
	}
}

func (m *calledOnceWithMatcher) Match(actual interface{}) (bool, error) {
	history, ok := getCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("calledOnceWithMatcher expects a mock function description.  Got:\n%s", format.Object(actual, 1))
	}

	n, err := getMatchingCallCounts(history, m.args)
	if n == 1 {
		return true, nil
	}

	return false, err
}

func (m *calledOnceWithMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to contain one call with argument list matching", m.args)
}

func (m *calledOnceWithMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to contain one call with argument list matching", m.args)
}
