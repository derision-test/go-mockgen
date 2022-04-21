package generation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMockInterfaceMethod(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodDo)
	code := generateMockInterfaceMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Do delegates to the next hook function in the queue and stores the
		// parameter and result values of this invocation.
		func (m *MockTestClient) Do(v0 string) bool {
			r0 := m.DoFunc.nextHook()(v0)
			m.DoFunc.appendCall(TestClientDoFuncCall{v0, r0})
			return r0
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockInterfaceMethodVariadic(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodDof)
	code := generateMockInterfaceMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Dof delegates to the next hook function in the queue and stores the
		// parameter and result values of this invocation.
		func (m *MockTestClient) Dof(v0 string, v1 ...string) bool {
			r0 := m.DofFunc.nextHook()(v0, v1...)
			m.DofFunc.appendCall(TestClientDofFuncCall{v0, v1, r0})
			return r0
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}
