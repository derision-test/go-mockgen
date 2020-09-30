package matchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnythingMatch(t *testing.T) {
	ok, err := BeAnything().Match(nil)
	assert.Nil(t, err)
	assert.True(t, ok)
}
