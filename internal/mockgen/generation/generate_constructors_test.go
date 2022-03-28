package generation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMockStructConstructor(t *testing.T) {
	code := generateMockStructConstructor(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))
	expected := strip(`
		// NewMockTestClient creates a new mock of the Client interface. All methods
		// return zero values for all results, unless overwritten.
		func NewMockTestClient() *MockTestClient {
			return &MockTestClient{
				StatusFunc: &TestClientStatusFunc{
					defaultHook: func() (r0 string, r1 bool) {
						return
					},
				},
				DoFunc: &TestClientDoFunc{
					defaultHook: func(string) (r0 bool) {
						return
					},
				},
				DofFunc: &TestClientDofFunc{
					defaultHook: func(string, ...string) (r0 bool) {
						return
					},
				},
			}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockStructStrictConstructor(t *testing.T) {
	code := generateMockStructStrictConstructor(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))
	expected := strip(`
		// NewStrictMockTestClient creates a new mock of the Client interface. All
		// methods panic on invocation, unless overwritten.
		func NewStrictMockTestClient() *MockTestClient {
			return &MockTestClient{
				StatusFunc: &TestClientStatusFunc{
					defaultHook: func() (string, bool) {
						panic("unexpected invocation of MockTestClient.Status")
					},
				},
				DoFunc: &TestClientDoFunc{
					defaultHook: func(string) bool {
						panic("unexpected invocation of MockTestClient.Do")
					},
				},
				DofFunc: &TestClientDofFunc{
					defaultHook: func(string, ...string) bool {
						panic("unexpected invocation of MockTestClient.Dof")
					},
				},
			}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockStructFromConstructor(t *testing.T) {
	code := generateMockStructFromConstructor(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))
	expected := strip(`
		// NewMockTestClientFrom creates a new mock of the MockTestClient interface.
		// All methods delegate to the given implementation, unless overwritten.
		func NewMockTestClientFrom(i test.Client) *MockTestClient {
			return &MockTestClient{
				StatusFunc: &TestClientStatusFunc{
					defaultHook: i.Status,
				},
				DoFunc: &TestClientDoFunc{
					defaultHook: i.Do,
				},
				DofFunc: &TestClientDofFunc{
					defaultHook: i.Dof,
				},
			}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockStructFromConstructorUnexported(t *testing.T) {
	iface := makeBareInterface(TestMethodStatus, TestMethodDo, TestMethodDof)
	iface.Name = "client"
	code := generateMockStructFromConstructor(wrapInterface(iface, TestPrefix, TestTitleName, TestMockStructName, ""), "")

	expected := strip(`
		// surrogateMockClient is a copy of the client interface (from the package
		// github.com/derision-test/go-mockgen/test). It is redefined here as it is
		// unexported in the source package.
		type surrogateMockClient interface {
			Status() (string, bool)
			Do(string) bool
			Dof(string, ...string) bool
		}

		// NewMockTestClientFrom creates a new mock of the MockTestClient interface.
		// All methods delegate to the given implementation, unless overwritten.
		func NewMockTestClientFrom(i surrogateMockClient) *MockTestClient {
			return &MockTestClient{
				StatusFunc: &TestClientStatusFunc{
					defaultHook: i.Status,
				},
				DoFunc: &TestClientDoFunc{
					defaultHook: i.Do,
				},
				DofFunc: &TestClientDofFunc{
					defaultHook: i.Dof,
				},
			}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}
