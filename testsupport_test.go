package testmatrix

import (
	"fmt"
	"testing"
)

func makeTestDims(count int, valueCountFunc func(index int) (valueCount int)) []Dimension {
	ds := make([]Dimension, count)
	for i := 0; i < count; i++ {
		ds[i] = makeTestDim(i, valueCountFunc(i))
	}
	return ds
}

func alwaysOneValue(int) int {
	return 1
}

func alwaysNValues(n int) func(int) int {
	return func(int) int { return n }
}

func TestMakeTestDims(t *testing.T) {
	t.Parallel()
	cases := []struct {
		namePrefix             string
		dimCount, wantValCount int
		valCountFunc           func(int) int
	}{
		{"alwaysOneValue", 1, 1, alwaysOneValue},
		{"alwaysTwoValues", 1, 2, alwaysNValues(2)},
		{"alwaysTenValues", 1, 10, alwaysNValues(10)},
		{"alwaysTenValues", 10, 10, alwaysNValues(10)},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%s/%d/%d", tc.namePrefix, tc.dimCount, tc.wantValCount), func(t *testing.T) {
			t.Parallel()
			got := makeTestDims(tc.dimCount, tc.valCountFunc)
			if len(got) != tc.dimCount {
				t.Fatalf("got %d Dimensions; want %d", len(got), tc.dimCount)
			}
			for _, d := range got {
				if len(d.values) != tc.wantValCount {
					t.Errorf("got %d values; want %d", len(d.values), tc.wantValCount)
				}
			}
		})

	}
}

// makeTestDim returns a Dimension named "dim<index>" with valueCount values
// for testing purposes.
func makeTestDim(index, valueCount int) Dimension {
	return Dim(fmt.Sprintf("dim%d", index), "", makeTestValues(index, valueCount))
}

func TestMakeTestDim(t *testing.T) {
	t.Parallel()
	cases := []struct {
		index, valCount int
		wantName        string
	}{
		{1, 1, "dim1"},
		{2, 1, "dim2"},
		{2, 2, "dim2"},
		{12, 12, "dim12"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%s/%d", tc.wantName, tc.valCount), func(t *testing.T) {
			t.Parallel()
			got := makeTestDim(tc.index, tc.valCount)
			if got.name != tc.wantName {
				t.Errorf("got name %q; want %q", got.name, tc.wantName)
			}
			if len(got.values) != tc.valCount {
				t.Errorf("got %d values; want %d", len(got.values), tc.valCount)
			}
		})
	}
}

// makeTestValues returns Values where each value equals its key
// for testing purposes.
// Each value name is in the format dim<dimIndex>val<valIndex>
// where valIndex is in the range 1..count inclusive.
func makeTestValues(dimIndex, count int) Values {
	v := Values{}
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("dim%dval%d", dimIndex, i+1)
		v[name] = name
	}
	return v
}

// TestMakeTestValues exists so we don't doubt our test data.
func TestMakeTestValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		dimIndex, count int
		wantNames       []string
	}{
		{1, 1, []string{"dim1val1"}},
		{2, 1, []string{"dim2val1"}},
		{1, 2, []string{"dim1val1", "dim1val2"}},
		{3, 3, []string{"dim3val1", "dim3val2", "dim3val3"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(fmt.Sprintf("%d/%d", tc.dimIndex, tc.count), func(t *testing.T) {
			t.Parallel()
			got := makeTestValues(tc.dimIndex, tc.count)
			if len(got) != len(tc.wantNames) {
				t.Fatalf("got %d values; want %d", len(got), len(tc.wantNames))
			}
			for _, wantName := range tc.wantNames {
				got, ok := got[wantName]
				if !ok {
					t.Fatalf("missing value %q", wantName)
				}
				if got != wantName {
					t.Errorf("got %q -> %q; want %q", wantName, got, wantName)
				}
			}
		})
	}
}
