package matchers

import (
	"fmt"

	"github.com/derision-test/go-mockgen/internal/testutil"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type calledMatcher struct {
	name string
}

var _ types.GomegaMatcher = &calledMatcher{}

// BeCalled constructs a matcher that asserts the mock function object was called at least once.
func BeCalled() types.GomegaMatcher {
	return &calledMatcher{
		name: "BeCalled",
	}
}

func (m *calledMatcher) Match(actual interface{}) (bool, error) {
	history, ok := testutil.GetCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("%s expects a mock function description. Got:\n%s", m.name, format.Object(actual, 1))
	}

	return len(history) > 0, nil
}

func (m *calledMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nto be called at least once", format.Object(actual, 1))
}

func (m *calledMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nnot to be called at least once", format.Object(actual, 1))
}

// BeCalledOnce constructs a matcher that asserts the mock function object was called exactly once.
func BeCalledOnce() types.GomegaMatcher {
	return &calledNMatcher{
		name: "BeCalledOnce",
		n:    1,
	}
}

type calledNMatcher struct {
	name string
	n    int
}

var _ types.GomegaMatcher = &calledNMatcher{}

// BeCalledN constructs a matcher that asserts the mock function object was called exactly n times.
func BeCalledN(n int) types.GomegaMatcher {
	return &calledNMatcher{
		name: "BeCalledN",
		n:    n,
	}
}

func (m *calledNMatcher) Match(actual interface{}) (bool, error) {
	history, ok := testutil.GetCallHistory(actual)
	if !ok {
		return false, fmt.Errorf("%s expects a mock function description. Got:\n%s", m.name, format.Object(actual, 1))
	}

	return len(history) == m.n, nil
}

func (m *calledNMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nto be called %d times", format.Object(actual, 1), m.n)
}

func (m *calledNMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nnot to be called %d times", format.Object(actual, 1), m.n)
}
