package integrationgcexportdata

import mockrequire "github.com/derision-test/go-mockgen/testutil/require"

type Banana interface {
	DoSomething() mockrequire.CallInstanceAsserter
}
