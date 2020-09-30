package mockassert

import (
	"fmt"

	"github.com/derision-test/go-mockgen/internal/testutil"
	"github.com/stretchr/testify/assert"
)

type CallInstanceAssertionFunc func(assert.TestingT, interface{}) bool

// Called asserts that the mock function object was called at least once.
func Called(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) bool {
	callCount, ok := callCount(t, mockFn, msgAndArgs...)
	if !ok {
		return false
	}
	if callCount == 0 {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called at least once", mockFn), msgAndArgs...)
	}

	return true
}

// NotCalled asserts that the mock function object was not called.
func NotCalled(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) bool {
	callCount, ok := callCount(t, mockFn, msgAndArgs...)
	if !ok {
		return false
	}
	if callCount != 0 {
		return assert.Fail(t, fmt.Sprintf("Did not expect %T to be called", mockFn), msgAndArgs...)
	}

	return true
}

// CalledOnce asserts that the mock function object was called exactly once.
func CalledOnce(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) bool {
	return CalledN(t, mockFn, 1, msgAndArgs...)
}

// CalledOnce asserts that the mock function object was called exactly n times.
func CalledN(t assert.TestingT, mockFn interface{}, n int, msgAndArgs ...interface{}) bool {
	callCount, ok := callCount(t, mockFn, msgAndArgs...)
	if !ok {
		return false
	}
	if callCount != n {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called exactly %d times, called %d times", mockFn, n, callCount), msgAndArgs...)
	}

	return true
}

// CalledWith asserts that the mock function object was called at least once with a set of
// arguments matching the given assertion function.
func CalledWith(t assert.TestingT, mockFn interface{}, assertion CallInstanceAssertionFunc, msgAndArgs ...interface{}) bool {
	matchingCallCount, ok := callCountWith(t, mockFn, assertion, msgAndArgs...)
	if !ok {
		return false
	}
	if matchingCallCount == 0 {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called at least once", mockFn), msgAndArgs...)
	}
	return true
}

// CalledWith asserts that the mock function object was not called with a set of arguments
// matching the given assertion function.
func NotCalledWith(t assert.TestingT, mockFn interface{}, assertion CallInstanceAssertionFunc, msgAndArgs ...interface{}) bool {
	matchingCallCount, ok := callCountWith(t, mockFn, assertion, msgAndArgs...)
	if !ok {
		return false
	}
	if matchingCallCount != 0 {
		return assert.Fail(t, fmt.Sprintf("Did not expect %T to be called", mockFn), msgAndArgs...)
	}
	return true
}

// CalledOnceWith asserts that the mock function object was called exactly once with a set of
// arguments matching the given assertion function.
func CalledOnceWith(t assert.TestingT, mockFn interface{}, assertion CallInstanceAssertionFunc, msgAndArgs ...interface{}) bool {
	return CalledNWith(t, mockFn, 1, assertion, msgAndArgs...)
}

// CalledNWith asserts that the mock function object was called exactly n times with a set of
// arguments matching the given assertion function.
func CalledNWith(t assert.TestingT, mockFn interface{}, n int, assertion CallInstanceAssertionFunc, msgAndArgs ...interface{}) bool {
	matchingCallCount, ok := callCountWith(t, mockFn, assertion, msgAndArgs...)
	if !ok {
		return false
	}
	if matchingCallCount != n {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called exactly %d times, called %d times", mockFn, n, matchingCallCount), msgAndArgs...)
	}
	return true
}

// callCount returns the number of times the given mock function was called.
func callCount(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) (int, bool) {
	return callCountWith(t, mockFn, func(t assert.TestingT, call interface{}) bool { return true }, msgAndArgs...)
}

// callCount returns the number of times the given mock function was called with a set of
// arguments matching the given assertion function.
func callCountWith(t assert.TestingT, mockFn interface{}, assertion CallInstanceAssertionFunc, msgAndArgs ...interface{}) (int, bool) {
	matchingHistory, ok := testutil.GetCallHistoryWith(mockFn, func(call testutil.CallInstance) bool {
		// Pass in a dummy non-erroring TestingT so that any assertions done inside
		// this function will not fail the enclosing test.
		return assertion(mockTestingT{}, call)
	})
	if !ok {
		return 0, assert.Fail(t, fmt.Sprintf("Parameters must be a mock function description, got %T", mockFn), msgAndArgs...)
	}

	return len(matchingHistory), true
}

type mockTestingT struct{}

func (mockTestingT) Errorf(format string, args ...interface{}) {}
