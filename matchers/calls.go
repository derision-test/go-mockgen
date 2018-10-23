package matchers

import (
	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

func getMatchingCallCounts(history []callInstance, args []interface{}) (n int, lastErr error) {
outer:
	for _, call := range history {
		if len(call.Args()) != len(args) {
			continue
		}

		for i, arg := range call.Args() {
			matcher, ok := args[i].(types.GomegaMatcher)
			if !ok {
				matcher = &matchers.EqualMatcher{Expected: args[i]}
			}

			success, err := matcher.Match(arg)
			if err != nil {
				lastErr = err
				continue outer
			}

			if !success {
				continue outer
			}
		}

		n++
	}

	return
}
