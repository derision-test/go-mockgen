package mockgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	assert.Equal(t, "", title(""))
	assert.Equal(t, "Foobar", title("foobar"))
	assert.Equal(t, "FooBar", title("fooBar"))
}
