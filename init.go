package testmatrix

import (
	"flag"
	"testing"
)

// Opts are global options for Supervisor.
type Opts struct {
	BeforeAll func() error
}

// DefaultOpts returns the default opts.
func DefaultOpts() Opts {
	return Opts{}
}

// Run wraps all initialisation logic, runs the tests, and returns the
// appropriate exit code.
//
// Your test package should declare a global *Supervisor and pass a pointer to
// that here, it will be configured an populated ready to use in creating tests.
func Run(m *testing.M, s **Supervisor, matrixFunc func() Matrix, config ...func(*Opts)) (exitCode int) {
	if !flag.Parsed() {
		flag.Parse()
	}
	opts := DefaultOpts()
	for _, c := range config {
		c(&opts)
	}
	*s = Init(matrixFunc, opts)
	defer (*s).PrintSummary()
	return m.Run()
}

// Init must be called from TestMain after flag.Parse, to initialise a new
// Supervisor. If Init returns nil, then tests will not be run this time (e.g.
// because we are just listing tests or printing the matrix def etc.)
func Init(defaultMatrix func() Matrix, opts Opts) *Supervisor {
	if !flag.Parsed() {
		flag.Parse()
	}
	runRealTests := !(Flags.PrintMatrix || Flags.PrintDimensions)
	if Flags.PrintDimensions {
		defaultMatrix().PrintDimensions()
	}
	if !runRealTests {
		return nil
	}
	if opts.BeforeAll != nil {
		if err := opts.BeforeAll(); err != nil {
			panic(err)
		}
	}
	return NewSupervisor()
}
