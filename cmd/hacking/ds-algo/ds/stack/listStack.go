package stack

import (
	"fmt"
	"github.com/Workiva/go-datastructures/list"
)

type ListStack struct {
	l list.PersistentList
}

func (l *ListStack) Push(data interface{}) {
	l.l = l.l.Add(data)
}

func (l *ListStack) Len() int {
	return int(l.l.Length())
}

func (l *ListStack) Pop() (data interface{}) {
	data, _ = l.l.Get(l.l.Length()-1)
	fmt.Println(data)
	return
}