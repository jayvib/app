package calculator

import (
	"fmt"
	"io"
)

// Calculator calculates the test coverage for a directory
// and it's sub-directories
type Calculator struct {
	// coverage data populated by `Calculate()` method
	data map[string]float64
}

// Calculate will calculate the coverage
func (c *Calculator) Calculate(path string) error {
	// run `go test -cover ./[path]/...` and store the results
	return nil
}

// Output will print the coverage data to the supplied writer
func (c *Calculator) Output(writer io.Writer) {
	for path, result := range c.data {
		fmt.Fprintf(writer, "%s -> %.1f\n", path, result)
	}
}

// OutputCSV will print the coverage data to the supplied writer
func (c Calculator) OutputCSV(writer io.Writer) {  // Decided that we also needed to ouput the results to CSV.
	for path, result := range c.data { // It added more responsibility to the struct and, in doing so, we have added complexity.
		fmt.Fprintf(writer, "%s,%.1f\n", path, result)
	}
}
