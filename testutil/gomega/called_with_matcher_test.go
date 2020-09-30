package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CalledWithMatcherSuite struct{}

func (s *CalledWithMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeCalledWith(ContainSubstring("foo"), 1, Not(Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"barbaz", 1, 2, 3}, nil},
	}})

	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}

func (s *CalledWithMatcherSuite) TestNoMatch(t sweet.T) {
	ok, err := BeCalledWith(ContainSubstring("foo"), 1, Not(Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"barbaz", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	Expect(err).To(BeNil())
	Expect(ok).To(BeFalse())
}

func (s *CalledWithMatcherSuite) TestMatchError(t sweet.T) {
	_, err := BeCalledWith("foo", 1, 2, 3).Match(nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(HavePrefix("BeCalledWith expects a mock function"))
}

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
	Expect(err.Error()).To(HavePrefix("BeCalledOnceWith expects a mock function"))
}

type CallsSuite struct{}

func (s *CallsSuite) TestGetMatchingCallCountsLiterals(t sweet.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, "foo", "bar", "baz")
	Expect(ok).To(BeTrue())
	Expect(matchingHistory).To(HaveLen(3))
}

func (s *CallsSuite) TestGetMatchingCallCountsMatchers(t sweet.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, HaveLen(3), HaveLen(3), HaveLen(3))
	Expect(ok).To(BeTrue())
	Expect(matchingHistory).To(HaveLen(3))
}

func (s *CallsSuite) TestGetMatchingCallCountsMixed(t sweet.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, "foo", "bar", ContainSubstring("bo"))
	Expect(ok).To(BeTrue())
	Expect(matchingHistory).To(HaveLen(1))
}
