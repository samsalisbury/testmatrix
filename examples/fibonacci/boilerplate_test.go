package fibonacci

import (
	"testing"

	"github.com/samsalisbury/testmatrix"
)

type runner struct{ *testmatrix.Runner }

type test func(*testing.T, *fixture)

func (r *runner) Run(name string, config func() *fixtureConfig, test test) {
	r.Runner.RunScenario(name, func(t *testing.T, s testmatrix.Scenario, lf *testmatrix.LateFixture) {
		fix := config().Init(t, s)
		lf.Set(fix)
		test(t, fix)
	})
}
