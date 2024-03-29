package generation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMockFuncCallArgsMethod(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodDo)
	code := generateMockFuncCallArgsMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Args returns an interface slice containing the arguments of this
		// invocation.
		func (c TestClientDoFuncCall) Args() []interface{} {
			return []interface{}{c.Arg0}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncCallArgsMethodVariadic(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodDof)
	code := generateMockFuncCallArgsMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Args returns an interface slice containing the arguments of this
		// invocation. The variadic slice argument is flattened in this array such
		// that one positional argument and three variadic arguments would result in
		// a slice of four, not two.
		func (c TestClientDofFuncCall) Args() []interface{} {
			trailing := []interface{}{}
			for _, val := range c.Arg1 {
				trailing = append(trailing, val)
			}

			return append([]interface{}{c.Arg0}, trailing...)
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncCallResultsMethod(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodDo)
	code := generateMockFuncCallResultsMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Results returns an interface slice containing the results of this
		// invocation.
		func (c TestClientDoFuncCall) Results() []interface{} {
			return []interface{}{c.Result0}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncCallResultsMethodMultiple(t *testing.T) {
	wrappedInterface := makeInterface(TestMethodStatus)
	code := generateMockFuncCallResultsMethod(wrappedInterface, wrappedInterface.wrappedMethods[0], "")
	expected := strip(`
		// Results returns an interface slice containing the results of this
		// invocation.
		func (c TestClientStatusFuncCall) Results() []interface{} {
			return []interface{}{c.Result0, c.Result1}
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}
