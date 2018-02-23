# go-mockgen

A mock interface code generator.

## Installation

Simply run `go get -u github.com/efritz/go-mockgen/...`.

### Usage

As an example, we generate a mock implementation for the Retry interface in
the [watchdog](https://github.com/efritz/watchdog) library. After running
the command ```go-mockgen github.com/efritz/watchdog Retry```, the following
code is generated and printed to standard out.

```
package test

import watchdog "github.com/efritz/watchdog"

type MockRetry struct {
	RetryFunc func() bool
}

var _ watchdog.Retry = NewMockRetry()

func NewMockRetry() *MockRetry {
	return &MockRetry{RetryFunc: func() bool {
		return false
	}}
}
func (m *MockRetry) Retry() bool {
	return m.RetryFunc()
}
```

If no interface (or list of interfaces) are given after the import name, a mock
for every exported interface defined in that package is mocked.

## License

Copyright (c) 2018 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
