package matchers

import (
	"testing"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestCalledWithMatch(t *testing.T) {
	ok, err := BeCalledWith(gomega.ContainSubstring("foo"), 1, gomega.Not(gomega.Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"barbaz", 1, 2, 3}, nil},
	}})

	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCalledWithNoMatch(t *testing.T) {
	ok, err := BeCalledWith(gomega.ContainSubstring("foo"), 1, gomega.Not(gomega.Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"barbaz", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledWithMatchError(t *testing.T) {
	_, err := BeCalledWith("foo", 1, 2, 3).Match(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BeCalledWith expects a mock function")
}

func TestCalledOnceWithMatch(t *testing.T) {
	ok, err := BeCalledOnceWith(gomega.ContainSubstring("foo"), 1, gomega.Not(gomega.Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"barbaz", 1, 2, 3}, nil},
	}})

	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCalledOnceWithNoMatch(t *testing.T) {
	ok, err := BeCalledOnceWith(gomega.ContainSubstring("foo"), 1, gomega.Not(gomega.Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"barbaz", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledOnceWithMultipleMatches(t *testing.T) {
	ok, err := BeCalledOnceWith(gomega.ContainSubstring("foo"), 1, gomega.Not(gomega.Equal(2)), 3).Match(litFunc{[]litCall{
		{[]interface{}{"foobar", 1, 2, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"foobar", 1, 4, 3}, nil},
		{[]interface{}{"foobaz", 1, 2, 3}, nil},
	}})

	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledOnceWithMatchError(t *testing.T) {
	_, err := BeCalledOnceWith("foo", 1, 2, 3).Match(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BeCalledOnceWith expects a mock function")
}

func TestGetMatchingCallCountsLiterals(t *testing.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, "foo", "bar", "baz")
	assert.True(t, ok)
	assert.Len(t, matchingHistory, 3)
}

func TestGetMatchingCallCountsMatchers(t *testing.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, gomega.HaveLen(3), gomega.HaveLen(3), gomega.HaveLen(3))
	assert.True(t, ok)
	assert.Len(t, matchingHistory, 3)
}

func TestGetMatchingCallCountsMixed(t *testing.T) {
	m := litFunc{[]litCall{
		{args: []interface{}{"foo", "bar"}},
		{args: []interface{}{"foo", "bar", "baz"}},
		{args: []interface{}{"foo", "bar", "baz", "bonk"}},
		{args: []interface{}{"foo", "bar", "bonk"}},
		{args: []interface{}{"foo", "bar", "baz"}},
	}}

	matchingHistory, ok := getCallHistoryWith(m, "foo", "bar", gomega.ContainSubstring("bo"))
	assert.True(t, ok)
	assert.Len(t, matchingHistory, 1)
}
