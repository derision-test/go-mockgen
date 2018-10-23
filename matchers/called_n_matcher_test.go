package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CalledNMatcherSuite struct{}

func (s *CalledNMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeCalledN(2).Match(litFunc{[]litCall{{}, {}}})
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *CalledNMatcherSuite) TestMatchEmptyHistory(t sweet.T) {
	ok, err := BeCalledN(1).Match(litFunc{})
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledNMatcherSuite) TestMatchMismatchedHistory(t sweet.T) {
	ok, err := BeCalledN(1).Match(litFunc{[]litCall{{}, {}}})
	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledNMatcherSuite) TestMatchError(t sweet.T) {
	_, err := BeCalledN(1).Match(nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(HavePrefix("calledNMatcher expects a mock function"))
}
