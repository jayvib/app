package calculator

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

var DefaultStdOutPrinter = &stdOutPrinter{w: os.Stdout}

type Calculator interface {
	Calculate(path string) error
}

type Printer interface {
	Output(data map[string]float64)
}

// NewCalculator accepts printer and return a calculator.
// When p is nil it will use the default stdout printer.
func NewCalculator(p Printer) Calculator {
	if p == nil {
		p = DefaultStdOutPrinter
	}

	return &calculatorImpl{
		data: make(map[string]float64),
		out: p,
	}
}

type calculatorImpl struct {
	data map[string]float64
	out Printer
}

func (c *calculatorImpl) Calculate(path string) error {
	// Do some calculation
	c.out.Output(c.data)
	return nil
}

type stdOutPrinter struct {
	w io.Writer
}

func (s *stdOutPrinter) Output(d map[string]float64) {
	for k, v := range d {
		fmt.Fprintf(s.w, "%s: %g\n", k, v)
	}
}

type csvPrinter struct {
	w io.Writer
}

func (s *csvPrinter) Output(d map[string]float64) {
	csvWriter := csv.NewWriter(s.w)
	defer csvWriter.Flush()
	for k, v := range d {
		line := []string{
			k,
			fmt.Sprintf("%g", v),
		}
		csvWriter.Write(line)
	}


}
