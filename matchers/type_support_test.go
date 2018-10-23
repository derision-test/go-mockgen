package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type TypeSupportSuite struct{}

func (s *TypeSupportSuite) TestGetCallHistoryNil(t sweet.T) {
	_, ok := getCallHistory(nil)
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryNoHistoryMethod(t sweet.T) {
	_, ok := getCallHistory(&TestNoHistory{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadParamArity(t sweet.T) {
	_, ok := getCallHistory(&TestHistoryBadParamArity{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadResultArity(t sweet.T) {
	_, ok := getCallHistory(&TestHistoryBadResultArity{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryNonSliceResult(t sweet.T) {
	_, ok := getCallHistory(&TestGetCallHistoryNonSliceResult{})
	Expect(ok).To(BeFalse())
}

func (s *TypeSupportSuite) TestGetCallHistoryBadSliceTypes(t sweet.T) {
	_, ok := getCallHistory(&TestGetCallHistoryBadSliceTypes{})
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

func (h *TestHistoryBadParamArity) History(n int) []callInstance       { return nil }
func (h *TestHistoryBadResultArity) History() ([]callInstance, error)  { return nil, nil }
func (h *TestGetCallHistoryNonSliceResult) History() string            { return "" }
func (h *TestGetCallHistoryBadSliceTypes) History() []halfCallInstance { return nil }
func (h *halfCallInstance) Args() []interface{}                        { return nil }
