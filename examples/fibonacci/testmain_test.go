package fibonacci

import (
	"os"
	"testing"

	"github.com/samsalisbury/testmatrix"
)

// sup is the global matrix supervisor, used to collate all test results.
var sup *testmatrix.Supervisor

func TestMain(m *testing.M) {
	os.Exit(testmatrix.Run(m, &sup, matrix))
}

// matrix() returns the full test matrix. Note if certain matrix combinations
// should not be tested, it is possible to make such derivations using the
// Matrix.FixedDimension method, and more like it. Do this in top-level tests.
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

// scenario holds one value for each matrix dimension as described above.
type scenario struct {
	fibProvider
}

func unwrapScenario(s testmatrix.Scenario) scenario {
	return scenario{
		fibProvider: s.Value("fib").(fibProvider),
	}
}
