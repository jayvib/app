package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	fmt.Println("vim-go")
	s := "Hello"
	r := strings.NewReader(s)
	b := make([]byte, 2)
	_, err := r.Read(b)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

	b = b[:]
	_, err = r.Read(b)
	fmt.Println(string(b))
}
