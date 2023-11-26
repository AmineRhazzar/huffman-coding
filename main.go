package main

import (
	"fmt"
	"huffman-coding/pkg/huffman"
)


func main() {

	h := huffman.Huffman{}

	err := h.Encode("example.txt", "output.huf")

	if(err != nil ){
		fmt.Println(err)
	}

	

}
