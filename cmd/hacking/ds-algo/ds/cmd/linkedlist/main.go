package main

import (
	"fmt"
	"github.com/Workiva/go-datastructures/list"
	"math"
	"strings"
)

type Person struct {
	Fname, Lname string
}

func main() {
	l := list.Empty
	l = l.Add(2).Add(Person{Fname: "Luffy", Lname: "Monkey"})
	fmt.Println(l.Length())
	h, _ := l.Head()
	fmt.Println(h, l.Length())
	valList := l.Map(func(d interface{}) interface{} {
		switch v := d.(type) {
		case int:
			return v * 2
		case Person:
			return Person{
				Fname: strings.ToUpper(v.Fname),
				Lname: strings.ToUpper(v.Lname),
			}
		}
		return d
	})
	fmt.Println(valList[0])
	fmt.Println(valList[1])
	h1, _ := l.Tail()
	fmt.Println(h1.Head())
	fmt.Println(l.IsEmpty())
	l , _= l.Remove(l.Length()-1) // pop
	fmt.Println(l.Length())
	h2, _ := l.Head()
	fmt.Println(h2)
	l = l.Add(Person{Lname: "Vinsmoke", Fname: "Sanji"})
	h3, _ := l.Head()
	fmt.Println(h3)

	l, _ = l.Insert(Person{
		"Zoro",
		"Roronoa",
	}, 1)
	h4, _ := l.Get(1)
	fmt.Println(h4)
	valList = l.Map(func(d interface{}) interface{} {
		switch v := d.(type) {
		case int:
			return v * 2
		case Person:
			return Person{
				Fname: strings.ToUpper(v.Fname),
				Lname: strings.ToUpper(v.Lname),
			}
		}
		return d
	})
	fmt.Println(valList)

	intList := list.Empty
	for i := 0; i < 1000000; i++ {
		intList = intList.Add(i+1)
	}
	fmt.Println(intList.Length())
	squareList := intList.Map(func(d interface{}) interface{}{
		val := d.(int)
		return math.Pow(float64(val), 2)
	})
	fmt.Println(len(squareList))
}
