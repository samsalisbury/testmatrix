package testmatrix

import "testing"

// T implements the *testing.T interface.
// We depend on this interface rather than *testing.T
// as it makes this package easier to test.
//
// You should always pass a real *testing.T into your own tests.
type T interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Parallel()
	Run(name string, f func(*testing.T)) bool
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
}

// M implements the *testing.M interface.
// We depend on this interface rather than *testing.M
// as it makes this package easier to test.
type M interface {
	Run() int
}
