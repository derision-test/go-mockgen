# go-mockgen

[![GoDoc](https://godoc.org/github.com/derision-test/go-mockgen?status.svg)](https://godoc.org/github.com/derision-test/go-mockgen)

A mock interface code generator.

## Generating Mocks

Install with `go get -u github.com/derision-test/go-mockgen/...`.

Mocks should be generated via `go generate` and should be regenerated on each update to the target interface. For example, in `gen.go`:

```go
package mocks

//go:generate go-mockgen -f github.com/example/package -i ExampleInterface -o mock_example_interface_test.go
```

Depending on how you prefer to structure your code, you can either

1. generate mocks next to the implementation (as a sibling or in a sibling `mocks` package), or
2. generate mocks as needed in test code (generating them into a `_test.go` file).

### Flags

The following flags are defined by the binary.

| Name       | Short Flag | Description  |
| ---------- | ---------- | ------------ |
| package    | p          | The name of the generated package. Is the name of target directory if dirname or filename is supplied by default. |
| prefix     |            | A prefix used in the name of each mock struct. Should be TitleCase by convention. |
| interfaces | i          | A whitelist of interfaces to generate given the import paths. |
| filename   | o          | The target output file. All mocks are writen to this file. |
| dirname    | d          | The target output directory. Each mock will be written to a unique file. |
| force      | f          | Do not abort if a write to disk would overwrite an existing file. |

If neither dirname nor filename are supplied, then the generated code is printed to standard out.

## Testing with Mocks

A mock value fulfills all of the methods of the target interface from which it was generated. Unless overridden, all methods of the mock will return zero values for everything. To override a specific method, you can set its `hook` or its `return values`.

A hook is a method that is called on each invocation and allows the test to specify complex behaviors in the mocked interface (conditionally returning values, synchronizing on external state, etc,). The default hook for a method is set with the `SetDefaultHook` method.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache()
    cache.GetFunc.SetDefaultHook(func (key string) (interface{}, bool) {
        if key == "expected" {
            return 42, true
        }
        return nil, false
    })

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

In the cases where you don't need specific behaviors but just need to return some data, the setup gets a bit easier with `SetDefaultReturn`.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache()
    cache.GetFunc.SetDefaultReturn(42, true)

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

Hook and return values can also be *stacked* when your test can anticipate multiple calls to the same function. Pushing a hook or a return value will set the hook or return value for *one* invocation of the mocked method. Once this hook or return value has been spent, it will be removed from the queue. Hooks and return values can be interleaved. If the queue is empty, the default hook will be invoked (or the default return values returned).

The following example will test a cache that returns values 50, 51, and 52 in sequence, then panic if there is an unexpected fourth call.

```go
func TestCache(t *testing.T) {
    cache := mocks.NewMockCache()
    cache.GetFunc.SetDefaultHook(func (key string) (interface{}, bool) {
        panic("unexpected call")
    })
    cache.GetFunc.PushReturn(50, true)
    cache.GetFunc.PushReturn(51, true)
    cache.GetFunc.PushReturn(52, true)

    testSubject := NewThingThatNeedsCache(cache)
    // ...
}
```

### Assertions

Mocks track their invocations and can be retrieved via the `History` method. Structs are generated for each method type containing fields for each argument and result type. Raw assertions can be performed on these values.

```go
allCalls := cache.GetFunc.History()
allCalls[0].Arg0 // key
allCalls[0].Result0 // value
allCalls[0].Result1 // exists flag
```

This library also contains a set of [Gomega](https://onsi.github.io/gomega/) matchers which simplify assertions over a mocked method's call history.

To use the matchers, import the matchers package anonymously.

```go
import . "github.com/derision-test/go-mockgen/matchers"
```

The following matchers are defined.

- `BeCalled()`
- `BeCalledN(n)`
- `BeCalledOnce()`
- `BeCalledWith(args...)`
- `BeCalledOnceWith(args...)`
- `BeAnything()`

These matchers can be used as follows.

```go
// cache.Get called 3 times
Expect(cache.GetFunc).To(BeCalledN(3)) 

// Ensure cache.Set("foo", "bar") was called
Expect(cache.SetFunc).To(BeCalledWith("foo", "bar")) 

// Ensure cache.Set("foo", _) was called
Expect(cache.SetFunc).To(BeCalledWith("foo", BeAnything())) 
```

## License

Copyright (c) 2020 Eric Fritz

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
