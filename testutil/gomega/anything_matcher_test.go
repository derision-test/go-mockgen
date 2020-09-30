package matchers

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type AnythingMatcherSuite struct{}

func (s *AnythingMatcherSuite) TestMatch(t sweet.T) {
	ok, err := BeAnything().Match(nil)
	Expect(err).To(BeNil())
	Expect(ok).To(BeTrue())
}
