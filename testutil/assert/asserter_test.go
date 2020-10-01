package mockassert

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	asserter := Values(
		Skip,
		123,
		func(a string) bool { return len(a) == 2 },
	)

	assert.True(t, asserter.Assert(newMockArgs(context.Background(), 123, "xx")))
	assert.True(t, asserter.Assert(newMockArgs(context.Background(), 123, "yy")))
	assert.True(t, asserter.Assert(newMockArgs(context.Background(), 123, "zz")))
	assert.False(t, asserter.Assert(newMockArgs(context.Background(), 789, "w")))
	assert.False(t, asserter.Assert(newMockArgs(context.Background(), 123, 123)))
	assert.False(t, asserter.Assert(newMockArgs(context.Background(), 123, nil)))
}

func TestCallTesterFunc(t *testing.T) {
	v1, ok := callTesterFunc(func(v int) bool { return v%2 == 0 }, 4)
	assert.True(t, ok)
	assert.True(t, v1)

	v2, ok := callTesterFunc(func(v int) bool { return v%2 == 0 }, 3)
	assert.True(t, ok)
	assert.False(t, v2)
}

func TestCallTesterFuncNil(t *testing.T) {
	_, ok := callTesterFunc(nil, nil)
	assert.False(t, ok)
}

func TestCallTesterFuncNonFunc(t *testing.T) {
	_, ok := callTesterFunc(123, nil)
	assert.False(t, ok)
}

func TestCallTesterFuncBadParamArity(t *testing.T) {
	_, ok := callTesterFunc(func(a, b string) bool { return false }, nil)
	assert.False(t, ok)
}

func TestCallTesterFuncBadResultArity(t *testing.T) {
	_, ok := callTesterFunc(func(a string) (bool, error) { return false, nil }, nil)
	assert.False(t, ok)
}

func TestCallTesterFuncBadResultType(t *testing.T) {
	_, ok := callTesterFunc(func(a string) string { return a }, nil)
	assert.False(t, ok)
}

func TestCallTesterFuncMismatchedTypes(t *testing.T) {
	_, ok := callTesterFunc(func(a string) bool { return true }, 123)
	assert.False(t, ok)
}

type testArgs interface {
	Args() []interface{}
}

type mockArgs struct {
	args []interface{}
}

func newMockArgs(args ...interface{}) testArgs {
	return &mockArgs{
		args: args,
	}
}

func (m *mockArgs) Args() []interface{} {
	return m.args
}
