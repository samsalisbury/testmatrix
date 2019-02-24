package testmatrix_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/samsalisbury/testmatrix"
)

// sup is the global matrix supervisor, used to collate all test results.
var sup *testmatrix.Supervisor

// matrix() returns the full test matrix. Note if certain matrix combinations
// should not be tested, it is possible to make such derivations using the
// Matrix.FixedDimension method, and more like it.
//
// The matrix func should always return a fresh matrix, so that references
// are not shared between tests.
func matrix() testmatrix.Matrix {
	return testmatrix.New(
		testmatrix.Dimension{
			Name: "fib",
			Desc: "fibbonaci func",
			Values: testmatrix.Values{
				"recur": &fibRecur{},
				"iter":  &fibIter{},
				"memo":  &fibMemo{},
			},
		},
	)
}

func TestMain(m *testing.M) {
	os.Exit(testmatrix.Run(m, &sup, matrix))
}

func TestExample(t *testing.T) {
	m := matrix()
	r := runner{sup.NewRunner(t, m)}
	r.Run("test one", defaultFixture, func(t *testing.T, f *fixture) {
		testCases := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
		for in, want := range testCases {
			t.Run(strconv.Itoa(in), func(t *testing.T) {
				got := f.Fib(in)
				if got != want {
					t.Errorf("Fib(%d) == %d; want %d", in, got, want)
				}
			})
		}
	})
}

type fibProvider interface {
	Fib(int) int
}

type fibRecur struct{}

func (f *fibRecur) Fib(n int) int {
	if n < 2 {
		return n
	}
	return f.Fib(n-1) + f.Fib(n-2)
}

type fibIter struct{}

func (f *fibIter) Fib(n int) int {
	if n < 2 {
		return n
	}
	acc, last := 1, 1
	for i := 2; i < n; i++ {
		acc, last = acc+last, acc
	}
	return acc
}

type fibMemo struct {
	results []int
}

func (f *fibMemo) Fib(n int) int {
	if len(f.results) > n {
		return f.results[n]
	}
	if n < 2 {
		f.results = append(f.results, n)
	} else {
		f.results = append(f.results, f.Fib(n-1)+f.Fib(n-2))
	}
	return f.Fib(n)
}

// scenario holds one value for each matrix dimension as described above.
type scenario struct {
	fibProvider
}

func unwrapScenario(s testmatrix.Scenario) scenario {
	return scenario{
		fibProvider: s.Value("fib").(fibProvider),
	}
}

type runner struct{ *testmatrix.Runner }

// fixtureConfig defines our desired initial conditions.
type fixtureConfig struct {
	PrewarmValues []int
}

func (fc *fixtureConfig) Init(t *testing.T, scenario testmatrix.Scenario) *fixture {
	s := unwrapScenario(scenario)
	for _, n := range fc.PrewarmValues {
		s.Fib(n)
	}
	return &fixture{
		s.fibProvider,
	}
}

// fixture represents a fully realised fixture, generated from the fixtureConfig
// and injected scenario.
type fixture struct {
	fibProvider
}

// defaultFixture returns the default fixture config.
func defaultFixture() *fixtureConfig {
	return &fixtureConfig{}
}

func prewarmed(n ...int) *fixtureConfig {
	return &fixtureConfig{PrewarmValues: n}
}

type test func(*testing.T, *fixture)

func (r *runner) Run(name string, config func() *fixtureConfig, test test) {
	r.Runner.RunScenario(name, func(t *testing.T, s testmatrix.Scenario, lf *testmatrix.LateFixture) {
		fix := config().Init(t, s)
		lf.Set(fix)
		test(t, fix)
	})
}
