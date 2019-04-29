package main

import (
	"fmt"
	gotypes "go/types"
	"strings"

	"github.com/aphistic/sweet"
	"github.com/dave/jennifer/jen"
	"github.com/efritz/go-genlib/types"
	. "github.com/onsi/gomega"
)

type GenerationSuite struct{}

const (
	TestPrefix         = "Test"
	TestTitleName      = "Client"
	TestMockStructName = "MockTestClient"
	TestImportPath     = "github.com/efritz/go-mockgen/test"
)

var (
	boolType        = getType(gotypes.Bool)
	stringType      = getType(gotypes.String)
	stringSliceType = gotypes.NewSlice(getType(gotypes.String))

	TestMethodStatus = &types.Method{
		Name:    "Status",
		Params:  []gotypes.Type{},
		Results: []gotypes.Type{stringType, boolType},
	}

	TestMethodDo = &types.Method{
		Name:    "Do",
		Params:  []gotypes.Type{stringType},
		Results: []gotypes.Type{boolType},
	}

	TestMethodDof = &types.Method{
		Name:     "Dof",
		Params:   []gotypes.Type{stringType, stringSliceType},
		Results:  []gotypes.Type{boolType},
		Variadic: true,
	}
)

func (s *GenerationSuite) TestGenerateInterface(t sweet.T) {
	expectedDecls := []string{
		// Structs
		"type MockTestClient struct",
		"type TestClientDoFunc struct",
		"type TestClientDoFuncCall struct",
		"type TestClientDofFunc struct",
		"type TestClientDofFuncCall struct",
		"func NewMockTestClient() *MockTestClient",
		// Overrides
		"func (m *MockTestClient) Do(v0 string) bool",
		"func (m *MockTestClient) Dof(v0 string, v1 ...string) bool",
		// DoFunc Methods
		"func (f *TestClientDoFunc) SetDefaultHook(hook func(string) bool)",
		"func (f *TestClientDoFunc) PushHook(hook func(string) bool)",
		"func (f *TestClientDoFunc) SetDefaultReturn(r0 bool)",
		"func (f *TestClientDoFunc) PushReturn(r0 bool)",
		"func (f *TestClientDoFunc) History() []TestClientDoFuncCall",
		// DoFuncCall methods
		"func (c TestClientDoFuncCall) Args() []interface{}",
		"func (c TestClientDoFuncCall) Results() []interface{}",
		// DofFunc Methods
		"func (f *TestClientDofFunc) SetDefaultHook(hook func(string, ...string) bool)",
		"func (f *TestClientDofFunc) PushHook(hook func(string, ...string) bool)",
		"func (f *TestClientDofFunc) SetDefaultReturn(r0 bool)",
		"func (f *TestClientDofFunc) PushReturn(r0 bool)",
		"func (f *TestClientDofFunc) History() []TestClientDofFuncCall",
		// DofFuncCall methods
		"func (c TestClientDofFuncCall) Args() []interface{}",
		"func (c TestClientDofFuncCall) Results() []interface{}",
	}

	file := jen.NewFile("test")
	g := &generator{""}
	g.generateInterface(file, makeBareInterface(TestMethodDo, TestMethodDof), TestPrefix)
	rendered := fmt.Sprintf("%#v\n", file)

	for _, decl := range expectedDecls {
		Expect(rendered).To(ContainSubstring(decl))
	}
}

func (s *GenerationSuite) TestGenerateMockStruct(t sweet.T) {
	g := &generator{""}
	code := g.generateMockStruct(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// MockTestClient is a mock impelementation of the Client interface (from
	// the package github.com/efritz/go-mockgen/test) used for unit testing.
	type MockTestClient struct {
		// StatusFunc is an instance of a mock function object controlling the
		// behavior of the method Status.
		StatusFunc *TestClientStatusFunc
		// DoFunc is an instance of a mock function object controlling the
		// behavior of the method Do.
		DoFunc *TestClientDoFunc
		// DofFunc is an instance of a mock function object controlling the
		// behavior of the method Dof.
		DofFunc *TestClientDofFunc
	}
	`)))
}

func (s *GenerationSuite) TestGenerateMockStructConstructor(t sweet.T) {
	g := &generator{""}
	code := g.generateMockStructConstructor(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// NewMockTestClient creates a new mock of the Client interface. All methods
	// return zero values for all results, unless overwritten.
	func NewMockTestClient() *MockTestClient {
		return &MockTestClient{
			StatusFunc: &TestClientStatusFunc{
				defaultHook: func() (string, bool) {
					return "", false
				},
			},
			DoFunc: &TestClientDoFunc{
				defaultHook: func(string) bool {
					return false
				},
			},
			DofFunc: &TestClientDofFunc{
				defaultHook: func(string, ...string) bool {
					return false
				},
			},
		}
	}
	`)))
}

func (s *GenerationSuite) TestGenerateMockStructFromConstructor(t sweet.T) {
	g := &generator{""}
	code := g.generateMockStructFromConstructor(makeInterface(TestMethodStatus, TestMethodDo, TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// NewMockTestClientFrom creates a new mock of the MockTestClient interface.
	// All methods delegate to the given implementation, unless overwritten.
	func NewMockTestClientFrom(i test.Client) *MockTestClient {
		return &MockTestClient{
			StatusFunc: &TestClientStatusFunc{
				defaultHook: i.Status,
			},
			DoFunc: &TestClientDoFunc{
				defaultHook: i.Do,
			},
			DofFunc: &TestClientDofFunc{
				defaultHook: i.Dof,
			},
		}
	}
	`)))
}

func (s *GenerationSuite) TestGenerateMockStructFromConstructorUnexported(t sweet.T) {
	iface := makeBareInterface(TestMethodStatus, TestMethodDo, TestMethodDof)
	iface.Name = "client"

	g := &generator{""}
	code := g.generateMockStructFromConstructor(g.wrapInterface(iface, TestPrefix, TestTitleName, TestMockStructName))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// surrogateMockClient is a copy of the client interface (from the package
	// github.com/efritz/go-mockgen/test). It is redefined here as it is
	// unexported in the source packge.
	type surrogateMockClient interface {
		Status() (string, bool)
		Do(string) bool
		Dof(string, ...string) bool
	}

	// NewMockTestClientFrom creates a new mock of the MockTestClient interface.
	// All methods delegate to the given implementation, unless overwritten.
	func NewMockTestClientFrom(i surrogateMockClient) *MockTestClient {
		return &MockTestClient{
			StatusFunc: &TestClientStatusFunc{
				defaultHook: i.Status,
			},
			DoFunc: &TestClientDoFunc{
				defaultHook: i.Do,
			},
			DofFunc: &TestClientDofFunc{
				defaultHook: i.Dof,
			},
		}
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncStruct(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncStruct(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// TestClientDoFunc describes the behavior when the Do method of the parent
	// MockTestClient instance is invoked.
	type TestClientDoFunc struct {
		defaultHook func(string) bool
		hooks       []func(string) bool
		history     []TestClientDoFuncCall
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncStructVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncStruct(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// TestClientDofFunc describes the behavior when the Dof method of the
	// parent MockTestClient instance is invoked.
	type TestClientDofFunc struct {
		defaultHook func(string, ...string) bool
		hooks       []func(string, ...string) bool
		history     []TestClientDofFuncCall
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFunc(t sweet.T) {
	g := &generator{""}
	code := g.generateFunc(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Do delegates to the next hook function in the queue and stores the
	// parameter and result values of this invocation.
	func (m *MockTestClient) Do(v0 string) bool {
		r0 := m.DoFunc.nextHook()(v0)
		m.DoFunc.history = append(m.DoFunc.history, TestClientDoFuncCall{v0, r0})
		return r0
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateFunc(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Dof delegates to the next hook function in the queue and stores the
	// parameter and result values of this invocation.
	func (m *MockTestClient) Dof(v0 string, v1 ...string) bool {
		r0 := m.DofFunc.nextHook()(v0, v1...)
		m.DofFunc.history = append(m.DofFunc.history, TestClientDofFuncCall{v0, v1, r0})
		return r0
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncSetHookMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncSetHookMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// SetDefaultHook sets function that is called when the Do method of the
	// parent MockTestClient instance is invoked and the hook queue is empty.
	func (f *TestClientDoFunc) SetDefaultHook(hook func(string) bool) {
		f.defaultHook = hook
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncSetHookMethodVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncSetHookMethod(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// SetDefaultHook sets function that is called when the Dof method of the
	// parent MockTestClient instance is invoked and the hook queue is empty.
	func (f *TestClientDofFunc) SetDefaultHook(hook func(string, ...string) bool) {
		f.defaultHook = hook
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncPushHookMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncPushHookMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// PushHook adds a function to the end of hook queue. Each invocation of the
	// Do method of the parent MockTestClient instance inovkes the hook at the
	// front of the queue and discards it. After the queue is empty, the default
	// hook function is invoked for any future action.
	func (f *TestClientDoFunc) PushHook(hook func(string) bool) {
		f.hooks = append(f.hooks, hook)
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncPushHookMethodVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncPushHookMethod(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// PushHook adds a function to the end of hook queue. Each invocation of the
	// Dof method of the parent MockTestClient instance inovkes the hook at the
	// front of the queue and discards it. After the queue is empty, the default
	// hook function is invoked for any future action.
	func (f *TestClientDofFunc) PushHook(hook func(string, ...string) bool) {
		f.hooks = append(f.hooks, hook)
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncSetReturnMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncSetReturnMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
	// the given values.
	func (f *TestClientDoFunc) SetDefaultReturn(r0 bool) {
		f.SetDefaultHook(func(string) bool {
			return r0
		})
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncPushReturnMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncPushReturnMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// PushReturn calls PushDefaultHook with a function that returns the given
	// values.
	func (f *TestClientDoFunc) PushReturn(r0 bool) {
		f.PushHook(func(string) bool {
			return r0
		})
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncNextHookMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncNextHookMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	func (f *TestClientDoFunc) nextHook() func(string) bool {
		if len(f.hooks) == 0 {
			return f.defaultHook
		}

		hook := f.hooks[0]
		f.hooks = f.hooks[1:]
		return hook
	}
	`)))
}

func (s *GenerationSuite) TestGenerateFuncHistoryMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateFuncHistoryMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// History returns a sequence of TestClientDoFuncCall objects describing the
	// invocations of this function.
	func (f *TestClientDoFunc) History() []TestClientDoFuncCall {
		return f.history
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallStruct(t sweet.T) {
	g := &generator{""}
	code := g.generateCallStruct(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// TestClientDoFuncCall is an object that describes an invocation of method
	// Do on an instance of MockTestClient.
	type TestClientDoFuncCall struct {
		// Arg0 is the value of the 1st argument passed to this method
		// invocation.
		Arg0 string
		// Result0 is the value of the 1st result returned from this method
		// invocation.
		Result0 bool
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallStructVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateCallStruct(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// TestClientDofFuncCall is an object that describes an invocation of method
	// Dof on an instance of MockTestClient.
	type TestClientDofFuncCall struct {
		// Arg0 is the value of the 1st argument passed to this method
		// invocation.
		Arg0 string
		// Arg1 is a slice containing the values of the variadic arguments
		// passed to this method invocation.
		Arg1 []string
		// Result0 is the value of the 1st result returned from this method
		// invocation.
		Result0 bool
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallArgMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateCallArgsMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Args returns an interface slice containing the arguments of this
	// invocation.
	func (c TestClientDoFuncCall) Args() []interface{} {
		return []interface{}{c.Arg0}
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallArgsMethodVariadic(t sweet.T) {
	g := &generator{""}
	code := g.generateCallArgsMethod(makeMethod(TestMethodDof))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Args returns an interface slice containing the arguments of this
	// invocation. The variadic slice argument is flattened in this array such
	// that one positional argument and three variadic arguments would result in
	// a slice of four, not two.
	func (c TestClientDofFuncCall) Args() []interface{} {
		trailing := []interface{}{}
		for _, val := range c.Arg1 {
			trailing = append(trailing, val)
		}

		return append([]interface{}{c.Arg0}, trailing...)
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallResultsMethod(t sweet.T) {
	g := &generator{""}
	code := g.generateCallResultsMethod(makeMethod(TestMethodDo))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Results returns an interface slice containing the results of this
	// invocation.
	func (c TestClientDoFuncCall) Results() []interface{} {
		return []interface{}{c.Result0}
	}
	`)))
}

func (s *GenerationSuite) TestGenerateCallResultsMethodMultiple(t sweet.T) {
	g := &generator{""}
	code := g.generateCallResultsMethod(makeMethod(TestMethodStatus))

	Expect(fmt.Sprintf("%#v", code)).To(Equal(strip(`
	// Results returns an interface slice containing the results of this
	// invocation.
	func (c TestClientStatusFuncCall) Results() []interface{} {
		return []interface{}{c.Result0, c.Result1}
	}
	`)))
}

//
// Helpers

func getType(kind gotypes.BasicKind) gotypes.Type {
	return gotypes.Typ[kind].Underlying()
}

func makeBareInterface(methods ...*types.Method) *types.Interface {
	return &types.Interface{
		Name:       TestTitleName,
		ImportPath: TestImportPath,
		Type:       types.InterfaceTypeInterface,
		Methods:    methods,
	}
}

func makeInterface(methods ...*types.Method) *wrappedInterface {
	g := &generator{""}
	return g.wrapInterface(makeBareInterface(methods...), TestPrefix, TestTitleName, TestMockStructName)
}

func makeMethod(methods ...*types.Method) (*wrappedInterface, *wrappedMethod) {
	wrapped := makeInterface(methods...)
	return wrapped, wrapped.wrappedMethods[0]
}

func strip(block string) string {
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if len(line) > 0 && line[0] == '\t' {
			lines[i] = line[1:]
		}
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
