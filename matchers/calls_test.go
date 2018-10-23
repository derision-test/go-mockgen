package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type CallsSuite struct{}

func (s *CallsSuite) TestGetMatchingCallCountsLiterals(t sweet.T) {
	count, err := getMatchingCallCounts(
		[]callInstance{
			litCall{args: []interface{}{"foo", "bar"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
			litCall{args: []interface{}{"foo", "bar", "baz", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
		},
		[]interface{}{"foo", "bar", "baz"},
	)

	Expect(err).To(BeNil())
	Expect(count).To(Equal(2))
}

func (s *CallsSuite) TestGetMatchingCallCountsMatchers(t sweet.T) {
	count, err := getMatchingCallCounts(
		[]callInstance{
			litCall{args: []interface{}{"foo", "bar"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
			litCall{args: []interface{}{"foo", "bar", "baz", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
		},
		[]interface{}{HaveLen(3), HaveLen(3), HaveLen(3)},
	)

	Expect(err).To(BeNil())
	Expect(count).To(Equal(2))
}

func (s *CallsSuite) TestGetMatchingCallCountsMixed(t sweet.T) {
	count, err := getMatchingCallCounts(
		[]callInstance{
			litCall{args: []interface{}{"foo", "bar"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
			litCall{args: []interface{}{"foo", "bar", "baz", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "bonk"}},
			litCall{args: []interface{}{"foo", "bar", "baz"}},
		},
		[]interface{}{"foo", ContainSubstring("ba"), "bonk"},
	)

	Expect(err).To(BeNil())
	Expect(count).To(Equal(1))
}

func (s *CallsSuite) TestGetMatchingCallCountsErrorNoMatch(t sweet.T) {
	_, err := getMatchingCallCounts(
		[]callInstance{
			litCall{args: []interface{}{1}},
		},
		[]interface{}{HaveLen(3)},
	)

	Expect(err).NotTo(BeNil())
	Expect(err.Error()).To(ContainSubstring("HaveLen matcher expects a string/array/map/channel/slice."))
}
