# testmatrix

## What?

Write a set of tests, have each of them run in a multitude of different ways.

## Why?

Sometimes you need to write software that behaves the same way no matter what underlying
libraries/external programs/services it is relying on. Writing a separate test for
each combination of these is impractical, so we don't do it.

Sometimes you may want to compare competing implementations of a type, not just for
correctness but for speed. That's why this lib also supports benchmarking.

This library aims to make it easier to test your code against a matrix of underlying
conditions, and provides some helpful features for dealing with an explosion in test
results.

## How?

Write tests using plan old Go. Use standard `go test` command to run tests.
Use the `go test -run` flag to run specific tests (with specific matrix combinations).
Additional flags help you navigate the matrix:

```sh
go test -matrix.expand     # List all matrix expansions.
go test -matrix.dimensions # Describe all dimensions of the matrix.
```

### Writing Tests

There is a little boilerplate in setting this up nicely. The steps are:

0. Top-level boilerplate.
1. Define your matrix
2. Define your fixture
3. Write a Run wrapper (optional but worthwhile)
4. Write your tests

#### Top-level boilerplate

We need to tell the Go test runner to run tests defined in the test matrix.
This lib also keeps track of tests run to provide summary information, and optional
hooks for various stages of the test lifecycle.

```go
func TestMain(m *testing.M) {
	os.Exit(testmatrix.Run(m, makeMatrix))
}
```


#### Define your matrix

Each value for each dimension in the matrix will be multiplied by all values for all
other dimensions.
This means adding a dimension can add a lot of computational overhead, but also
means you end up with broader coverage.

```go
func makeMatrix() *testmatrix.Matrix {
	return testmatrix.New(
		testmatrix.NewDimension("git", "version of git", "1", "2", "3"),
		testmatrix.NewDimension("docker", "version of docker", "1", "2", "3"),
		testmatrix.NewDimension("kubectl", "version kubectl", "1", "2", "3"),
		testmatrix.NewDimension("kubeapi", "version kubeapi", "1", "2", "3"),
	)
}
```

#### Define your fixture

testmatrix insists you pass a fixture to each test. The fixture can be anything
you want, but is typically a struct containing information about a test environment
you have spun up for this test in particular. It may have helper methods attached etc.

You must define a func that takes a `*testing.T` and a `testmatrix.Scenario` to create
your fixture. This is invoked automatically by testmatrix just before each test is run,
and the resulting fixture is passed to that test.

```go
type fixture struct {
	// Anything you want to pass to the tests.
}

// makeFixture should create an isolated fixture for each test.
// It can for example use the test name to help with that isolation.
// The scenario passed in will be a single combination from the matrix,
// and should be used to set things up appropriately (e.g. launch and
// configure docker containers, acquire the right version of binaries
// as specified in the scenario etc.)
func makeFixture(t *testing.T, scenario *testmatrix.Scenario) *fixture {
	return &fixture{
		// Whatevs.
	}
}
```

#### Write a Run wrapper

You can make your tests somewhat easier to write, and less ugly, by
adding a Run wrapper. This wraps the call to `testmatrix.Matrix.RunScenario` and
performs any necessary unwrapping of types, so your tests can rely on strongly-typed
fixtures.

At the moment, this is pretty ugly, but typically looks like this:

```go
// Define your test type (analogous to testmatrix.Test but takes a *fixture
// instead of interface{}).
type test func(*testing.T, *fixture)

// runner wraps the *testmatrix.Runner so you get all of its methods by default.
type runner struct{ *testmatrix.Runner }

type fixtureFunc func(*testing.T, testmatrix.Scenario) (*fixture)

// Run is analogous to *testing.T.Run in that it creates a subtest.
func (r *runner) Run(name string, makeFixture fixtureFunc, test test) {
	r.Runner.RunScenario(name, func(t *testing.T, s testmatrix.Scenario, lf *testmatrix.LateFixture) {
		fix := makeFixture(t, s)
		lf.Set(fix)
		test(t, fix)
	})
}

// NewRunner returns a newly configured runner. You need one of these for each top-level
// test.
func newRunner() *runner {
	return &runner{sup.NewRunner()}
}
```

#### Write your tests

Now you have the boilerplate set up, you're ready to write some tests.

```go
func TestBlahBlah(t *testing.T) {
	r := newRunner()
	r.Run("test one", makeFixture, func(t *testing.T, f *fixture) {
		// Write a standard go test, using info from your fixture.
		f.Fatalf("this test blew up!")
	})
}
```
