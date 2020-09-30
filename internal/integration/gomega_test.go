package integration

import (
	"fmt"
	"testing"

	"github.com/derision-test/go-mockgen/internal/integration/testdata"
	"github.com/derision-test/go-mockgen/internal/integration/testdata/mocks"
	. "github.com/derision-test/go-mockgen/testutil/gomega"
	. "github.com/onsi/gomega"
)

func TestGomegaCalls(t *testing.T) {
	RegisterTestingT(t)

	mock := mocks.NewMockClient()
	Expect(mock.CloseFunc).NotTo(BeCalled())
	Expect(mock.Close()).To(BeNil())
	Expect(mock.CloseFunc).To(BeCalled())
	Expect(mock.CloseFunc).To(BeCalledOnce())
}

func TestGomegaCallsWithArgs(t *testing.T) {
	RegisterTestingT(t)

	mock := mocks.NewMockClient()
	mock.Do("foo")
	Expect(mock.DoFunc).To(BeCalled())
	Expect(mock.DoFunc).To(BeCalledOnce())
	Expect(mock.DoFunc).To(BeCalledWith("foo"))
	Expect(mock.DoFunc).NotTo(BeCalledWith("bar"))
}

func TestGomegaCallsWithVariadicArgs(t *testing.T) {
	RegisterTestingT(t)

	mock := mocks.NewMockClient()
	mock.DoArgs("foo", 1, 2, 3)
	Expect(mock.DoArgsFunc).To(BeCalledWith("foo", 1, 2, 3))
	Expect(mock.DoArgsFunc).To(BeCalledWith(Equal("foo"), Equal(1), Equal(2), Equal(3)))

	mock.DoArgs("bar", 42)
	mock.DoArgs("baz")
	Expect(mock.DoArgsFunc).To(BeCalledN(3))
	Expect(mock.DoArgsFunc).To(BeCalledNWith(2, ContainSubstring("a")))

	// Mismatched variadic arg
	Expect(mock.DoArgsFunc).NotTo(BeCalledWith("baz", BeAnything()))
}

func TestGomegaPushHook(t *testing.T) {
	RegisterTestingT(t)

	child1 := mocks.NewMockChild()
	child2 := mocks.NewMockChild()
	child3 := mocks.NewMockChild()
	parent := mocks.NewMockParent()

	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child1, nil })
	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child2, nil })
	parent.GetChildFunc.PushHook(func(i int) (testdata.Child, error) { return child3, nil })
	parent.GetChildFunc.SetDefaultHook(func(i int) (testdata.Child, error) {
		return nil, fmt.Errorf("uh-oh")
	})

	for _, expected := range []interface{}{child1, child2, child3} {
		Expect(parent.GetChild(0)).To(Equal(expected))
	}

	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("uh-oh"))
}

func TestGomegaSetDefaultReturn(t *testing.T) {
	RegisterTestingT(t)

	parent := mocks.NewMockParent()
	parent.GetChildFunc.SetDefaultReturn(nil, fmt.Errorf("uh-oh"))
	_, err := parent.GetChild(0)
	Expect(err).To(MatchError("uh-oh"))
}

func TestGomegaPushReturn(t *testing.T) {
	RegisterTestingT(t)

	parent := mocks.NewMockParent()
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil, nil})

	Expect(parent.GetChildren()).To(HaveLen(1))
	Expect(parent.GetChildren()).To(HaveLen(2))
	Expect(parent.GetChildren()).To(HaveLen(3))
	Expect(parent.GetChildren()).To(HaveLen(0))
}
