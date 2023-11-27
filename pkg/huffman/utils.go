package huffman

import (
	"bytes"
	"encoding/binary"
)

func insert[T any](a *[]T, index int, value T) {
	if len((*a)) == index { // nil or empty slice or after last element
		(*a) = append((*a), value)
		return
	}
	(*a) = append((*a)[:index+1], (*a)[index:]...) // index < len(a)
	(*a)[index] = value
}

func byteEncode(num uint32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func getBit(b byte, i uint8) uint8 {
	idx := i % 8
	var b_mask byte = 0b1000_0000 >> idx

	result := b & b_mask

	if result == byte(0b1000_0000>>idx) {
		return 1
	}

	return 0
}

func AreByteArraysEqual(b1, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}

	for idx, b := range b1 {
		if b != b2[idx] {
			return false
		}
	}
	return true
}

func areWritersEqual(w1, w2 Writer) bool {
	return AreByteArraysEqual(w1.buffer, w2.buffer) && w1.curr_byte == w2.curr_byte && w1.cursor == w2.cursor
}
