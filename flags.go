package testmatrix

import (
	"flag"
)

// Flags are extra go test flags you can pass to testmatrix.
var Flags = struct {
	PrintMatrix     bool
	PrintDimensions bool
}{}

func init() {
	flag.BoolVar(&Flags.PrintDimensions, "matrix.print-dimensions", false, "list test matrix dimensions")
	flag.BoolVar(&Flags.PrintMatrix, "matrix.print-top-level", false, "list test matrix names")
}
