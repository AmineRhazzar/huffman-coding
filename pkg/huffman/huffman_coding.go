package huffman

import (
	"fmt"
	"os"
	"sort"
)

type Huffman struct {
	tree *Node
}

func (h *Huffman)DisplayTree() {
	h.tree.Display(0);
}

func (h *Huffman) constructTree(t []byte) {

	// count occurence of each byte in the string
	occurences := make(map[byte]int)
	for _, char := range t {
		count, exists := occurences[char]
		if exists {
			occurences[byte(char)] = count + 1
		} else {
			occurences[byte(char)] = 1
		}
	}

	// construct nodes from each byte, occurence
	nodes := make([]Node, len(occurences))
	i := 0
	for k, v := range occurences {
		nodes[i] = Node{
			ch:     k,
			weight: v,
		}
		i++
	}

	// sort nodes in increasing order
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].weight < nodes[j].weight
	})

	for len(nodes) != 0 {

		if len(nodes) >= 2 {
			// take the first 2 elements, create their mother node, then remove them
			// insert mother node back into nodes[] with conserving sort

			a, b := (nodes)[0], (nodes)[1]
			motherNode := Node{
				weight: a.weight + b.weight,
				Left:   &b,
				Right:  &a,
			}

			(nodes) = (nodes)[2:]

			// insert mother node into nodes[] and keep sort intact
			index := len(nodes)
			for i, node := range nodes {
				if node.weight >= motherNode.weight {
					index = i
					break
				}
			}
			Insert[Node](&nodes, index, motherNode)

		} else {
			h.tree = &(nodes[0])
			break
		}
	}
}

func (h *Huffman) Encode(inputFile string, outputFile string) (error) {
	data, err := os.ReadFile(inputFile)
	if(err != nil) {
		return err;
	}

	h.constructTree(data)
	h.tree.Display(0)
	fmt.Println(h.tree.encode())

	return nil;
}

