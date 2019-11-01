package main

import (
	"fmt"
	"github.com/jayvib/app/cmd/hacking/ds-algo/ds/LinkedList/singlev2/review"
)

func main() {
	l := list.List{}
	l.AddTail(1)
	l.AddHead(2)
	l.AddHead(3)
	l.AddTail(4)
	l.Print()

	fmt.Println("Is 1 present?", l.IsPresent(1))
	fmt.Println("Is 5 present?", l.IsPresent(5))
}
