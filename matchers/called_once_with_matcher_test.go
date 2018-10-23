package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CalledOnceWithMatcherSuite struct{}

func (s *CalledOnceWithMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeCalledOnceWith(ContainSubstring("foo"), 1, Not(Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"barbaz", 1, 2, 3}, nil},
	}})

	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *CalledOnceWithMatcherSuite) TestNoMatch(t sweet.T) {
	ok, err := BeCalledOnceWith(ContainSubstring("foo"), 1, Not(Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"barbaz", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledOnceWithMatcherSuite) TestMultipleMatches(t sweet.T) {
	ok, err := BeCalledOnceWith(ContainSubstring("foo"), 1, Not(Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledOnceWithMatcherSuite) TestMatchError(t sweet.T) {
	_, err := BeCalledOnceWith("foo", 1, 2, 3).Match(nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(HavePrefix("calledOnceWithMatcher expects a mock function"))
}
