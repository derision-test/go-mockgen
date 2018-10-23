package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CalledOnceMatcherSuite struct{}

func (s *CalledOnceMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeCalledOnce().Match(litFunc{[]litCall{{}}})
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *CalledOnceMatcherSuite) TestMatchEmptyHistory(t sweet.T) {
	ok, err := BeCalledOnce().Match(litFunc{})
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledOnceMatcherSuite) TestMatchMismatchedHistory(t sweet.T) {
	ok, err := BeCalledOnce().Match(litFunc{[]litCall{{}, {}}})
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledOnceMatcherSuite) TestMatchError(t sweet.T) {
	_, err := BeCalledOnce().Match(nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(HavePrefix("calledNMatcher expects a mock function"))
}
