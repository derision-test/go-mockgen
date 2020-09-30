package main

import (
	"fmt"

	"github.com/aphistic/sweet"
	"github.com/derision-test/go-mockgen/internal/testdata"
	"github.com/derision-test/go-mockgen/internal/testdata/mocks"
	. "github.com/derision-test/go-mockgen/testutil/gomega"
	. "github.com/onsi/gomega"
)

type BinaryTestSuite struct{}

func (s *BinaryTestSuite) TestCalls(t sweet.T) {
	mock := mocks.NewMockClient()
	Expect(mock.CloseFunc).NotTo(BeCalled())
	Expect(mock.Close()).To(BeNil())
	Expect(mock.CloseFunc).To(BeCalled())
	Expect(mock.CloseFunc).To(BeCalledOnce())
	Expect(mock.CloseFunc).To(BeCalledWith())
}

func (s *BinaryTestSuite) TestCallsWithArgs(t sweet.T) {
	mock := mocks.NewMockClient()
	mock.Do("foo")
	Expect(mock.DoFunc).To(BeCalled())
	Expect(mock.DoFunc).To(BeCalledWith("foo"))
}

func (s *BinaryTestSuite) TestCallsWithVariadicArgs(t sweet.T) {
	mock := mocks.NewMockClient()
	mock.DoArgs("foo", 1, 2, 3)
	Expect(mock.DoArgsFunc).To(BeCalledWith("foo", 1, 2, 3))
	Expect(mock.DoArgsFunc).To(BeCalledWith(Equal("foo"), Equal(1), Equal(2), Equal(3)))

	mock.DoArgs("bar", 42)
	mock.DoArgs("baz")
	Expect(mock.DoArgsFunc).To(BeCalledN(3))
	Expect(mock.DoArgsFunc).To(BeCalledOnceWith(ContainSubstring("a")))

	// Mismatched variadic arg
	Expect(mock.DoArgsFunc).NotTo(BeCalledWith("baz", BeAnything()))
}

func (s *BinaryTestSuite) TestPushHook(t sweet.T) {
	child1 := mocks.NewMockChild()
	child2 := mocks.NewMockChild()
	child3 := mocks.NewMockChild()
	parent := mocks.NewMockParent()

	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child1, nil })
	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child2, nil })
	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child3, nil })

	parent.GetChildFunc.SetDefaultHook(func(i int) (testdata.Child, error) {
		return nil, fmt.Errorf("utoh")
	})

	Expect(parent.GetChild(0)).To(Equal(child1))
	Expect(parent.GetChild(0)).To(Equal(child2))
	Expect(parent.GetChild(0)).To(Equal(child3))

	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("utoh"))
}

func (s *BinaryTestSuite) TestSetDefaultReturn(t sweet.T) {
	parent := mocks.NewMockParent()
	parent.GetChildFunc.SetDefaultReturn(nil, fmt.Errorf("utoh"))
	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("utoh"))
}

func (s *BinaryTestSuite) TestPushReturn(t sweet.T) {
	parent := mocks.NewMockParent()
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil, nil})

	Expect(parent.GetChildren()).To(HaveLen(1))
	Expect(parent.GetChildren()).To(HaveLen(2))
	Expect(parent.GetChildren()).To(HaveLen(3))
	Expect(parent.GetChildren()).To(HaveLen(0))
}
