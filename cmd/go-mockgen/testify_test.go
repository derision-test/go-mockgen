package main

import (
	"fmt"
	"testing"

	"github.com/derision-test/go-mockgen/internal/testdata"
	"github.com/derision-test/go-mockgen/internal/testdata/mocks"
	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
)

func TestTestifyCalls(t *testing.T) {
	mock := mocks.NewMockClient()
	mockassert.NotCalled(t, mock.CloseFunc)
	assert.Nil(t, mock.Close())
	mockassert.Called(t, mock.CloseFunc)
	mockassert.CalledOnce(t, mock.CloseFunc)
}

func TestTestifyCallsWithArgs(t *testing.T) {
	mock := mocks.NewMockClient()
	mock.Do("foo")
	mockassert.Called(t, mock.DoFunc)
	mockassert.CalledOnce(t, mock.DoFunc)
	mockassert.CalledWith(t, mock.DoFunc, func(t assert.TestingT, v interface{}) bool { // TODO - ergonomics
		return v.(mocks.ClientDoFuncCall).Arg0 == "foo"
	})
	mockassert.NotCalledWith(t, mock.DoFunc, func(t assert.TestingT, v interface{}) bool { // TODO - ergonomics
		return v.(mocks.ClientDoFuncCall).Arg0 == "bar"
	})
}

func TestTestifyCallsWithVariadicArgs(t *testing.T) {
	mock := mocks.NewMockClient()
	mock.DoArgs("foo", 1, 2, 3)
	mockassert.CalledWith(t, mock.DoArgsFunc, func(t assert.TestingT, v interface{}) bool {
		// TODO - ergonomics
		call := v.(mocks.ClientDoArgsFuncCall)
		return call.Arg0 == "foo" && assert.Equal(t, []interface{}{1, 2, 3}, call.Arg1)
	})

	mock.DoArgs("bar", 42)
	mock.DoArgs("baz")
	mockassert.CalledN(t, mock.DoArgsFunc, 3)
	mockassert.CalledNWith(t, mock.DoArgsFunc, 2, func(t assert.TestingT, v interface{}) bool {
		// TODO - ergonomics
		return assert.Contains(t, v.(mocks.ClientDoArgsFuncCall).Arg0, "a")
	})

	// Mismatched variadic arg
	mockassert.NotCalledWith(t, mock.DoArgsFunc, func(t assert.TestingT, v interface{}) bool {
		// TODO - ergonomics
		call := v.(mocks.ClientDoArgsFuncCall)
		return call.Arg0 == "baz" && len(call.Arg1) > 0
	})
}

func TestTestifyPushHook(t *testing.T) {
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
		child, _ := parent.GetChild(0)
		assert.Equal(t, expected, child)
	}

	_, err := parent.GetChild(0)
	assert.EqualError(t, err, "uh-oh")
}

func TestTestifySetDefaultReturn(t *testing.T) {
	parent := mocks.NewMockParent()
	parent.GetChildFunc.SetDefaultReturn(nil, fmt.Errorf("uh-oh"))
	_, err := parent.GetChild(0)
	assert.EqualError(t, err, "uh-oh")
}

func TestTestifyPushReturn(t *testing.T) {
	parent := mocks.NewMockParent()
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil})
	parent.GetChildrenFunc.PushReturn([]testdata.Child{nil, nil, nil})

	assert.Len(t, parent.GetChildren(), 1)
	assert.Len(t, parent.GetChildren(), 2)
	assert.Len(t, parent.GetChildren(), 3)
	assert.Len(t, parent.GetChildren(), 0)
}
