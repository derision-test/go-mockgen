package matchers

import (
	"fmt"

	"github.com/derision-test/go-mockgen/v2/internal/testutil"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

type calledWithMatcher struct {
	name string
	args []interface{}
}

var _ types.GomegaMatcher = &calledWithMatcher{}

// BeCalledWith constructs a matcher that asserts the mock function object was called at least once
// with a set of arguments matching the given values. The values can be another matcher or a literal
// value. In the latter case, the values will be checked for equality.
func BeCalledWith(args ...interface{}) types.GomegaMatcher {
	return &calledWithMatcher{
		name: "BeCalledWith",
		args: args,
	}
}

func (m *calledWithMatcher) Match(actual interface{}) (bool, error) {
	matchingHistory, ok := getCallHistoryWith(actual, m.args...)
	if !ok {
		return false, fmt.Errorf("%s expects a mock function description. Got:\n%s", m.name, format.Object(actual, 1))
	}

	return len(matchingHistory) > 0, nil
}

func (m *calledWithMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to contain at least one call with argument list matching", m.args)
}

func (m *calledWithMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to contain at least one call with argument list matching", m.args)
}

// BeCalledOnceWith constructs a matcher that asserts the mock function object was called exactly once
// with a set of arguments matching the given values. The values can be another matcher or a literal
// value. In the latter case, the values will be checked for equality.
func BeCalledOnceWith(args ...interface{}) types.GomegaMatcher {
	return &calledNWithMatcher{
		name: "BeCalledOnceWith",
		n:    1,
		args: args,
	}
}

type calledNWithMatcher struct {
	name string
	n    int
	args []interface{}
}

var _ types.GomegaMatcher = &calledNWithMatcher{}

// BeCalledNWith constructs a matcher that asserts the mock function object was called exactly n times
// with a set of arguments matching the given values. The values can be another matcher or a literal
// value. In the latter case, the values will be checked for equality.
func BeCalledNWith(n int, args ...interface{}) types.GomegaMatcher {
	return &calledNWithMatcher{
		name: "BeCalledNWith",
		n:    n,
		args: args,
	}
}

func (m *calledNWithMatcher) Match(actual interface{}) (bool, error) {
	matchingHistory, ok := getCallHistoryWith(actual, m.args...)
	if !ok {
		return false, fmt.Errorf("%s expects a mock function description. Got:\n%s", m.name, format.Object(actual, 1))
	}

	return len(matchingHistory) == m.n, nil
}

func (m *calledNWithMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to contain one call with argument list matching", m.args)
}

func (m *calledNWithMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to contain one call with argument list matching", m.args)
}

// getCallHistoryWith returns the set of call instances matching the given values. The values can
// be another matcher or a literal value. In the latter case, the values will be checked for equality.
func getCallHistoryWith(actual interface{}, args ...interface{}) ([]testutil.CallInstance, bool) {
	return testutil.GetCallHistoryWith(actual, func(v testutil.CallInstance) bool {
		if len(args) > len(v.Args()) {
			return false
		}

		for i, expectedArg := range args {
			matcher, ok := expectedArg.(types.GomegaMatcher)
			if !ok {
				matcher = &matchers.EqualMatcher{Expected: expectedArg}
			}

			success, err := matcher.Match(v.Args()[i])
			if err != nil || !success {
				return false
			}
		}

		return true
	})
}
