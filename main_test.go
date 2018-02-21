package main

import (
	"testing"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

var testPath = "github.com/efritz/go-mockgen/testing"

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.AddSuite(&ParserSuite{})
		s.AddSuite(&NamesSuite{})
		s.AddSuite(&InterfacesSuite{})
	})
}
