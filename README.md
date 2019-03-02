# testmatrix [![CircleCI](https://circleci.com/gh/samsalisbury/testmatrix/tree/master.svg?style=svg)](https://circleci.com/gh/samsalisbury/testmatrix/tree/master) [![GoDoc](https://godoc.org/github.com/samsalisbury/testmatrix?status.svg)](https://godoc.org/github.com/samsalisbury/testmatrix)

testmatrix is a fully `go test` compatible library
for expanding your test coverage
with minimal overhead.

## What?

1. Define an N-dimensional matrix of factors that might affect your system.
2. Have your tests run against every possible combination of those factors.
3. Discover hard to find bugs before your end users.

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
go test . -tm.info # Print matrix info without running tests.
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
		testmatrix.Dim("git", "version of git", testmatrix.Values{
			"2.19.0": struct{}{},
			"1.0.0": struct{}{},
		}),
		testmatrix.Dim("docker", "version of docker", testmatrix.Values{
			"1.0.0": "https://download.docker.com/v1.0.0",			
			"2.0.0": "https://download.docker.com/v2.0.0",			
		}),
		...
	)
}
```

#### Define your fixture

testmatrix insists you pass a fixture to each test. The fixture can be anything
you want, but is typically a struct containing information about a test environment
you have spun up for this test in particular. It may have helper methods attached etc.

The fixture may implement a Teardown method, see below for details.

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

#### Do some type wrapping

You can make your tests somewhat easier to write
by defining your own `testFunc` and `fixtureFunc` types.
(Rather than relying on the weakly typed
`testmatrix.Test` and `testmatrix.FixtureFactory` respectively,
in your own tests.)
You can then add a strongly-typed wrapper around testmatrix.Runner,
and redefine the `Run` method to map between these weakly and strongly typed
tests and fixtures.

At the moment, this is a little ugly, but typically looks like this:

```go

// runner wraps the *testmatrix.Runner so we can add our own Run method.
type runner struct{ *testmatrix.Runner }

// newRunner returns a newly configured runner.
// You need one of these for each top-level test.
func newRunner(t *testing.T) *runner {
	return &runner{testmatrix.NewRunner(t)}
}

// testFunc is the strongly typed test function signature you will use to write your test.
type testFunc func(*testing.T, *fixture)

// fixtureFunc is a strongly typed fixture generation function.
type fixtureFunc func(*testing.T, testmatrix.Scenario) (*fixture)

// Run is analogous to *testing.T.Run in that it creates a subtest.
// Run accepts your strongly typed fixtureFunc and testFunc, wraps them up and
// passes them through to the generic testmatrix.Runner.Run for execution.
func (r *runner) Run(name string, makeFixture fixtureFunc, test testFunc) {
	r.Runner.Run(name,
		// Return testmatrix.Fixture which is really a *fixture.
		func(t *testing.T, s testmatrix.Scenario) testmatrix.Fixture {
			return makeFixture(t, s)
		},
		// Cast that testmatrix.Fixture back to the strongly typed *fixture
		// we know it really to be...
		func(t *testing.T, f testmatrix.Fixture) {
			test(t, f.(*fixture))
		},
	)
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

### Fixture Teardown

TODO: Document this.
