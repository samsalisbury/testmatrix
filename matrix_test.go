package testmatrix

import (
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		matrixFunc func() Matrix
		wantCount  int
	}{
		{"empy", func() Matrix {
			return New()
		}, 0},
		{"one", func() Matrix {
			return New(makeTestDim(1, 1))
		}, 1},
		{"two", func() Matrix {
			return New(makeTestDims(2, alwaysOneValue)...)
		}, 2},
		{"5", func() Matrix {
			return New(makeTestDims(5, alwaysOneValue)...)
		}, 5},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotMatrix, want := tc.matrixFunc(), tc.wantCount
			got := len(gotMatrix.dimensions)
			if want != got {
				t.Errorf("got %d dimensions; want %d", got, want)
			}
		})
	}
}

func TestNew_error(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		matrixFunc func() Matrix
		wantPanic  string
	}{
		{"dupe", func() Matrix {
			return New(
				Dim("dim1", "", Values{"a": struct{}{}}),
				Dim("dim1", "", Values{"b": struct{}{}}),
			)
		}, `duplicate dimension name "dim1"`},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			want := tc.wantPanic
			func() {
				defer func() {
					gotPanic := recover()
					if gotPanic == nil {
						t.Fatalf("did not panic; want panic with %q", want)
					}
					got := fmt.Sprint(gotPanic)
					if !strings.Contains(got, want) {
						t.Errorf("got panic %q; want it to contain %q", got, want)
					}
				}()
				tc.matrixFunc()
			}()
		})
	}
}

func TestMatrix_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   Matrix
		want string
	}{
		{
			"empty",
			New(),
			"",
		},
		{
			"onedim-noval",
			New(Dim("dimname", "yolo", Values{})),
			"dimname\t\n-\t\n",
		},
		{
			"onedim-oneval",
			New(Dim("dim1", "yolo", Values{"dim1val1": 1})),
			"dim1\t\n-\t\ndim1val1\t\n",
		},
		{
			"twodim-twoval",
			New(
				Dim("dim1", "yolo", Values{"dim1val1": 1, "dim1val2": 2}),
				Dim("dim2", "yolo", Values{"dim2val1": 1, "dim2val2": 2}),
			),
			"dim1\tdim2\t\n-\t-\t\ndim1val1\tdim2val1\t\ndim1val2\tdim2val2\t\n",
		},
		{
			"twodim-unevenval1",
			New(
				Dim("dim1", "yolo", Values{"dim1val1": 1}),
				Dim("dim2", "yolo", Values{"dim2val1": 1, "dim2val2": 2}),
			),
			"dim1\tdim2\t\n-\t-\t\ndim1val1\tdim2val1\t\n\tdim2val2\t\n",
		},
		{
			"twodim-unevanval2",
			New(
				Dim("dim1", "yolo", Values{"dim1val1": 1, "dim1val2": 2}),
				Dim("dim2", "yolo", Values{"dim2val2": 2}),
			),
			"dim1\tdim2\t\n-\t-\t\ndim1val1\tdim2val2\t\ndim1val2\t\t\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, want := tc.in.String(), tc.want
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		})
	}

}
