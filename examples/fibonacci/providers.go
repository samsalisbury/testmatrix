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

// Recursive gets Fibonacci sequences the hard way.
type Recursive struct{}

// Iterative gets Fibonacci sequences without burning the stack.
type Iterative struct{}

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
