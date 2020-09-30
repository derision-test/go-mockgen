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
	Expect(err.Error()).To(HavePrefix("BeCalled expects a mock function"))
}

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
	Expect(err.Error()).To(HavePrefix("BeCalledN expects a mock function"))
}

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
	Expect(err.Error()).To(HavePrefix("BeCalledOnce expects a mock function"))
}
