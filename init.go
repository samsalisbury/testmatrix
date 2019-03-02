// Package testmatrix enables you to define an N-dimensional matrix
// of named variables, and have each test you write be run once for
// each possible combination of those variables.
//
// This was largely written to support testing of github.com/opentable/sous
// but should be applicable to any situation where you can write tests that
// expect to see the same behaviour externally for differing implementations
// of underlying libraries/types/functions/external binaries etc.
//
// See the examples directory for some contrived examples to get you started.
package testmatrix

import (
	"flag"
)

var opts = DefaultOpts()

// Opts are global options.
type Opts struct {
	BeforeAll     func()
	PrintInfoOnly bool
}

// ShouldRunTests returns true if we want to actually run tests, not just print
// info about the matrix.
func (o Opts) ShouldRunTests() bool {
	return !o.PrintInfoOnly
}

// DefaultOpts returns the default opts.
func DefaultOpts() Opts {
	return Opts{}
}

// Run wraps all initialisation logic, runs the tests, and returns the
// appropriate exit code. This should only be called once, in TestMain.
//
// You should pass the *testing.M from TestMain as the first parameter.
// We depend in the interface M for testing purposes.
func (m *Matrix) Run(testingM M, config ...func(*Opts)) (exitCode int) {
	if !m.Init(config...).ShouldRunTests() {
		return 0
	}
	defer m.sup.PrintSummary()
	return testingM.Run()
}

// Init ensures flags are parsed, and makes decision on whether to actually run
// tests, or just print summaries etc. It invokes the BeforeAll hook in the case
// that we are actually intending to run tests.
// If Init returns false, then the user does not intend to run tests, only to
// print diagnostic information.
func (m *Matrix) Init(config ...func(*Opts)) Opts {
	if !flag.Parsed() {
		flag.Parse()
	}
	for _, c := range config {
		c(&opts)
	}
	if *printInfo {
		m.PrintDimensions()
		opts.PrintInfoOnly = true
		return opts
	}
	if opts.BeforeAll != nil {
		opts.BeforeAll()
	}
	return opts
}

// PrintSummary prints the summary of all tests run/passed/failed.
// It must be called after all tests have run to completion.
//
// If using the Run func, you don't need to additionally call this.
func (m *Matrix) PrintSummary() {
	m.sup.PrintSummary()
}
