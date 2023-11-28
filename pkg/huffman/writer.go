package huffman

import (
	"fmt"
	"io"
)

type Writer struct {
	io_writer io.Writer
	buffer    []byte
	curr_byte byte  // we write bits to this byte because we can't write a single bit to file
	cursor    uint8 // index of current bit in curr_byte 0..7
	debug     bool
}

func (w *Writer) WriteBit(bit uint8) int {
	// if bit = 0, we just increment the cursor since currByte is always 00000000 in the beginning
	if bit == 1 {
		w.curr_byte = w.curr_byte | (0b1000_0000 >> w.cursor)
	}
	w.cursor++
	if w.cursor == 8 {
		w.buffer = append(w.buffer, w.curr_byte)
		w.cursor = 0
		w.curr_byte = 0
	}

	return 1
}

func (w *Writer) WriteMultipleBits(bits ...uint8) {
	for _, b := range bits {
		w.WriteBit(b)
	}
}

func (w *Writer) WriteByte(b byte) int {
	if w.cursor == 0 {
		w.curr_byte = 0
		w.buffer = append(w.buffer, b)
	} else {
		// b_mask is a byte where the first (8-cursor) bits are 1 and the rest is 0
		var b_mask byte = (0xFF << w.cursor)

		var rest byte = b & (^b_mask) << (8 - w.cursor)
		w.buffer = append(w.buffer, (b&b_mask>>w.cursor)|w.curr_byte)
		w.curr_byte = rest
	}
	return 8
}

func (w *Writer) WriteTree(n *Node) uint32 {
	if n.Left == nil && n.Right == nil {
		return uint32(w.WriteBit(1) + w.WriteByte(n.ch))
	}

	return uint32(w.WriteBit(0)) + w.WriteTree(n.Left) + w.WriteTree(n.Right)
}

func (w *Writer) Flush() (int, error) {
	if w.cursor != 0 {
		w.buffer = append(w.buffer, w.curr_byte)
	}

	// we store the cursor in the last byte so that we can determine how many bits to read in the second-to-last byte
	var lastByte byte = byte(w.cursor)
	w.buffer = append(w.buffer, lastByte)
	
	if(w.debug) {
		fmt.Println("buffer:")
		for _, b:= range(w.buffer) {
			fmt.Printf("%08b ", b)
		}
		fmt.Printf("\n")
	}

	n, err := w.io_writer.Write(w.buffer)
	if err != nil {
		return 0, err
	}

	w.buffer = []byte{}
	w.curr_byte = 0
	w.cursor = 0
	return n, nil
}
