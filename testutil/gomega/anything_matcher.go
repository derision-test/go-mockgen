package matchers

import (
	"fmt"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type anythingMatcher struct {
	n int
}

// BeAnything returns a matcher that never fails.
func BeAnything() types.GomegaMatcher {
	return &anythingMatcher{}
}

func (m *anythingMatcher) Match(actual interface{}) (bool, error) {
	return true, nil
}

func (m *anythingMatcher) FailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nto be anything", format.Object(actual, 1))
}

func (m *anythingMatcher) NegatedFailureMessage(actual interface{}) string {
	return fmt.Sprintf("Expected\n%s\nnot to be anything", format.Object(actual, 1))
}
