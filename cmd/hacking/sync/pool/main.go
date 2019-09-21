package main

import (
	"bytes"
	"fmt"
	"sync"
)

// Tutorial from:
// https://riptutorial.com/go/example/16314/sync-pool

var pool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func main() {
	s := pool.Get().(*bytes.Buffer)

	s.Write([]byte("Luffy"))
	// put it back
	pool.Put(s)

	s = pool.Get().(*bytes.Buffer)
	s.Write([]byte(" Monkey"))

	fmt.Println(s)

	// clean up so that it can be reuse
	s.Reset()
	pool.Put(s)
	s = pool.Get().(*bytes.Buffer)
	// Defer to make sure the it don't leak
	defer pool.Put(s)
	s.Write([]byte("reset!"))
	fmt.Println(s)
}
