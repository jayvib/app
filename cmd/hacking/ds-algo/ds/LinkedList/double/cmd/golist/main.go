package main

import (
	"container/list"
	"fmt"
)

func main() {
	l := list.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	_, _ = e4, e1
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)
	l.Remove(e4)
	tail := l.Back()
	l.InsertAfter(struct{ Fn, Ln string}{Fn: "Luffy", Ln: "Monkey"}, tail)
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

