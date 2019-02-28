package testmatrix

import (
	"flag"
	"testing"
)

// sup is global state, and keeps track of tests started and finished, allowing
// us to print a summary at the end. The rationale for using global state here
// is that we want summary information for a single invocation of `go test`, on
// a single package, which implies that this state should be global to that
// invocation.
var sup = newSupervisor()

var opts = DefaultOpts()

// Opts are global options.
type Opts struct {
	BeforeAll func()
}

// DefaultOpts returns the default opts.
func DefaultOpts() Opts {
	return Opts{}
}

// Run wraps all initialisation logic, runs the tests, and returns the
// appropriate exit code. This should only be called once, in TestMain.
//
// Your test package should declare a global *Supervisor and pass a pointer to
// that here, it will be configured an populated ready to use in creating tests.
func Run(m *testing.M, matrixFunc func() Matrix, config ...func(*Opts)) (exitCode int) {
	Init(matrixFunc, config...)
	defer sup.PrintSummary()
	return m.Run()
}

// Init ensures flags are parsed, and makes decision on whether to actually run
// tests, or just print summaries etc. It invokes the BeforeAll hook in the case
// that we are actually intending to run tests.
func Init(matrixFunc func() Matrix, config ...func(*Opts)) {
	if !flag.Parsed() {
		flag.Parse()
	}
	for _, c := range config {
		c(&opts)
	}
	sup.matrixFunc = matrixFunc
	if *printInfo {
		matrixFunc().PrintDimensions()
		return
	}
	if opts.BeforeAll != nil {
		opts.BeforeAll()
	}
}

// PrintSummary prints the summary of all tests run/passed/failed etc.
// It must be called after all tests have run to completion.
//
// If using the Run func, you don't need to additionally call this.
func PrintSummary() {
	sup.PrintSummary()
}
