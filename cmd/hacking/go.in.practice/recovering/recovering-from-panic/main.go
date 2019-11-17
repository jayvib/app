package main

import (
	"fmt"
)

func main() {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				fmt.Println("error:", err.Error())
			}

			if mt, ok := e.(MyStype); ok {
				fmt.Println("my type:", mt)
			}
		}
	}()
	yikes()
}

type MyStype struct {
}

func (MyStype) String() string {
	return "yeeeeyyyy my type"
}

func yikes() {
	//panic(errors.New("something bad happend."))
	panic(MyStype{})
}
