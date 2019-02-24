package testmatrix

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

// Runner runs tests defined in a Matrix.
type Runner struct {
	t                  *testing.T
	matrix             Matrix
	testNames          map[string]struct{}
	testNamesMu        sync.RWMutex
	testNamesPassed    map[string]struct{}
	testNamesPassedMu  sync.Mutex
	testNamesSkipped   map[string]struct{}
	testNamesSkippedMu sync.Mutex
	testNamesFailed    map[string]struct{}
	testNamesFailedMu  sync.Mutex
	parent             *Supervisor
}

func (pf *Runner) recordTestStarted(t *testing.T) {
	t.Helper()
	name := t.Name()
	pf.testNamesMu.Lock()
	defer pf.testNamesMu.Unlock()
	if _, ok := pf.testNames[name]; ok {
		t.Fatalf("duplicate test name: %q", name)
	}
	pf.testNames[name] = struct{}{}
}

// Test is a generic test.
type Test func(t *testing.T, fixture Fixture)

// FixtureFactory generates Fixtures from test and combination.
type FixtureFactory func(*testing.T, Scenario) Fixture

// Fixture is just the empty interface, it can be anything.
type Fixture interface{}

// TearableDown is a kind of Fixture that can be torn down after a test has
// finished.
type TearableDown interface {
	Teardown(*testing.T)
}

// Context as well as *testing.T. Run runs the defined test with all possible
// matrix combinations in parallel.
func (pf *Runner) teardown(t *testing.T, f Fixture) {
	if tear, ok := f.(TearableDown); ok {
		tear.Teardown(t)
	}
}

// Run is analogous to *testing.T.Run, but takes a method makeFixture that
// generates a fixture from the test and scenario, and passes that to the
// test func along with the *testing.T.
func (pf *Runner) Run(name string, makeFixture FixtureFactory, test Test) {
	for _, c := range pf.matrix.scenarios() {
		c := c
		pf.t.Run(c.String()+"/"+name, func(t *testing.T) {
			pf.recordTestStarted(t)
			defer pf.recordTestStatus(t)
			pf.parent.wg.Add(1)
			fix := makeFixture(t, c)
			defer func() {
				timeout := 10 * time.Second
				defer pf.parent.wg.Done()
				select {
				case <-time.After(10 * time.Second):
					rtLog("ERROR: Teardown took longer than %s", timeout)
				case <-func() <-chan struct{} {
					c := make(chan struct{})
					go func() {
						pf.teardown(t, fix)
						close(c)
					}()
					return c
				}():
				}
			}()
			t.Parallel()
			test(t, fix)
		})
	}
}

func (pf *Runner) recordTestStatus(t *testing.T) {
	t.Helper()
	name := t.Name()
	pf.testNamesMu.RLock()
	_, started := pf.testNames[name]
	pf.testNamesMu.RUnlock()

	statusString := "UNKNOWN"
	status := &statusString
	defer func() { rtLog("Finished running %s: %s", name, *status) }()

	if !started {
		t.Fatalf("test %q reported as finished, but not started", name)
		*status = "ERROR: Not Started"
		return
	}
	switch {
	default:
		*status = "PASSED"
		pf.testNamesPassedMu.Lock()
		pf.testNamesPassed[name] = struct{}{}
		pf.testNamesPassedMu.Unlock()
		return
	case t.Skipped():
		*status = "SKIPPED"
		pf.testNamesSkippedMu.Lock()
		pf.testNamesSkipped[name] = struct{}{}
		pf.testNamesSkippedMu.Unlock()
		return
	case t.Failed():
		*status = "FAILED"
		pf.testNamesFailedMu.Lock()
		pf.testNamesFailed[name] = struct{}{}
		pf.testNamesFailedMu.Unlock()
		return
	}
}

func testNamesSlice(m map[string]struct{}) []string {
	var s, i = make([]string, len(m)), 0
	for n := range m {
		s[i] = n
		i++
	}
	return s
}

// Quiet causes less output to be produced if set to true.
// You must not change the value of Quiet after calling Init.
var Quiet bool

func rtLog(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func (pf *Runner) printSummary() (total, passed, skipped, failed, missing []string) {
	t := pf.t
	t.Helper()
	total = testNamesSlice(pf.testNames)
	passed = testNamesSlice(pf.testNamesPassed)
	skipped = testNamesSlice(pf.testNamesSkipped)
	failed = testNamesSlice(pf.testNamesFailed)

	if !Quiet {
		summary := fmt.Sprintf("%s summary: %d failed; %d skipped; %d passed (total %d)",
			t.Name(), len(failed), len(skipped), len(passed), len(total))
		t.Log(summary)
		fmt.Fprintln(os.Stdout, summary)
	}

	missingCount := len(total) - (len(passed) + len(failed) + len(skipped))
	if missingCount != 0 {
		for t := range pf.testNamesPassed {
			delete(pf.testNames, t)
		}
		for t := range pf.testNamesSkipped {
			delete(pf.testNames, t)
		}
		for t := range pf.testNamesFailed {
			delete(pf.testNames, t)
		}
		for t := range pf.testNames {
			missing = append(missing, t)
		}
	}
	return total, passed, skipped, failed, missing
}
