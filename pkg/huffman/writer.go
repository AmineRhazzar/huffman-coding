package huffman

import (
	"fmt"
	"io"
)

var offsets []byte = []byte{0b11111111, 0b11111110, 0b11111100, 0b11111000, 0b11110000, 0b11100000, 0b11000000, 0b10000000}

type Writer struct {
	ioWriter io.Writer
	buffer   []byte
	currByte byte
	cursor   int8
}

func (w *Writer) WriteBit(bit int) {
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
}

func (w *Writer) WriteByte(b byte) {
	if w.cursor == 0 {
		w.currByte = 0
		w.buffer = append(w.buffer, b)
	} else {
		var rest byte = b & (^offsets[w.cursor]) << (8 - w.cursor)
		w.buffer = append(w.buffer, (b & offsets[w.cursor] >> w.cursor) | w.currByte)
		w.currByte = rest
	}
}

func (w *Writer) Print() {
	fmt.Println("+++++++++++++++++++++++++++++")
	fmt.Printf("cursor:   %d\n", w.cursor)
	fmt.Printf("currByte: %08b\n", w.currByte)
	fmt.Println("buffer:  ", w.buffer)
	fmt.Println("+++++++++++++++++++++++++++++")
}

func (w *Writer) Flush() {
	w.ioWriter.Write(w.buffer)
}
