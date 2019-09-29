package main

type node struct {
	data  int
	left  *node
	right *node
}

func main() {
	topParent := node{
		data: 1,
		left: &node{
			data:2,
		},
		right: &node{
			data: 3,
		},
	}

	topParent.left.left = &node{data: 4}

}