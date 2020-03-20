package apperr

import (
	"errors"
	"fmt"
)

func ExampleError() {
	err := Error{Op: "Get", Kind: IO, Err: errors.New("unable to Get item")}
	fmt.Print(err)

	// Output:
	// Get: I/O error: unable to Get item
}

func ExampleErrorNested() {
	err1 := Error{Op: "Get", Kind: IO, Err: errors.New("unable to Get item")}
	err2 := Error{Op: "Put", Kind: IO, Err: err1}
	fmt.Println(err2)

	// Output:
	// Put: I/O error
	// 	Get: I/O error: unable to Get item
}

func ExampleErrorTwoNested() {
	err1 := Error{Op: "Get", Kind: IO, Err: errors.New("unable to Get item")}
	err2 := Error{Op: "Put", Kind: IO, Err: err1}
	err3 := Error{Op: "Delete", Kind: Invalid, Err: err2}
	fmt.Println(err3)

	// Output:
	// Delete: invalid operation
	// 	Put: I/O error
	// 	Get: I/O error: unable to Get item
}