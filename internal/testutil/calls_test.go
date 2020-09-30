package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallHistory(t *testing.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	history, ok := GetCallHistory(m)
	assert.True(t, ok)
	assert.Len(t, history, 5)
}

func TestGetCallHistoryWith(t *testing.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := GetCallHistoryWith(m, func(v CallInstance) bool { return len(v.Args()) == 3 })
	assert.True(t, ok)
	assert.Len(t, matchingHistory, 3)
}

func TestGetCallHistoryNil(t *testing.T) {
	_, ok := GetCallHistory(nil)
	assert.False(t, ok)
}

func TestGetCallHistoryNoHistoryMethod(t *testing.T) {
	_, ok := GetCallHistory(&testNoHistory{})
	assert.False(t, ok)
}

func TestGetCallHistoryBadParamArity(t *testing.T) {
	_, ok := GetCallHistory(&testHistoryBadParamArity{})
	assert.False(t, ok)
}

func TestGetCallHistoryBadResultArity(t *testing.T) {
	_, ok := GetCallHistory(&testHistoryBadResultArity{})
	assert.False(t, ok)
}

func TestGetCallHistoryNonSliceResult(t *testing.T) {
	_, ok := GetCallHistory(&testGetCallHistoryNonSliceResult{})
	assert.False(t, ok)
}

func TestGetCallHistoryBadSliceTypes(t *testing.T) {
	_, ok := GetCallHistory(&testGetCallHistoryBadSliceTypes{})
	assert.False(t, ok)
}

//
// Helpers

type (
	testNoHistory                    struct{}
	testHistoryBadParamArity         struct{}
	testHistoryBadResultArity        struct{}
	testGetCallHistoryNonSliceResult struct{}
	testGetCallHistoryBadSliceTypes  struct{}
	halfCallInstance                 struct{}
)

func (h *testHistoryBadParamArity) History(n int) []CallInstance       { return nil }
func (h *testHistoryBadResultArity) History() ([]CallInstance, error)  { return nil, nil }
func (h *testGetCallHistoryNonSliceResult) History() string            { return "" }
func (h *testGetCallHistoryBadSliceTypes) History() []halfCallInstance { return nil }
func (h *halfCallInstance) Args() []interface{}                        { return nil }

//
//

type (
	litFunc struct {
		history []litCall
	}

	litCall struct {
		args    []interface{}
		results []interface{}
	}
)

func (f litFunc) History() []litCall     { return f.history }
func (i litCall) Args() []interface{}    { return i.args }
func (i litCall) Results() []interface{} { return i.results }
