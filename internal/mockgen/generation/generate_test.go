package generation

import (
	"fmt"
	"testing"

	"github.com/dave/jennifer/jen"
	"github.com/stretchr/testify/assert"
)

func TestGenerateInterface(t *testing.T) {
	expectedDecls := []string{
		// Structs
		"type MockTestClient struct",
		"type TestClientDoFunc struct",
		"type TestClientDoFuncCall struct",
		"type TestClientDofFunc struct",
		"type TestClientDofFuncCall struct",
		"func NewMockTestClient() *MockTestClient",
		// Overrides
		"func (m *MockTestClient) Do(v0 string) bool",
		"func (m *MockTestClient) Dof(v0 string, v1 ...string) bool",
		// DoFunc Methods
		"func (f *TestClientDoFunc) SetDefaultHook(hook func(string) bool)",
		"func (f *TestClientDoFunc) PushHook(hook func(string) bool)",
		"func (f *TestClientDoFunc) SetDefaultReturn(r0 bool)",
		"func (f *TestClientDoFunc) PushReturn(r0 bool)",
		"func (f *TestClientDoFunc) History() []TestClientDoFuncCall",
		// DoFuncCall methods
		"func (c TestClientDoFuncCall) Args() []interface{}",
		"func (c TestClientDoFuncCall) Results() []interface{}",
		// DofFunc Methods
		"func (f *TestClientDofFunc) SetDefaultHook(hook func(string, ...string) bool)",
		"func (f *TestClientDofFunc) PushHook(hook func(string, ...string) bool)",
		"func (f *TestClientDofFunc) SetDefaultReturn(r0 bool)",
		"func (f *TestClientDofFunc) PushReturn(r0 bool)",
		"func (f *TestClientDofFunc) History() []TestClientDofFuncCall",
		// DofFuncCall methods
		"func (c TestClientDofFuncCall) Args() []interface{}",
		"func (c TestClientDofFuncCall) Results() []interface{}",
	}

	file := jen.NewFile("test")

	generateInterface(file, makeBareInterface(TestMethodDo, TestMethodDof), TestPrefix, "")
	rendered := fmt.Sprintf("%#v\n", file)

	for _, decl := range expectedDecls {
		assert.Contains(t, rendered, decl)
	}
}
