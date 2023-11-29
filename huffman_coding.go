package main

import (
	"fmt"
	"os"
	"sort"
)

type Huffman struct {
	tree  *Node
	codes map[byte]([]uint8)
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
			insert[Node](&nodes, index, motherNode)

		} else {
			h.tree = &(nodes[0])
			break
		}
	}

	if __DEBUG__ {
		h.tree.Display(0)
	}

	h.codes = make(map[byte][]uint8)
	for b := range occurences {
		var path []uint8
		h.tree.Find(b, &path, true, 0)
		h.codes[b] = path
		if __DEBUG__ {
			fmt.Printf("path(%s)=%v\n", string(b), path)
		}
	}
}

// Encodes a file into using Huffman encoding
// Returns ratio of outputsize / inputsize, and whatever error that may have resulted
func (h *Huffman) Encode(inputFile string, outputFile string) (float64, error) {
	data, read_err := os.ReadFile(inputFile)
	if read_err != nil {
		return 0, read_err
	}

	h.constructTree(data)

	w := Writer{}

	tree_size := w.WriteTree(h.tree)

	// get byte representation of tree_size, it's uint32 so it's 4 bytes
	tree_size_bytes, byte_encode_err := intToBytes(tree_size)

	if byte_encode_err != nil {
		return 0, byte_encode_err
	}
	w.buffer = append(tree_size_bytes, w.buffer...)

	for _, b := range data {
		w.WriteMultipleBits(h.codes[b]...)
	}

	f, write_err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0600)
	if write_err != nil {
		return 0, write_err
	}
	defer f.Close()
	w.io_writer = f

	n, flush_err := w.Flush()

	return float64(n) / float64(len(data)), flush_err
}

// Decodes huffman-encoded input file to output file.
// Returns ratio of outputsize / inputsize written to outputFile, and whatever error that may have resulted
func (h *Huffman) Decode(inputFile string, outputFile string) (float64, error) {
	data, read_err := os.ReadFile(inputFile)
	if read_err != nil {
		return 0, read_err
	}

	r := GetReader(data)

	// read first 4 bytes, they represent tree size in bits
	tree_size_bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b, read_tree_size_err := r.ReadByte()
		if read_tree_size_err != nil {
			return 0, fmt.Errorf("error reading tree size %v", read_tree_size_err)
		}
		tree_size_bytes[i] = b
	}

	tree_size := bytesToInt(tree_size_bytes)
	root, read_tree_err := r.ReadTree(r.idx*8+int(r.cursor), int(tree_size))

	if read_tree_err != nil {
		return 0, fmt.Errorf("error reading tree %v", read_tree_err)
	}

	if(__DEBUG__) {
		root.Display(0)
	}
	
	var decoded_data []byte
	current := root
	for {
		bit, read_bit_err := r.ReadBit()

		if(read_bit_err != nil) {
			break;
		}

		if root.isLeaf() {
			// single-noded tree
			decoded_data = append(decoded_data, root.ch)
		} else {
			if bit == 0 {
				current = current.Left
			} else {
				current = current.Right
			}

			if current.Left == nil && current.Right == nil {
				decoded_data = append(decoded_data, current.ch)
				current = root
			}
		}
	}

	f, open_write_err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0600)
	if open_write_err != nil {
		return 0, open_write_err
	}
	defer f.Close()
	n, write_err := f.Write(decoded_data)

	return float64(len(data)) / float64(n), write_err
}
