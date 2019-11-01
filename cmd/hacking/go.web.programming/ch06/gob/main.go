package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type P struct {
	X, Y, Z int
}

func main() {
	ps := []P{
		{X: 1, Y: 2, Z:3 },
		{X: 4, Y: 5, Z:6 },
		{X: 7, Y: 8, Z:9 },
	}
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(ps)
	if err != nil {
		log.Fatal(err)
	}

	var newPs []P
	err = gob.NewDecoder(&b).Decode(&newPs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(newPs)
}

func example1() {
	var b bytes.Buffer
	w := gob.NewEncoder(&b)
	err := w.Encode(&P{X: 1, Y: 2, Z: 3})
	if err != nil {
		log.Fatal(err)
	}

	d := gob.NewDecoder(&b)

	p := new(P)
	err = d.Decode(p)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", p)
	fmt.Println(p.X, p.Y, p.Z)

}
