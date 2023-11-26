package huffman

import (
	"os"
	"sort"
)

type Huffman struct {
	tree *Node
}

func (h *Huffman) DisplayTree() {
	h.tree.Display(0)
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

func (h *Huffman) Encode(inputFile string, outputFile string) error {
	data, read_err := os.ReadFile(inputFile)
	if read_err != nil {
		return read_err
	}

	f, write_err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0600)
	if write_err != nil {
		return write_err
	}
	defer f.Close()

	h.constructTree(data)
	h.tree.Display(0)

	w := Writer{
		debug:    true,
		ioWriter: f,
	}

	// since we use a uint32 for tree size, it's gonna be encoded on 4 bytes even if it can be encoded on less
	// we can further optimise this, we can use the first 3 bits to indicates how many bytes (SIZE_L) the size is encoded on (0->7 or 000 -> 111)
	// and then we take the next SIZE_L bytes and that's the size of our tree. it doesn't seem like it's worth it
	tree_size := w.WriteTree(h.tree)
	tree_size_bytes, byte_encode_err := ByteEncode(tree_size)

	if byte_encode_err != nil {
		return byte_encode_err
	}
	// fmt.Printf("tree size : %d, ", tree_size,)
	// for _, v := range buf.Bytes() {
	// 	fmt.Printf("%08b ", v)
	// }
	// fmt.Printf("\n")

	w.buffer = append(tree_size_bytes, w.buffer...)

	w.Flush()
	return nil
}
