package generation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMockStruct(t *testing.T) {
	code := generateMockStruct(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))
	expected := strip(`
		// MockTestClient is a mock implementation of the Client interface (from the
		// package github.com/derision-test/go-mockgen/test) used for unit testing.
		type MockTestClient struct {
			// StatusFunc is an instance of a mock function object controlling the
			// behavior of the method Status.
			StatusFunc *TestClientStatusFunc
			// DoFunc is an instance of a mock function object controlling the
			// behavior of the method Do.
			DoFunc *TestClientDoFunc
			// DofFunc is an instance of a mock function object controlling the
			// behavior of the method Dof.
			DofFunc *TestClientDofFunc
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateFuncStruct(t *testing.T) {
	code := generateMockFuncStruct(makeMethod(TestMethodDo))
	expected := strip(`
		// TestClientDoFunc describes the behavior when the Do method of the parent
		// MockTestClient instance is invoked.
		type TestClientDoFunc struct {
			defaultHook func(string) bool
			hooks       []func(string) bool
			history     []TestClientDoFuncCall
			mutex       sync.Mutex
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateFuncStructVariadic(t *testing.T) {
	code := generateMockFuncStruct(makeMethod(TestMethodDof))
	expected := strip(`
		// TestClientDofFunc describes the behavior when the Dof method of the
		// parent MockTestClient instance is invoked.
		type TestClientDofFunc struct {
			defaultHook func(string, ...string) bool
			hooks       []func(string, ...string) bool
			history     []TestClientDofFuncCall
			mutex       sync.Mutex
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncCallStruct(t *testing.T) {
	code := generateMockFuncCallStruct(makeMethod(TestMethodDo))
	expected := strip(`
		// TestClientDoFuncCall is an object that describes an invocation of method
		// Do on an instance of MockTestClient.
		type TestClientDoFuncCall struct {
			// Arg0 is the value of the 1st argument passed to this method
			// invocation.
			Arg0 string
			// Result0 is the value of the 1st result returned from this method
			// invocation.
			Result0 bool
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncCallStructVariadic(t *testing.T) {
	code := generateMockFuncCallStruct(makeMethod(TestMethodDof))
	expected := strip(`
		// TestClientDofFuncCall is an object that describes an invocation of method
		// Dof on an instance of MockTestClient.
		type TestClientDofFuncCall struct {
			// Arg0 is the value of the 1st argument passed to this method
			// invocation.
			Arg0 string
			// Arg1 is a slice containing the values of the variadic arguments
			// passed to this method invocation.
			Arg1 []string
			// Result0 is the value of the 1st result returned from this method
			// invocation.
			Result0 bool
		}
	`)

	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}
