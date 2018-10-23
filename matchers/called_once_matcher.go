package matchers

import (
	"github.com/onsi/gomega/types"
)

func BeCalledOnce() types.GomegaMatcher {
	return BeCalledN(1)
}
