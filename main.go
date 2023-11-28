package main

import (
	"fmt"
	"huffman-coding/pkg/huffman"
)

func main() {

	h := huffman.Huffman{
		Debug: true,
	}

	err := h.Encode("example.txt", "output.huf")

	if err != nil {
		fmt.Println(err)
	}

}
