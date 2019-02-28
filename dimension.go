package testmatrix

// Dimension represents a dimension in the test matrix.
type Dimension struct {
	// name is the name of this dimension, used to look up values in calculated
	// Scenarios.
	name string
	// desc is used in help text when using the -matrix flag.
	desc string
	// values is a map of named possible values for this Dimension.
	// The name used here forms part of the sub-test path.
	values Values
}

// Dim returns a new Dimension.
func Dim(name, desc string, values Values) Dimension {
	return Dimension{
		name:   name,
		desc:   desc,
		values: values,
	}
}
