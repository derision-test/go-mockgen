package main

import (
	"path/filepath"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type ParserSuite struct{}

func (s *ParserSuite) TestSimple(t sweet.T) {
	pkg, _, err := parseDir(testPath + "/parser/nonnested")
	Expect(err).To(BeNil())
	Expect(pkg.Name).To(Equal("nonnested"))
	Expect(pkg.Files).To(HaveLen(3))
	Expect(pkg.Files).To(HaveKey(filepath.Join(gopath(), "src", testPath+"/parser/nonnested/x.go")))
	Expect(pkg.Files).To(HaveKey(filepath.Join(gopath(), "src", testPath+"/parser/nonnested/y.go")))
	Expect(pkg.Files).To(HaveKey(filepath.Join(gopath(), "src", testPath+"/parser/nonnested/z.go")))
}

func (s *ParserSuite) TestEmpty(t sweet.T) {
	_, _, err := parseDir(testPath + "/parser/empty")
	Expect(err).NotTo(BeNil())
}

func (s *ParserSuite) TestTwoPackages(t sweet.T) {
	_, _, err := parseDir(testPath + "/parser/twopackages")
	Expect(err).NotTo(BeNil())
}
