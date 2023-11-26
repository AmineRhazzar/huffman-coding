package huffman

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

