// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ ./go-mockgen github.com/efritz/watchdog -p foo -d tmp -f

package foo

import watchdog "github.com/efritz/watchdog"

type MockRetry struct {
	RetryFunc func() bool
}

var _ watchdog.Retry = NewMockRetry()

func NewMockRetry() *MockRetry {
	m := &MockRetry{}
	m.RetryFunc = m.defaultRetryFunc
	return m
}
func (m *MockRetry) Retry() bool {
	return m.RetryFunc()
}
func (m *MockRetry) defaultRetryFunc() bool {
	return false
}
