package main

import (
	"huffman-coding/pkg/huffman"
)


func main() {
	w := huffman.Writer{
	}

	w.WriteBit(1)
	w.WriteBit(1)
	w.WriteByte(0b01001001)
	w.Print()
	

}
