package mockgen

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type UtilSuite struct{}

func (s *UtilSuite) TestTitle(t sweet.T) {
	Expect(title("")).To(Equal(""))
	Expect(title("foobar")).To(Equal("Foobar"))
	Expect(title("fooBar")).To(Equal("FooBar"))
}
