package testmatrix

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

// supervisor supervises a set of Runners, and collates their results.
// There should be exactly one global supervisor per `go test` invocation.
type supervisor struct {
	mu         sync.Mutex
	matrixFunc MatrixFunc
	GetAddrs   func(int) []string
	fixtures   map[string]*Runner
	wg         sync.WaitGroup
}

func newSupervisor() *supervisor {
	return &supervisor{
		fixtures: map[string]*Runner{},
	}
}

// newRunner returns a new *Runner ready to run tests with all possible
// combinations of the provided Matrix. NewRunner should be called exactly once
// in each top-level TestXXX(t *testing.T) function in your package. Calling it
// more than once per top-level test may cause undefined behaviour and may
// panic.
func (s *supervisor) newRunner(t *testing.T) *Runner {
	matrix := s.matrixFunc()
	if *printInfo {
		scenarios := matrix.scenarios()
		for _, s := range scenarios {
			fmt.Printf("%s/%s\n", t.Name(), s)
		}
		t.Skip("Just printing test matrix.")
	}
	t.Helper()
	t.Parallel()
	pf := &Runner{
		t:                t,
		matrix:           matrix,
		testNames:        map[string]struct{}{},
		testNamesPassed:  map[string]struct{}{},
		testNamesSkipped: map[string]struct{}{},
		testNamesFailed:  map[string]struct{}{},
		parent:           s,
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fixtures[t.Name()] = pf
	return pf
}

// TestCount returns the number of tests that have been registered so far.
func (s *supervisor) TestCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.fixtures)
}

// PrintSummary prints a summary of tests run by top-level test and as a sum
// total. It reports tests failed, skipped, passed, and missing (when a test has
// failed to report back any status, which should not happen under normal
// circumstances.
func (s *supervisor) PrintSummary() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var total, passed, skipped, failed, missing []string
	for _, pf := range s.fixtures {
		t, p, s, f, m := pf.printSummary()
		total = append(total, t...)
		passed = append(passed, p...)
		skipped = append(skipped, s...)
		failed = append(failed, f...)
		missing = append(missing, m...)
	}

	if len(failed) != 0 {
		fmt.Printf("These tests failed:\n")
		for _, n := range failed {
			fmt.Printf("FAILED> %s\n", n)
		}
	}

	if len(missing) != 0 {
		fmt.Printf("These tests did not report status:\n")
		for _, n := range missing {
			fmt.Printf("MISSING> %s\n", n)
		}
	}

	// By default, don't print anything for 'missing' if all tests reported
	// back. In general this should happen rarely so showing it is just noise.
	var missingStr string
	if len(missing) != 0 {
		missingStr = fmt.Sprintf("%d missing ", len(missing))
	}

	summary := fmt.Sprintf("Summary: %d failed; %d skipped; %d passed; %s(total %d)",
		len(failed), len(skipped), len(passed), missingStr, len(total))
	fmt.Fprintln(os.Stdout, summary)
}
