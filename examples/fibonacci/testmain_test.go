package fibonacci

import (
	"os"
	"testing"

	"github.com/samsalisbury/testmatrix"
)

// matrix is our matrix definition.
var matrix = testmatrix.New(
	testmatrix.Dim("fib", "fibbonaci func", testmatrix.Values{
		"recur": &Recursive{},
		"iter":  &Iterative{},
	}),
	testmatrix.Dim("decorator", "decorator to use", testmatrix.Values{
		"none": NoDecorator,
		"memo": Memoize,
	}),
)

// TestMain tells us to run tests as defined in the matrix.
// Note: if you need more control over TestMain, you can manually call
// matrix.Init and matrix.PrintSummary instead of matrix.Run.
// See the implementation of matrix.Run for details on how to do this.
func TestMain(m *testing.M) {
	os.Exit(matrix.Run(m))
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
	decorated := s.Value("decorator").(Decorator)
	return &fixture{
		Provider: decorated(provider),
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
	return &runner{matrix.NewRunner(t)}
}
