package huffman

import (
	"fmt"
	"io"
)

type Writer struct {
	ioWriter io.Writer
	buffer   []byte
	currByte byte
	cursor   int8
	debug    bool
}

func (w *Writer) WriteBit(bit int) (int) {
	// if bit = 0, we just increment the cursor since currByte is always 00000000 in the beginning
	if bit == 1 {
		w.currByte = w.currByte | (0b10000000 >> w.cursor)
	}
	w.cursor++
	if w.cursor == 8 {
		w.buffer = append(w.buffer, w.currByte)
		w.cursor = 0
		w.currByte = 0
	}

	return 1;
}

func (w *Writer) WriteByte(b byte) int {
	if w.cursor == 0 {
		w.currByte = 0
		w.buffer = append(w.buffer, b)
	} else {
		// b_mask is a byte where the first (8-cursor) bits are 1 and the rest is 0
		var b_mask byte = (0xFF << w.cursor) & 0xFF

		var rest byte = b & (^b_mask) << (8 - w.cursor)
		w.buffer = append(w.buffer, (b&b_mask>>w.cursor)|w.currByte)
		w.currByte = rest
	}
	return 8
}

func (w *Writer) WriteTree(n *Node) (uint32) {
	if n.Left == nil && n.Right == nil {
		return uint32(w.WriteBit(1) + w.WriteByte(n.ch))
	}

	return uint32(w.WriteBit(0)) + w.WriteTree(n.Left) + w.WriteTree(n.Right)
}

func (w *Writer) Debug() {
	fmt.Println("+++++++++++++++++++++++++++++")
	fmt.Printf("cursor:   %d\n", w.cursor)
	fmt.Printf("currByte: %08b\n", w.currByte)
	fmt.Printf("buffer:  [ ")
	for _, b := range w.buffer {
		fmt.Printf("%08b ", b)
	}
	fmt.Printf("]\n")
	fmt.Println("+++++++++++++++++++++++++++++")
}

func (w *Writer) Flush() {
	if w.cursor != 0 {
		w.buffer = append(w.buffer, w.currByte)

	}

	// we store the cursor in the last byte so that we can determine how many bits to read in the second-to-last byte
	var lastByte byte = byte(w.cursor)
	w.buffer = append(w.buffer, lastByte)

	if w.debug {
		w.Debug()
	}

	w.ioWriter.Write(w.buffer)
	w.buffer = []byte{}
	w.currByte = 0
	w.cursor = 0
}
