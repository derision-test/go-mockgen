package main

import (
	"testing"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.AddSuite(&NamesSuite{})
	})
}
