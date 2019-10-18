package main

import (
	"fmt"
	"github.com/jayvib/app/apperr/v2"
	"github.com/pkg/errors"
	e "errors"
)

type GetErr struct {}

func (g GetErr) Error() string {
	return "unable to Get Item"
}

func main() {
	err1 := apperr.Error{Op: "Get", Kind: apperr.IO, Err: GetErr{}}
	err2 := apperr.Error{Op: "Put", Kind: apperr.IO, Err: err1}
	err3 := apperr.Error{Op: "pocker/search/Elasticsearch.Delete", Kind: apperr.Invalid, Err: err2}
	err4 := apperr.Error{Op: "pocker/search/Elasticsearch.Add", Kind: apperr.Invalid, Err: err3}

	cause := errors.Cause(err4)

	if e.Is(cause, GetErr{}) {
		fmt.Println("yaykksss this is an GetErr type... yohoooo using Cause")
		fmt.Println(cause)
	}

	cause2 := e.Unwrap(err4)
	if e.Is(cause2, GetErr{}) {
		fmt.Println("yaykksss this is an GetErr type... yohoooo using Unwrap")
		fmt.Println(cause)
	}
	fmt.Println(err4)
}
