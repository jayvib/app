package main

import (
	"fmt"
	"github.com/Workiva/go-datastructures/list"
	"log"
)

func main() {
	l := list.Empty
	l = l.Add(1).Add(3).Add(4).Add(5)

	fmt.Println(l.Map(func(d interface{})interface{}{
		v := d.(int)
		return v * 2
	}))
	l, err := l.Insert(2, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(l.Length())
	val, _ := l.Get(1)
	fmt.Println("Position 1:", val)
	fmt.Println(l.Length())
	l, _ = l.Remove(1)
	fmt.Println(l.Length())

	prodheader := list.Empty
	prodheader = prodheader.
		Add("vendorPart"). // vendor part
		Add("Product Code").
		Add("Name").
		Add("Description")

	headers := prodheader.Map(func(d interface{})interface{}{
		return d
	})
	fmt.Println(headers)
}

