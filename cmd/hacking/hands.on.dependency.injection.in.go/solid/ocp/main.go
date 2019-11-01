package main

import "fmt"

// Open-Close Principle: Software entities should be
// open for extension but close for modification.

// -------------Application-----------------


// Suppose that I have this types
type ProcessorFunc func(d interface{}) interface{}

func Compute(data []interface{}, fn ProcessorFunc) (result []interface{}) {
	for _, n := range data {
		result = append(result, fn(n))
	}
	return
}

type ComputeFunc func(data []interface{}, fn ProcessorFunc) (result []interface{})

// This is called function decorator
func PrintParamData(computeFunc ComputeFunc) (ComputeFunc) {
	return func(data []interface{}, fn ProcessorFunc) (result []interface{}) {
		fmt.Println("Data:", data)
		return computeFunc(data, fn)
	}
}

// Then I want to extend the functionality of compute....
func main() {
	data := []interface{}{1, 2, 3, 4}
	compute := PrintParamData(Compute)
	res := compute(data, addByOne)
	fmt.Println("Result:", res)
}

func addByOne(d interface{}) interface{} {
	n := d.(int)
	return n + 1
}

func multiplyByTwo(d interface{}) interface{} {
	n := d.(int)
	return n * 2
}
