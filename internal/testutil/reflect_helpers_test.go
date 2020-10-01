package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallHistory(t *testing.T) {
	value := newHistory(
		mockCall{args: []interface{}{"foo", "bar"}},
		mockCall{args: []interface{}{"foo", "bar", "baz"}},
		mockCall{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		mockCall{args: []interface{}{"foo", "bar", "bonk"}},
		mockCall{args: []interface{}{"foo", "bar", "baz"}},
	)

	history, ok := GetCallHistory(value)
	assert.True(t, ok)
	assert.Len(t, history, 5)
}

func TestGetCallHistoryWith(t *testing.T) {
	value := newHistory(
		mockCall{args: []interface{}{"foo", "bar"}},
		mockCall{args: []interface{}{"foo", "bar", "baz"}},
		mockCall{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		mockCall{args: []interface{}{"foo", "bar", "bonk"}},
		mockCall{args: []interface{}{"foo", "bar", "baz"}},
	)

	matchingHistory, ok := GetCallHistoryWith(value, func(v CallInstance) bool { return len(v.Args()) == 3 })
	assert.True(t, ok)
	assert.Len(t, matchingHistory, 3)
}

func TestGetCallHistoryNil(t *testing.T) {
	_, ok := GetCallHistory(nil)
	assert.False(t, ok)
}

func TestGetCallHistoryNoHistoryMethod(t *testing.T) {
	_, ok := GetCallHistory(struct{}{})
	assert.False(t, ok)
}

func TestGetCallHistoryBadParamArity(t *testing.T) {
	_, ok := GetCallHistory(&historyFuncBadParamArity{})
	assert.False(t, ok)
}

type historyFuncBadParamArity struct{}

func (h *historyFuncBadParamArity) History(n int) []CallInstance {
	return nil
}

func TestGetCallHistoryBadResultArity(t *testing.T) {
	_, ok := GetCallHistory(&historyFuncBadResultArity{})
	assert.False(t, ok)
}

type historyFuncBadResultArity struct{}

func (h *historyFuncBadResultArity) History() ([]CallInstance, error) {
	return nil, nil
}

func TestGetCallHistoryNonSliceResult(t *testing.T) {
	_, ok := GetCallHistory(&historyFuncNonSliceResult{})
	assert.False(t, ok)
}

type historyFuncNonSliceResult struct{}

func (h *historyFuncNonSliceResult) History() string {
	return ""
}

func TestGetCallHistoryBadSliceTypes(t *testing.T) {
	_, ok := GetCallHistory(&historyFuncBadSliceTypes{})
	assert.False(t, ok)
}

type historyFuncBadSliceTypes struct{}

func (h *historyFuncBadSliceTypes) History() []string {
	return nil
}

func TestGetArgs(t *testing.T) {
	expectedArgs := []interface{}{"foo", "bar", "baz", "bonk"}
	args, ok := GetArgs(mockCall{args: expectedArgs})
	assert.True(t, ok)
	assert.Equal(t, expectedArgs, args)
}

func TestGetArgsNil(t *testing.T) {
	_, ok := GetArgs(nil)
	assert.False(t, ok)
}

func TestGetArgsNoArgsMethod(t *testing.T) {
	_, ok := GetArgs(struct{}{})
	assert.False(t, ok)
}

func TestGetArgsBadParamArity(t *testing.T) {
	_, ok := GetArgs(&argsFuncBadParamArity{})
	assert.False(t, ok)
}

type argsFuncBadParamArity struct{}

func (m *argsFuncBadParamArity) Args(n int) []interface{} {
	return nil
}

func TestGetArgsBadResultArity(t *testing.T) {
	_, ok := GetArgs(&argsFuncBadResultArity{})
	assert.False(t, ok)
}

type argsFuncBadResultArity struct{}

func (m *argsFuncBadResultArity) Args() ([]interface{}, error) {
	return nil, nil
}

func TestGetArgsNonSliceResult(t *testing.T) {
	_, ok := GetArgs(&argsFuncNonSliceResult{})
	assert.False(t, ok)
}

type argsFuncNonSliceResult struct{}

func (m *argsFuncNonSliceResult) Args() string {
	return ""
}

func TestGetArgsBadSliceTypes(t *testing.T) {
	_, ok := GetArgs(&argsFuncBadSliceTypes{})
	assert.False(t, ok)
}

type argsFuncBadSliceTypes struct{}

func (m *argsFuncBadSliceTypes) Args() []string {
	return nil
}
