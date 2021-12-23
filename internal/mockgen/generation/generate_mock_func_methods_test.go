package generation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMockFuncSetHookMethod(t *testing.T) {
	code := generateMockFuncSetHookMethod(makeMethod(TestMethodDo))
	expected := strip(`
		// SetDefaultHook sets function that is called when the Do method of the
		// parent MockTestClient instance is invoked and the hook queue is empty.
		func (f *TestClientDoFunc) SetDefaultHook(hook func(string) bool) {
			f.defaultHook = hook
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncSetHookMethodVariadic(t *testing.T) {
	code := generateMockFuncSetHookMethod(makeMethod(TestMethodDof))
	expected := strip(`
		// SetDefaultHook sets function that is called when the Dof method of the
		// parent MockTestClient instance is invoked and the hook queue is empty.
		func (f *TestClientDofFunc) SetDefaultHook(hook func(string, ...string) bool) {
			f.defaultHook = hook
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncPushHookMethod(t *testing.T) {
	code := generateMockFuncPushHookMethod(makeMethod(TestMethodDo))
	expected := strip(`
		// PushHook adds a function to the end of hook queue. Each invocation of the
		// Do method of the parent MockTestClient instance invokes the hook at the
		// front of the queue and discards it. After the queue is empty, the default
		// hook function is invoked for any future action.
		func (f *TestClientDoFunc) PushHook(hook func(string) bool) {
			f.mutex.Lock()
			f.hooks = append(f.hooks, hook)
			f.mutex.Unlock()
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncPushHookMethodVariadic(t *testing.T) {
	code := generateMockFuncPushHookMethod(makeMethod(TestMethodDof))
	expected := strip(`
		// PushHook adds a function to the end of hook queue. Each invocation of the
		// Dof method of the parent MockTestClient instance invokes the hook at the
		// front of the queue and discards it. After the queue is empty, the default
		// hook function is invoked for any future action.
		func (f *TestClientDofFunc) PushHook(hook func(string, ...string) bool) {
			f.mutex.Lock()
			f.hooks = append(f.hooks, hook)
			f.mutex.Unlock()
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncSetReturnMethod(t *testing.T) {
	code := generateMockFuncSetReturnMethod(makeMethod(TestMethodDo))
	expected := strip(`
		// SetDefaultReturn calls SetDefaultHook with a function that returns the
		// given values.
		func (f *TestClientDoFunc) SetDefaultReturn(r0 bool) {
			f.SetDefaultHook(func(string) bool {
				return r0
			})
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncPushReturnMethod(t *testing.T) {
	code := generateMockFuncPushReturnMethod(makeMethod(TestMethodDo))
	expected := strip(`
		// PushReturn calls PushHook with a function that returns the given values.
		func (f *TestClientDoFunc) PushReturn(r0 bool) {
			f.PushHook(func(string) bool {
				return r0
			})
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncNextHookMethod(t *testing.T) {
	code := generateMockFuncNextHookMethod(makeMethod(TestMethodDo))
	expected := strip(`
		func (f *TestClientDoFunc) nextHook() func(string) bool {
			f.mutex.Lock()
			defer f.mutex.Unlock()

			if len(f.hooks) == 0 {
				return f.defaultHook
			}

			hook := f.hooks[0]
			f.hooks = f.hooks[1:]
			return hook
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncAppendCallMethod(t *testing.T) {
	code := generateMockFuncAppendCallMethod(makeMethod(TestMethodDo))
	expected := strip(`
		func (f *TestClientDoFunc) appendCall(r0 TestClientDoFuncCall) {
			f.mutex.Lock()
			f.history = append(f.history, r0)
			f.mutex.Unlock()
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}

func TestGenerateMockFuncHistoryMethod(t *testing.T) {
	code := generateMockFuncHistoryMethod(makeMethod(TestMethodDo))
	expected := strip(`
		// History returns a sequence of TestClientDoFuncCall objects describing the
		// invocations of this function.
		func (f *TestClientDoFunc) History() []TestClientDoFuncCall {
			f.mutex.Lock()
			history := make([]TestClientDoFuncCall, len(f.history))
			copy(history, f.history)
			f.mutex.Unlock()

			return history
		}
	`)
	assert.Equal(t, expected, fmt.Sprintf("%#v", code))
}
