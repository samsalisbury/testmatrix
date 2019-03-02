package testmatrix

import (
	"fmt"
	"sort"
	"strings"
)

// Dimensions is a complete Matrix description, with all dimension's names and
// possible values.
type Dimensions map[string]Values

// Values is mapping of all the possible values' names (for a single dimension)
// to their values for use in tests.
type Values map[string]interface{}

// Matrix is a compiled set of named dimensions and possible values.
// Every combination of a single value from each dimension forms a Scenario.
// Each test you define will be run once for each possible Scenario.
type Matrix struct {
	sup                   *supervisor
	orderedDimensionNames []string
	orderedDimensionDescs []string
	dimensions            Dimensions
}

// Scenario is a single combination of values from a Matrix.
type Scenario []Binding

// Binding is a named value belonging to a particular dimension.
type Binding struct {
	Dimension, Name string
	Value           interface{}
}

// New returns a new Matrix.
// Note that the order of Dimensions is significant in determining the name of
// each sub-test run. For example if you have 3 dimensions,
// and a sub-test is being run with values
// "a" for the first, "b" for the second, and "c" for the third dimension,
// the test name will be "<root>/a/b/c/<subtest>".
func New(dimensions ...Dimension) Matrix {
	m := Matrix{
		sup:        newSupervisor(),
		dimensions: Dimensions{},
	}
	for _, d := range dimensions {
		m.addDimension(d.name, d.desc, d.values)
	}
	return m
}

// FixedDimension returns a new Matrix based on m with one of its dimensions
// fixed to a particular value. This can be used when writing tests where
// only a single value for one particular dimension is appropriate.
func (m Matrix) FixedDimension(dimensionName, valueName string) Matrix {
	return m.clone(func(dimension, value string) bool {
		return dimension != dimensionName || value == valueName
	})
}

// String returns a description of this matrix.
func (m Matrix) String() string {
	cols := make([][]string, len(m.orderedDimensionNames))
	var maxRows int
	for i, name := range m.orderedDimensionNames {
		cols[i] = append(cols[i], name, "-")
		rowCount := 2 // 2 for the column header and divider
		for valueName := range m.dimensions[name] {
			rowCount++
			cols[i] = append(cols[i], valueName)
		}
		sort.Strings(cols[i][1:])
		if rowCount > maxRows {
			maxRows = rowCount
		}
	}
	var out string
	for i := 0; i < maxRows; i++ {
		for _, c := range cols {
			if len(c) > i {
				out += c[i]
			}
			out += "\t"
		}
		out += "\n"
	}
	return out
}

// PrintDimensions writes the dimensions and allowed values
// as a table to stdout.
func (m Matrix) PrintDimensions() {
	fmt.Println(m.String())
}

// addDimension adds a new dimension to this matrix with the provided name
// and desc which is used in help text when using -matrix flag on 'go test'.
// The values are a map of short value names to concrete representations, which
// are passed to tests. The names of values map to parts of the sub-test path
// for 'go test -run' flag.
func (m *Matrix) addDimension(name, desc string, values Values) {
	if _, ok := m.dimensions[name]; ok {
		panic(fmt.Sprintf("duplicate dimension name %q", name))
	}
	if len(values) == 0 {
		panic(fmt.Sprintf("no values for dimension %q", name))
	}
	m.dimensions[name] = values
	m.orderedDimensionNames = append(m.orderedDimensionNames, name)
	m.orderedDimensionDescs = append(m.orderedDimensionDescs, desc)
}

func (m Matrix) clone(include func(dimension, value string) bool) Matrix {
	n := m
	n.dimensions = Dimensions{}
	for name, values := range m.dimensions {
		nv := map[string]interface{}{}
		for vn, v := range values {
			if include(name, vn) {
				nv[vn] = v
			}
		}
		n.dimensions[name] = nv
	}
	return n
}

func (m *Matrix) scenarios() []Scenario {
	combos := [][]Scenario{}
	for _, d := range m.orderedDimensionNames {
		c := []Scenario{}
		dim := m.dimensions[d]
		valNames := []string{}
		for name := range dim {
			valNames = append(valNames, name)
		}
		sort.Strings(valNames)
		for _, name := range valNames {
			c = append(c, Scenario{
				Binding{
					Dimension: d,
					Name:      name,
					Value:     dim[name],
				},
			})
		}
		combos = append(combos, c)
	}
	return product(combos...)
}

func product(slices ...[]Scenario) []Scenario {
	res := slices[0]
	for _, s := range slices[1:] {
		res = mult(res, s)
	}
	return res
}

func mult(a, b []Scenario) []Scenario {
	res := make([][]Scenario, len(a)*len(b))
	n := 0
	for _, aa := range a {
		for _, bb := range b {
			res[n] = []Scenario{aa, bb}
			n++
		}
	}
	slice := make([]Scenario, len(res))
	for i, r := range res {
		slice[i] = concat(r)
	}
	return slice
}

func concat(scenarios []Scenario) Scenario {
	res := scenarios[0]
	for _, c := range scenarios[1:] {
		res = append(res, c...)
	}
	return res
}

// String returns the sub-test path of this Scenario. E.g.
// dim1ValueName/dim2ValueName[/...].
func (c Scenario) String() string {
	var names []string
	for _, p := range c {
		names = append(names, p.Name)
	}
	return strings.Join(names, "/")
}

// Map returns a map of dimension name to specific value for this Scenario.
// This is useful where we want to look up a value by dimension name.
func (c Scenario) Map() map[string]interface{} {
	res := make(map[string]interface{}, len(c))
	for _, p := range c {
		res[p.Dimension] = p.Value
	}
	return res
}

// Value returns the value for the named dimension in this Scenario.
func (c Scenario) Value(dimension string) interface{} {
	v, ok := c.Map()[dimension]
	if !ok {
		panic(fmt.Sprintf("scenario contains no value for dimension %q", dimension))
	}
	return v
}
