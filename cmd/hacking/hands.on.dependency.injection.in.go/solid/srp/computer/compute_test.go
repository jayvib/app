package computer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCompute(t *testing.T) {

	cases := []struct {
		name      string
		input     []interface{}
		want      []interface{}
		processor ProcessorFunc
	}{
		{
			name: "It accepts int to be calculate",
			input: []interface{}{1, 2, 3},
			want: []interface{}{2, 3, 4},
			processor: func(d interface{}) interface{} {
				n := d.(int)
				return n + 1
			},
		},
		{
			name: "It accepts string to be calculate",
			input: []interface{}{"luffy", "sanji", "zoro"},
			want: []interface{}{"LUFFY", "SANJI", "ZORO"},
			processor: func(d interface{}) interface{} {
				s, ok := d.(string)
				if !ok {
					return d
				}
				return strings.ToUpper(s)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T){
			got := Compute(tc.input, tc.processor)
			assert.Equal(t, tc.want, got)
		})
	}
}

// Requirement Compute should print the result
func TestComputeOutput(t *testing.T) {
	cases := []struct {
		name      string
		input     []interface{}
		want      string
		processor ProcessorFunc
		separator string
	}{
		{
			name: "Print result in each line",
			input: []interface{}{1, 2, 3},
			want: "2\n3\n4\n",
			separator: "\n",
			processor: func(d interface{}) interface{} {
				n := d.(int)
				return n + 1
			},
		},
		{
			name: "Print result in csv",
			input: []interface{}{1, 2, 3},
			want: "2,3,4,",
			separator: ",",
			processor: func(d interface{}) interface{} {
				n := d.(int)
				return n + 1
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T){
			Separator = tc.separator
			out := new(bytes.Buffer)
			Fcompute(out, tc.input, tc.processor)
			got := out.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func Example_ComputePrint() {
	processorFunc := func(d interface{}) interface{} {
		n := d.(int)
		return n+1
	}
	PrintComputation([]interface{}{1, 2, 3}, processorFunc)
	// Output:
	// 2
	// 3
	// 4
}

