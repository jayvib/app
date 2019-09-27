package stack

import "github.com/Workiva/go-datastructures/list"

type ListStack struct {
	l list.PersistentList
}

func (l *ListStack) Push(data interface{}) {
	l.l = l.l.Add(data)
}

func (l *ListStack) Len() int {
	return int(l.l.Length())
}