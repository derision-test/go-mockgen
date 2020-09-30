package matchers

import (
	"testing"

	"github.com/aphistic/sweet"
	junit "github.com/aphistic/sweet-junit"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	RegisterFailHandler(sweet.GomegaFail)

	sweet.Run(m, func(s *sweet.S) {
		s.RegisterPlugin(junit.NewPlugin())

		s.AddSuite(&AnythingMatcherSuite{})
		s.AddSuite(&CalledMatcherSuite{})
		s.AddSuite(&CalledNMatcherSuite{})
		s.AddSuite(&CalledOnceMatcherSuite{})
		s.AddSuite(&CalledOnceWithMatcherSuite{})
		s.AddSuite(&CalledWithMatcherSuite{})
		s.AddSuite(&CallsSuite{})
	})
}

//
//

type (
	litFunc struct {
		history []litCall
	}

	litCall struct {
		args    []interface{}
		results []interface{}
	}
)

func (f litFunc) History() []litCall     { return f.history }
func (i litCall) Args() []interface{}    { return i.args }
func (i litCall) Results() []interface{} { return i.results }
