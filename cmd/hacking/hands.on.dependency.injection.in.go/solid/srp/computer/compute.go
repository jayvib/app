package computer

import (
	"fmt"
	"io"
	"os"
)

type ProcessorFunc func(d interface{}) interface{}

func Compute(data []interface{}, fn ProcessorFunc) (result []interface{}) {
	for _, n := range data {
		result = append(result, fn(n))
	}
	return
}

var Separator = "\n"

func Fcompute(w io.Writer, d []interface{}, processor ProcessorFunc) {
	res := Compute(d, processor)
	for _, r := range res {
		fmt.Fprintf(w, "%v%s", r, Separator)
	}
}

func PrintComputation(d []interface{}, processorFunc ProcessorFunc) {
	Fcompute(os.Stdout, d, processorFunc)
}