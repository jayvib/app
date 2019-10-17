package main

import (
	"fmt"
	"github.com/Workiva/go-datastructures/trie/ctrie"
	"strconv"
	"sync"
)

type Person struct {
	Name string
}

type Node struct {
	id string
	parentID string
	child []*Node
}

func main() {
	concurrency()
}

func node() {
	parent := &Node{id: "_",}
	child1 := &Node{id:"1", parentID: "_"}
	child2 := &Node{id:"2", parentID: "_"}

	trie := ctrie.New(nil)
	trie.Insert([]byte(parent.id), parent)
	trie.Insert([]byte(child1.id), child1)
	trie.Insert([]byte(child2.id), child2)

}

func simpleUsage() {
	trie := ctrie.New(nil)
	trie.Insert([]byte("Jayson"), Person{Name: "Jayson"})

	val, ok := trie.Lookup([]byte("Jayson"))
	if ok {
		p := val.(Person)
		fmt.Println(p.Name)
	}
}

func concurrency() {
	trie := ctrie.New(nil)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for i := 0; i < 10000; i++ {
			trie.Insert([]byte(strconv.Itoa(i)), i)
		}
		wg.Done()
	}()

	//go func() {
	//	for i := 10001; i < 20000; i++ {
	//		trie.Insert([]byte(strconv.Itoa(i)), i)
	//	}
	//	wg.Done()
	//}()

	//for i := 0; i < 10000; i++ {
	//	time.Sleep(5)
	//	trie.Remove([]byte(strconv.Itoa(i)))
	//}

	wg.Wait()
	fmt.Println("size", trie.Size())

	// resource leak if no cancel channel has been passed.
	doneChan := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		var v []int
		for val := range trie.Iterator(doneChan) {
			fmt.Println(val.Value.(int))
			v = append(v, val.Value.(int))
		}
		wg.Done()
		fmt.Println("total:", len(v))
	}()

	wg.Wait()
	doneChan <- struct{}{}
}