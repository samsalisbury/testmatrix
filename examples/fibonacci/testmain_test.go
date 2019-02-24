package fibonacci

import (
	"os"
	"testing"

	"github.com/samsalisbury/testmatrix"
)

// TestMain tells us to run tests as defined in the matrix.
// Note: if you need more control over TestMain, you can manually call
// testmatrix.Init and testmatrix.PrintSummary instead of testmatrix.Run.
// See the implementation of testmatrix.Run for details on how to do this.
func TestMain(m *testing.M) {
	os.Exit(testmatrix.Run(m, makeMatrix))
}

// makeMatrix() returns the full test matrix. If certain matrix combinations
// should not be tested, it is possible to make such derivations using the
// Matrix.FixedDimension method, and more like it. Do this in top-level tests.
//
// The matrix func should always return a fresh matrix, so that references
// are not shared between tests.
func makeMatrix() testmatrix.Matrix {
	return testmatrix.New(
		testmatrix.Dimension{
			Name: "fib",
			Desc: "fibbonaci func",
			Values: testmatrix.Values{
				"recur": &Recursive{},
				"iter":  &Iterative{},
			},
		},
		testmatrix.Dimension{
			Name: "enhancement",
			Desc: "which enhancement to use",
			Values: testmatrix.Values{
				"plain":    Enhancer(func(p Provider) Provider { return p }),
				"memoized": Enhancer(func(p Provider) Provider { return NewMemoized(p) }),
			},
		},
	)
}

// fixture represents a fully realised fixture, generated from the injected
// scenario.
type fixture struct {
	Provider
}

// fixtureFunc returns a strongly-typed *fixture.
type fixtureFunc func(*testing.T, testmatrix.Scenario) *fixture

// makeFixture is a fixtureFunc that creates a fixture from the give *testing.T
// and scenario.
func makeFixture(t *testing.T, s testmatrix.Scenario) *fixture {
	provider := s.Value("fib").(Provider)
	enhanced := s.Value("enhancement").(Enhancer)
	return &fixture{
		Provider: enhanced(provider),
	}
}

// testFunc is a test function that takes a stongly-typed *fixture.
type testFunc func(*testing.T, *fixture)

// runner wraps *testmatrix.Runner and adds a strongly-typed Run func.
type runner struct{ *testmatrix.Runner }

// Run accepts your strongly typed fixtureFunc and testFunc, wraps them up and
// passes them through to the generic testmatrix.Runner.Run for execution.
func (r *runner) Run(name string, makeFixture fixtureFunc, test testFunc) {
	r.Runner.Run(name,
		// Return a strongly typed fixture.
		func(t *testing.T, s testmatrix.Scenario) testmatrix.Fixture {
			return makeFixture(t, s)
		},
		// Unwrap strongly typed feature, pass to strongly typed test.
		func(t *testing.T, f testmatrix.Fixture) {
			test(t, f.(*fixture))
		},
	)
}

func newRunner(t *testing.T) *runner {
	return &runner{testmatrix.NewRunner(t)}
}
