package fibonacci

// Memoized remembers its results.
type Memoized struct {
	results  []int
	provider Provider
}

// NewMemoized returns a Memoized version of p.
func NewMemoized(p Provider) *Memoized {
	return &Memoized{provider: p}
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
