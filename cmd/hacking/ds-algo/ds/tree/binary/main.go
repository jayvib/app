package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Tree struct {
	Root *Node `json:"root"`
}

func (t *Tree) PrintPostOrder() {
	printPostOrder(t.Root, "root", "\n")
	fmt.Println()
}

func printPostOrder(n *Node, where string, sep string) {
	if n == nil {
		return
	}

	fmt.Print(where, " ", n.Value, sep)
	printPostOrder(n.Left, "left", "\n\t")
	printPostOrder(n.Right, "right", "\n\t")
}

type Node struct {
	Value int
	Left  *Node `json:"left"`
	Right *Node `json:"right"`
}

func levelOrderBinaryTree(arr []int, start int, size int) *Node {
	curr := &Node{arr[start], nil, nil}
	left := 2*start + 1
	right := 2*start + 2
	if left < size {
		curr.Left = levelOrderBinaryTree(arr, left, size)
	}

	if right < size {
		curr.Right = levelOrderBinaryTree(arr, right, size)
	}
	return curr
}

func LevelOrderBinaryTree(arr []int) *Tree {
	tree := new(Tree)
	tree.Root = levelOrderBinaryTree(arr, 0, len(arr))
	return tree
}

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t2 := LevelOrderBinaryTree(arr)
	t2.PrintPostOrder()
	bite, err := json.Marshal(t2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bite))
}
