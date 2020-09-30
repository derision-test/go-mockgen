package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalledMatch(t *testing.T) {
	ok, err := BeCalled().Match(litFunc{[]litCall{{}}})
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCalledMatchEmptyHistory(t *testing.T) {
	ok, err := BeCalled().Match(litFunc{})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledMatchError(t *testing.T) {
	_, err := BeCalled().Match(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BeCalled expects a mock function")
}

func TestCalledNMatch(t *testing.T) {
	ok, err := BeCalledN(2).Match(litFunc{[]litCall{{}, {}}})
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCalledNMatchEmptyHistory(t *testing.T) {
	ok, err := BeCalledN(1).Match(litFunc{})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledNMatchMismatchedHistory(t *testing.T) {
	ok, err := BeCalledN(1).Match(litFunc{[]litCall{{}, {}}})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledNMatchError(t *testing.T) {
	_, err := BeCalledN(1).Match(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BeCalledN expects a mock function")
}

func TestCalledOnceMatch(t *testing.T) {
	ok, err := BeCalledOnce().Match(litFunc{[]litCall{{}}})
	assert.Nil(t, err)
	assert.True(t, ok)
}

func TestCalledOnceMatchEmptyHistory(t *testing.T) {
	ok, err := BeCalledOnce().Match(litFunc{})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledOnceMatchMismatchedHistory(t *testing.T) {
	ok, err := BeCalledOnce().Match(litFunc{[]litCall{{}, {}}})
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestCalledOnceMatchError(t *testing.T) {
	_, err := BeCalledOnce().Match(nil)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "BeCalledOnce expects a mock function")
}
