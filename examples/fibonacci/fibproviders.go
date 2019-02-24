package fibonacci

type fibProvider interface {
	Fib(int) int
}

type (
	fibRecur struct{}
	fibIter  struct{}
	fibMemo  struct {
		results []int
	}
)

func (f *fibRecur) Fib(n int) int {
	if n < 2 {
		return n
	}
	return f.Fib(n-1) + f.Fib(n-2)
}

func (f *fibIter) Fib(n int) int {
	if n < 2 {
		return n
	}
	acc, last := 1, 1
	for i := 2; i < n; i++ {
		acc, last = acc+last, acc
	}
	return acc
}

func (f *fibMemo) Fib(n int) int {
	if len(f.results) > n {
		return f.results[n]
	}
	if n < 2 {
		f.results = append(f.results, n)
	} else {
		f.results = append(f.results, f.Fib(n-1)+f.Fib(n-2))
	}
	return f.Fib(n)
}
