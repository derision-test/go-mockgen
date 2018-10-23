package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CalledMatcherSuite struct{}

func (s *CalledMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeCalled().Match(litFunc{[]litCall{{}}})
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *CalledMatcherSuite) TestMatchEmptyHistory(t sweet.T) {
	ok, err := BeCalled().Match(litFunc{})
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledMatcherSuite) TestMatchError(t sweet.T) {
	_, err := BeCalled().Match(nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(HavePrefix("calledMatcher expects a mock function"))
}
