package fibonacci

import (
	"strconv"
	"testing"
)

func TestFib(t *testing.T) {
	r := newRunner(t)
	r.Run("test one", makeFixture, func(t *testing.T, f *fixture) {
		testCases := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
		for in, want := range testCases {
			t.Run(strconv.Itoa(in), func(t *testing.T) {
				got := f.Fib(in)
				if got != want {
					t.Errorf("Fib(%d) == %d; want %d", in, got, want)
				}
			})
		}
	})
}
