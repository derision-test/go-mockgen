// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ ./go-mockgen github.com/efritz/watchdog -p foo -d tmp -f

package foo

import watchdog "github.com/efritz/watchdog"

type MockWatcher struct {
	CheckFunc func()
	StartFunc func() <-chan struct{}
	StopFunc  func()
}

var _ watchdog.Watcher = NewMockWatcher()

func NewMockWatcher() *MockWatcher {
	m := &MockWatcher{}
	m.CheckFunc = m.defaultCheckFunc
	m.StartFunc = m.defaultStartFunc
	m.StopFunc = m.defaultStopFunc
	return m
}
func (m *MockWatcher) Start() <-chan struct{} {
	return m.StartFunc()
}
func (m *MockWatcher) Stop() {
	m.StopFunc()
}
func (m *MockWatcher) Check() {
	m.CheckFunc()
}
func (m *MockWatcher) defaultCheckFunc() {
	return
}
func (m *MockWatcher) defaultStartFunc() <-chan struct{} {
	return nil
}
func (m *MockWatcher) defaultStopFunc() {
	return
}
