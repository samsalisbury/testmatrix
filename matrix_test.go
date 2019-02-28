package testmatrix

import (
	"testing"
)

func TestMatrix_String(t *testing.T) {

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

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, want := c.in.String(), c.want
			if got != want {
				t.Errorf("got %q; want %q", got, want)
			}
		})
	}

}
