package main

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type NamesSuite struct{}

func (s *NamesSuite) TestNameExtractor(t sweet.T) {
	pkg, _, err := parseDir("./testing/names")
	Expect(err).To(BeNil())

	v := newNameExtractor()
	walk(pkg, v)
	Expect(v.names).To(ConsistOf([]string{
		"OuterType",
		"StructType",
		"IntefaceType",
		"SimpleType",
	}))
}
