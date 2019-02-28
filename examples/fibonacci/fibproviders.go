// Package fibonacci is a contrived example to show
// github.com/samsalisbury/testmatrix in action.
//
// You are recommended not to use this for anything other than demonstration
// purposes.
package fibonacci

// Provider is a provider of a Fibonacci function.
type Provider interface {
	Fib(int) int
}

// Decorator decorates a Provider.
type Decorator func(Provider) Provider

type (
	// Recursive gets Fibonacci sequences the hard way.
	Recursive struct{}

	// Iterative gets Fibonacci sequences without burning the stack.
	Iterative struct{}

	// Memoized remembers its results.
	Memoized struct {
		results  []int
		provider Provider
	}
)

// Fib returns the nth Fibonacci number the hard way.
func (f *Recursive) Fib(n int) int {
	if n < 2 {
		return n
	}
	return f.Fib(n-1) + f.Fib(n-2)
}

// Fib returns the nth Fibonacci number the iterative way.
func (f *Iterative) Fib(n int) int {
	if n < 2 {
		return n
	}
	acc, last := 1, 1
	for i := 2; i < n; i++ {
		acc, last = acc+last, acc
	}
	return acc
}

// NewMemoized returns a Memoized version of p.
func NewMemoized(p Provider) *Memoized {
	return &Memoized{provider: p}
}

// Memoize decorates p with memoization.
func Memoize(p Provider) Provider {
	return NewMemoized(p)
}

// NoDecorator is a noop Decorator.
func NoDecorator(p Provider) Provider {
	return p
}

// Fib returns the nth Fibonacci number, sometimes from memory.
func (f *Memoized) Fib(n int) int {
	if len(f.results) > n {
		return f.results[n]
	}
	if n < 2 {
		f.results = append(f.results, n)
	} else {
		f.results = append(f.results, f.provider.Fib(n-1)+f.provider.Fib(n-2))
	}
	return f.Fib(n)
}
