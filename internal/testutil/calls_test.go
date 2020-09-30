package testutil

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type TypeSupportSuite struct{}

func (s *TypeSupportSuite) TestGetCallHistory(t sweet.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	history, ok := GetCallHistory(m)
	Expect(ok).To(BeTrue())
	Expect(history).To(HaveLen(5))
}

func (s *TypeSupportSuite) TestGetCallHistoryWith(t sweet.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := GetCallHistoryWith(m, func(v CallInstance) bool { return len(v.Args()) == 3 })
	Expect(ok).To(BeTrue())
	Expect(matchingHistory).To(HaveLen(3))
}

func (s *TypeSupportSuite) TestGetCallHistoryNil(t sweet.T) {
	_, ok := GetCallHistory(nil)
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryNoHistoryMethod(t sweet.T) {
	_, ok := GetCallHistory(&TestNoHistory{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadParamArity(t sweet.T) {
	_, ok := GetCallHistory(&TestHistoryBadParamArity{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadResultArity(t sweet.T) {
	_, ok := GetCallHistory(&TestHistoryBadResultArity{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryNonSliceResult(t sweet.T) {
	_, ok := GetCallHistory(&TestGetCallHistoryNonSliceResult{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadSliceTypes(t sweet.T) {
	_, ok := GetCallHistory(&TestGetCallHistoryBadSliceTypes{})
	Expect(ok).To(BeFalse())
}

//
// Helpers

type (
	TestNoHistory                    struct{}
	TestHistoryBadParamArity         struct{}
	TestHistoryBadResultArity        struct{}
	TestGetCallHistoryNonSliceResult struct{}
	TestGetCallHistoryBadSliceTypes  struct{}
	halfCallInstance                 struct{}
)

func (h *TestHistoryBadParamArity) History(n int) []CallInstance       { return nil }
func (h *TestHistoryBadResultArity) History() ([]CallInstance, error)  { return nil, nil }
func (h *TestGetCallHistoryNonSliceResult) History() string            { return "" }
func (h *TestGetCallHistoryBadSliceTypes) History() []halfCallInstance { return nil }
func (h *halfCallInstance) Args() []interface{}                        { return nil }

//
//

type (
	litFunc struct {
		history []litCall
	}

	litCall struct {
		args    []interface{}
		results []interface{}
	}
)

func (f litFunc) History() []litCall     { return f.history }
func (i litCall) Args() []interface{}    { return i.args }
func (i litCall) Results() []interface{} { return i.results }
