package main

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ParserSuite struct{}

func (s *ParserSuite) TestSimple(t sweet.T) {
	pkg, _, err := parseDir("./testing/parser/nonnested")
	Expect(err).To(BeNil())
	Expect(pkg.Name).To(Equal("nonnested"))
	Expect(pkg.Files).To(HaveLen(3))
	Expect(pkg.Files).To(HaveKey("testing/parser/nonnested/x.go"))
	Expect(pkg.Files).To(HaveKey("testing/parser/nonnested/y.go"))
	Expect(pkg.Files).To(HaveKey("testing/parser/nonnested/z.go"))
}

func (s *ParserSuite) TestEmpty(t sweet.T) {
	_, _, err := parseDir("./testing/parser/empty")
	Expect(err).To(MatchError("could not import package ./testing/parser/empty"))
}

func (s *ParserSuite) TestTwoPackages(t sweet.T) {
	_, _, err := parseDir("./testing/parser/twopackages")
	Expect(err).To(MatchError("could not import package ./testing/parser/twopackages"))
}
