package main

import "fmt"

const count int = 10

type Node struct {
	Left   *Node
	Right  *Node
	ch     byte
	weight int
}

func (n *Node) Display(space int) {
	// Base case
	if n == nil {
		return
	}
	// Increase distance between levels
	space += count

	// Process right child first
	n.Right.Display(space)

	// Print current node after space
	// count
	fmt.Printf("\n")
	for i := count; i < space; i++ {
		fmt.Printf(" ")
	}
	if n.ch == 0 {
		fmt.Printf("%d\n", n.weight)
	} else {
		fmt.Printf("\"%s\"-%d\n", string(n.ch), n.weight)
	}

	// Process left child
	n.Left.Display(space)
}

// val=0 if we go left, val=1 if we go right
func (n *Node) Find(b byte, path *[]uint8, isRoot bool, val uint8) bool {
	if n == nil {
		return false
	}
	if !isRoot {
		*path = append(*path, val)
	}
	if n.ch == b {
		return true
	}

	if n.Left.Find(b, path, false, 0) || n.Right.Find(b, path, false, 1) {
		return true
	}

	*path = (*path)[:len(*path)-1]
	return false
}
