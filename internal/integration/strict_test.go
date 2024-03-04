package integration

import (
	"context"
	"testing"

	"github.com/derision-test/go-mockgen/v2/internal/integration/testdata/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStrictConstructor(t *testing.T) {
	// All invocations panic by default
	mock := mocks.NewStrictMockRetrier()

	assert.Panics(t, func() {
		_ = mock.Retry(context.Background(), func() error {
			return nil
		})
	})

	// Should not panic if overwritten
	mock.RetryFunc.SetDefaultReturn(nil)

	assert.NotPanics(t, func() {
		_ = mock.Retry(context.Background(), func() error {
			return nil
		})
	})
}
