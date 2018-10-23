package main

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"

	"github.com/efritz/go-mockgen/internal/e2e-tests/iface"
	"github.com/efritz/go-mockgen/internal/e2e-tests/mock"
)

type E2ESuite struct{}

func (s *E2ESuite) TestCalls(t sweet.T) {
	mock := mock.NewMockClient()
	Expect(mock.CloseFunc).NotTo(BeCalled())
	Expect(mock.Close()).To(BeNil())
	Expect(mock.CloseFunc).To(BeCalled())
	Expect(mock.CloseFunc).To(BeCalledOnce())
	Expect(mock.CloseFunc).To(BeCalledWith())
}

func (s *E2ESuite) TestCallsWithArgs(t sweet.T) {
	mock := mock.NewMockClient()
	mock.Do("foo")
	Expect(mock.DoFunc).To(BeCalled())
	Expect(mock.DoFunc).To(BeCalledWith("foo"))
}

func (s *E2ESuite) TestCallsWithVariadicArgs(t sweet.T) {
	mock := mock.NewMockClient()
	mock.DoArgs("foo", 1, 2, 3)
	Expect(mock.DoArgsFunc).To(BeCalledWith("foo", 1, 2, 3))
	Expect(mock.DoArgsFunc).To(BeCalledWith(Equal("foo"), Equal(1), Equal(2), Equal(3)))

	mock.DoArgs("bar", 42)
	mock.DoArgs("baz")
	Expect(mock.DoArgsFunc).To(BeCalledN(3))
	Expect(mock.DoArgsFunc).To(BeCalledOnceWith(ContainSubstring("a")))
	Expect(mock.DoArgsFunc).To(BeCalledOnceWith(ContainSubstring("a"), BeAnything()))

	// Mismatched variadic arg
	Expect(mock.DoArgsFunc).NotTo(BeCalledWith("baz", BeAnything()))
}

func (s *E2ESuite) TestPushHook(t sweet.T) {
	child1 := mock.NewMockChild()
	child2 := mock.NewMockChild()
	child3 := mock.NewMockChild()
	parent := mock.NewMockParent()

	parent.GetChildFunc.PushHook(func(i int) (iface.Child, error) { return child1, nil })
	parent.GetChildFunc.PushHook(func(i int) (iface.Child, error) { return child2, nil })
	parent.GetChildFunc.PushHook(func(i int) (iface.Child, error) { return child3, nil })

	parent.GetChildFunc.SetDefaultHook(func(i int) (iface.Child, error) {
		return nil, fmt.Errorf("utoh")
	})

	Expect(parent.GetChild(0)).To(Equal(child1))
	Expect(parent.GetChild(0)).To(Equal(child2))
	Expect(parent.GetChild(0)).To(Equal(child3))

	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("utoh"))
}

func (s *E2ESuite) TestSetDefaultReturn(t sweet.T) {
	parent := mock.NewMockParent()
	parent.GetChildFunc.SetDefaultReturn(nil, fmt.Errorf("utoh"))
	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("utoh"))
}

func (s *E2ESuite) TestPushReturn(t sweet.T) {
	parent := mock.NewMockParent()
	parent.GetChildrenFunc.PushReturn([]iface.Child{nil})
	parent.GetChildrenFunc.PushReturn([]iface.Child{nil, nil})
	parent.GetChildrenFunc.PushReturn([]iface.Child{nil, nil, nil})

	Expect(parent.GetChildren()).To(HaveLen(1))
	Expect(parent.GetChildren()).To(HaveLen(2))
	Expect(parent.GetChildren()).To(HaveLen(3))
	Expect(parent.GetChildren()).To(HaveLen(0))
}
